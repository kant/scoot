package scootconfig

import (
	"fmt"
	"time"

	"github.com/twitter/scoot/cloud"
	"github.com/twitter/scoot/cloud/local"
	"github.com/twitter/scoot/ice"
)

// Parameters for configuring an in-memory Scoot cluster
// Count - number of in-memory workers
type ClusterMemoryConfig struct {
	Type  string
	Count int
}

func (c *ClusterMemoryConfig) Install(bag *ice.MagicBag) {
	bag.Put(c.Create)
}

func (c *ClusterMemoryConfig) Create() (*cloud.Cluster, error) {
	workerNodes := make([]cloud.Node, c.Count)
	for i := 0; i < c.Count; i++ {
		workerNodes[i] = cloud.NewIdNode(fmt.Sprintf("inmemory%d", i))
	}
	return cloud.NewCluster(workerNodes, nil), nil
}

// Parameters for configuring a Scoot cluster that will have locally-run components.
type ClusterLocalConfig struct {
	Type string
}

func (c *ClusterLocalConfig) Install(bag *ice.MagicBag) {
	bag.Put(c.Create)
}

func (c *ClusterLocalConfig) Create() (*cloud.Cluster, error) {
	f := local.MakeFetcher("workerserver", "thrift_addr")
	updates := cloud.MakeFetchCron(f, time.NewTicker(time.Second).C)
	return cloud.NewCluster(nil, updates), nil
}
