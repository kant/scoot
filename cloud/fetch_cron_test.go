package cloud

import (
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestFetchCron(t *testing.T) {
	h := makeCronHelper(t)
	h.assertFetch(t)
	h.assertFetch(t, "host1:1234")
	h.assertFetch(t, "host1:1234", "host2:8888")
	h.assertFetch(t)
	h.assertFetch(t, "host1:1234")
	h.assertFetch(t)
}

type cronHelper struct {
	t      *testing.T
	tickCh chan time.Time
	f      *fakeFetcher
	ch     chan ClusterUpdate
}

func makeCronHelper(t *testing.T) *cronHelper {
	h := &cronHelper{
		t:      t,
		tickCh: make(chan time.Time),
		f:      &fakeFetcher{},
	}
	h.ch = MakeFetchCron(h.f, h.tickCh)
	return h
}

func (h *cronHelper) assertFetch(t *testing.T, expectedNames ...string) {
	nodes := nodes(expectedNames)
	sort.Sort(NodeSorter(nodes))
	h.f.setResult(nodes)
	expected := nodes
	h.tickCh <- time.Now()
	actual := <-h.ch
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("got %v, expected %v", actual, expected)
	}
}

func nodes(ids []string) []Node {
	n := []Node{}
	for _, name := range ids {
		n = append(n, NewIdNode(name))
	}
	return n
}

// fakeFetcher for testing fetch cron
type fakeFetcher struct {
	mutex sync.Mutex
	nodes []Node
}

func (f *fakeFetcher) Fetch() ([]Node, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.nodes, nil
}

func (f *fakeFetcher) setResult(nodes []Node) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.nodes = nodes
}
