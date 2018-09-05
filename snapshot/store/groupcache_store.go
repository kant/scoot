package store

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/twitter/groupcache"
	"github.com/twitter/scoot/cloud/cluster"
	"github.com/twitter/scoot/common/stats"
)

//TODO: we should consider modifying google groupcache lib further to:
// 1) It makes more sense given our use-case to cache bundles loaded via peer 100% of the time (currently 10%).
// 2) Modify peer proto to support setting bundle data on the peer that owns the bundlename. (via PopulateCache()).
// 3) For populate cache requests from user, fill the right cache, main or hot.
//
//TODO: Add a doneCh/Done() to stop the created goroutine.

// Called periodically in a goroutine. Must include the current instance among the fetched nodes.
type PeerFetcher interface {
	Fetch() ([]cluster.Node, error)
}

// Note: Endpoint is concatenated with Name in groupcache internals, and AddrSelf is expected as HOST:PORT.
type GroupcacheConfig struct {
	Name         string
	Memory_bytes int64
	AddrSelf     string
	Endpoint     string
	Cluster      *cluster.Cluster
}

// Add in-memory caching to the given store.
func MakeGroupcacheStore(underlying Store, cfg *GroupcacheConfig, stat stats.StatsReceiver) (Store, http.Handler, error) {
	stat = stat.Scope("bundlestoreCache")
	go stats.StartUptimeReporting(stat, stats.BundlestoreUptime_ms, "", stats.DefaultStartupGaugeSpikeLen)

	// Create the cache which knows how to retrieve the underlying bundle data.
	var cache = groupcache.NewGroup(cfg.Name, cfg.Memory_bytes, groupcache.GetterFunc(
		func(ctx groupcache.Context, bundleName string, dest groupcache.Sink) error {
			log.Info("Not cached, try to fetch bundle and populate cache: ", bundleName)
			stat.Counter(stats.GroupcacheReadUnderlyingCounter).Inc(1)
			defer stat.Latency(stats.GroupcacheReadUnderlyingLatency_ms).Time().Stop()
			reader, err := underlying.OpenForRead(bundleName)
			if err != nil {
				return err
			}
			defer reader.Close()
			data, err := ioutil.ReadAll(reader)
			if err != nil {
				return err
			}
			dest.SetBytes(data)
			return nil
		},
	))

	// Create and initialize peer group.
	// The HTTPPool constructor will register as a global PeerPicker on our behalf.
	poolOpts := &groupcache.HTTPPoolOptions{BasePath: cfg.Endpoint}
	pool := groupcache.NewHTTPPoolOpts("http://"+cfg.AddrSelf, poolOpts)
	go loop(cfg.Cluster, pool, cache, stat)

	return &groupcacheStore{underlying: underlying, cache: cache, stat: stat}, pool, nil
}

// Convert 'host:port' node ids to the format expected by groupcache peering, http URLs.
func toPeers(nodes []cluster.Node, stat stats.StatsReceiver) []string {
	peers := []string{}
	for _, node := range nodes {
		peers = append(peers, "http://"+string(node.Id()))
	}
	log.Info("New groupcacheStore peers: ", peers)
	stat.Counter(stats.GroupcachePeerDiscoveryCounter).Inc(1)
	stat.Gauge(stats.GroupcachePeerCountGauge).Update(int64(len(peers)))
	return peers
}

// Loop will listen for cluster updates and create a list of peer addresses to update groupcache.
// Cluster is expected to include the current node.
// Also updates cache stats, every 1s for now to account for arbitrary stat latch time.
func loop(c *cluster.Cluster, pool *groupcache.HTTPPool, cache *groupcache.Group, stat stats.StatsReceiver) {
	sub := c.Subscribe()
	pool.Set(toPeers(c.Members(), stat)...)
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-sub.Updates:
			pool.Set(toPeers(c.Members(), stat)...)
		case <-ticker.C:
			updateCacheStats(cache, stat)
		}
	}
}

