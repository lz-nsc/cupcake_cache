package hash

import (
	"strconv"
	"testing"
)

func TestConsistent(t *testing.T) {
	// Use hash function that the result of test case can be predict
	mgr := NewManager(3, func(node []byte) uint32 {
		hash, _ := strconv.Atoi(string(node))
		return uint32(hash)
	})

	// virtual nodes: "1":"1","11","20";"10":"10","110","210";"30":"30","130","230"
	mgr.AddNodes("1", "10", "30")

	testCases := map[string]string{
		"350": "1",
		"12":  "1",
		"100": "10",
		"111": "30",
	}

	for key, res := range testCases {
		node := mgr.GetNode(key)
		if res != node {
			t.Errorf("failed to get node, want: %s, got: %s", res, node)
		}
	}
	// virtual nodes: "15":"115","215","30"
	mgr.AddNodes("15")
	if node := mgr.GetNode("111"); node != "15" {
		t.Errorf("failed to get node after added new node, want: 15, got: %s", node)
	}

}
