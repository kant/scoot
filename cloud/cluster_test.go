package cloud

import (
	"testing"
)

func TestMembers(t *testing.T) {
	h := makeHelper(t)
	defer h.close()
	h.assertMembers()
	h.add("node1")
	h.assertMembers("node1")
	h.remove("node1")
	h.assertMembers()
	h.add("node1", "node2")
	h.assertMembers("node1", "node2")
	h.add("node1")
	h.assertMembers("node1", "node2")
	h.remove("node1", "node2")
	h.assertMembers()
	// Try confusing it by removing a nonexisting node
	h.remove("node3")
	h.assertMembers()
	h.add("node3", "node4")
	h.assertMembers("node3", "node4")
}

func TestSubscribe(t *testing.T) {
	h := makeHelper(t)
	defer h.close()
	h.assertMembers()
	s := h.subscribe()
	defer s.Closer.Close()
	h.assertInitialMembers(s)
	h.add("node1")
	h.assertUpdates(s, add("node1"))
	h.add("node2")
	h.add("node3")
	h.remove("node1")
	h.assertUpdates(s, add("node2"), add("node3"), remove("node1"))

	// Add a second subscription
	s2 := h.subscribe()
	defer s2.Closer.Close()
	h.assertInitialMembers(s2, "node2", "node3")

	// Now test that a subscriber that's not pulling doesn't block others
	h.add("node4")
	h.assertUpdates(s, add("node4"))
	h.remove("node2")
	h.assertUpdates(s, remove("node2"))
	h.assertMembers("node3", "node4")
	h.assertUpdates(s2, add("node4"), remove("node2"))

	// Test a subscription gets updates when state is changed directly
	h.changeStateTo("node3")
	h.assertMembers("node3")
	h.assertUpdates(s, remove("node4"))
	h.assertUpdates(s2, remove("node4"))

	h.changeStateTo("node1", "node2")
	h.assertMembers("node1", "node2")
	h.assertUpdates(s, add("node1"), add("node2"), remove("node3"))
	h.assertUpdates(s2, add("node1"), add("node2"), remove("node3"))
}

// Below here are helpers that make it easy to write more fluent tests.

type helper struct {
	t        *testing.T
	c        *Cluster
	updateCh chan []NodeUpdate
	f        Fetcher
	ch       chan ClusterUpdate
}

func makeHelper(t *testing.T) *helper {
	h := &helper{t: t}
	h.ch = make(chan ClusterUpdate)
	h.c = NewCluster(nil, h.ch)
	return h
}

func (h *helper) close() {
	h.c.Close()
}

func (h *helper) assertMembers(node ...string) {
	assertMembersEqual(makeNodes(node...), h.c.Members(), h.t)
}

func (h *helper) add(node ...string) {
	updates := []NodeUpdate{}
	for _, n := range node {
		updates = append(updates, add(n))
	}
	h.ch <- updates
}

func (h *helper) remove(node ...string) {
	updates := []NodeUpdate{}
	for _, n := range node {
		updates = append(updates, remove(n))
	}
	h.ch <- updates
}

func (h *helper) changeStateTo(node ...string) {
	nodes := makeNodes(node...)
	h.ch <- nodes
}

func (h *helper) subscribe() Subscription {
	return h.c.Subscribe()
}

func (h *helper) assertInitialMembers(s Subscription, node ...string) {
	assertMembersEqual(s.InitialMembers, makeNodes(node...), h.t)
}

func (h *helper) assertUpdates(s Subscription, expected ...NodeUpdate) {
	// Calling Members makes sure that any updates sent have propagated from the cluster to the subscription
	h.c.Members()
	actual := <-s.Updates
	h.assertUpdatesEqual(expected, actual)
}

func (h *helper) assertUpdatesEqual(expected []NodeUpdate, actual []NodeUpdate) {
	if len(expected) != len(actual) {
		h.t.Fatalf("unequal updates: %v %v", expected, actual)
	}
	for i, ex := range expected {
		act := actual[i]
		if ex.UpdateType != act.UpdateType || ex.Id != act.Id {
			h.t.Fatalf("unequal updates: %v %v (from %v %v)", ex, act, expected, actual)
		}
	}
}

func makeNodes(node ...string) []Node {
	r := []Node{}
	for _, n := range node {
		r = append(r, NewIdNode(n))
	}
	return r
}

func add(node string) NodeUpdate {
	return NewAdd(NewIdNode(node))
}

func remove(node string) NodeUpdate {
	return NewRemove(NodeId(node))
}

func assertMembersEqual(expected []Node, actual []Node, t *testing.T) {
	if len(expected) != len(actual) {
		t.Fatalf("unequal members: %v %v", expected, actual)
	}
	for i, ex := range expected {
		act := actual[i]
		if ex.Id() != act.Id() {
			t.Fatalf("unequal members: %v %v (from %v %v)", ex.Id(), act.Id(), expected, actual)
		}
	}
}