// The groupcache lib updates its stats in the background - we need to convert those to our own stat representation.
// Gauges are expected to fluctuate, counters are expected to only ever increase.
func updateCacheStats(cache *groupcache.Group, stat stats.StatsReceiver) {
	stat.Gauge(stats.GroupcacheMainBytesGauge).Update(cache.CacheStats(groupcache.MainCache).Bytes)
	stat.Gauge(stats.GroupcacheMainItemsGauge).Update(cache.CacheStats(groupcache.MainCache).Items)
	stat.Counter(stats.GroupcacheMainGetsCounter).Update(cache.CacheStats(groupcache.MainCache).Gets)
	stat.Counter(stats.GroupcacheMainHitsCounter).Update(cache.CacheStats(groupcache.MainCache).Hits)
	stat.Counter(stats.GroupcacheMainEvictionsCounter).Update(cache.CacheStats(groupcache.MainCache).Evictions)

	stat.Gauge(stats.GroupcacheHotBytesGauge).Update(cache.CacheStats(groupcache.HotCache).Bytes)
	stat.Gauge(stats.GroupcacheHotItemsGauge).Update(cache.CacheStats(groupcache.HotCache).Items)
	stat.Counter(stats.GroupcacheHotGetsCounter).Update(cache.CacheStats(groupcache.HotCache).Gets)
	stat.Counter(stats.GroupcacheHotHitsCounter).Update(cache.CacheStats(groupcache.HotCache).Hits)
	stat.Counter(stats.GroupcacheHotEvictionsCounter).Update(cache.CacheStats(groupcache.HotCache).Evictions)

	stat.Counter(stats.GroupcacheGetCounter).Update(cache.Stats.Gets.Get())
	stat.Counter(stats.GroupcacheHitCounter).Update(cache.Stats.CacheHits.Get())
	stat.Counter(stats.GroupcacheLoadCounter).Update(cache.Stats.Loads.Get())
	stat.Counter(stats.GroupcachePeerGetsCounter).Update(cache.Stats.PeerLoads.Get())
	stat.Counter(stats.GroupcachPeerErrCounter).Update(cache.Stats.PeerErrors.Get())
	stat.Counter(stats.GroupcacheLocalLoadCounter).Update(cache.Stats.LocalLoads.Get())
	stat.Counter(stats.GroupcacheLocalLoadErrCounter).Update(cache.Stats.LocalLoadErrs.Get())
	stat.Counter(stats.GroupcacheIncomingRequestsCounter).Update(cache.Stats.ServerRequests.Get())
}

type groupcacheStore struct {
	underlying Store
	cache      *groupcache.Group
	stat       stats.StatsReceiver
}

func (s *groupcacheStore) OpenForRead(name string) (io.ReadCloser, error) {
	log.Info("Read() checking for cached bundle: ", name)
	defer s.stat.Latency(stats.GroupcacheReadLatency_ms).Time().Stop()
	s.stat.Counter(stats.GroupcacheReadCounter).Inc(1)
	var data []byte
	if err := s.cache.Get(nil, name, groupcache.AllocatingByteSliceSink(&data)); err != nil {
		return nil, err
	}
	s.stat.Counter(stats.GroupcacheReadOkCounter).Inc(1)
	return ioutil.NopCloser(bytes.NewReader(data)), nil
}

func (s *groupcacheStore) Exists(name string) (bool, error) {
	log.Info("Exists() checking for cached bundle: ", name)
	defer s.stat.Latency(stats.GroupcachExistsLatency_ms).Time().Stop()
	s.stat.Counter(stats.GroupcacheExistsCounter).Inc(1)
	if err := s.cache.Get(nil, name, groupcache.TruncatingByteSliceSink(&[]byte{})); err != nil {
		return false, nil
	}
	s.stat.Counter(stats.GroupcacheExistsOkCounter).Inc(1)
	return true, nil
}

func (s *groupcacheStore) Write(name string, data io.Reader, ttl *TTLValue) error {
	log.Info("Write() populating cache: ", name)
	defer s.stat.Latency(stats.GroupcacheWriteLatency_ms).Time().Stop()
	s.stat.Counter(stats.GroupcacheWriteCounter).Inc(1)

	// Read data into a []byte and make a right-sized copy, as ReadAll will reserve at 2x capacity
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}
	c := make([]byte, len(b))
	copy(c, b)

	err = s.underlying.Write(name, data, ttl)
	if err != nil {
		return err
	}

	s.cache.PopulateCache(name, c)
	s.stat.Counter(stats.GroupcacheWriteOkCounter).Inc(1)
	return nil
}

func (s *groupcacheStore) Root() string {
	return s.underlying.Root()
}
