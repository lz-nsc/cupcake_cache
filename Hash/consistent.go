package hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFunc func(data []byte) uint32

type Manager struct {
	hashFn   HashFunc
	replicas int // number or virtual node for each node
	nodes    []int
	hashMap  map[int]string
}

func NewManager(replicas int, hashFn HashFunc) *Manager {
	if hashFn == nil {
		hashFn = crc32.ChecksumIEEE
	}

	return &Manager{
		hashFn:   hashFn,
		replicas: replicas,
		nodes:    make([]int, 0),
		hashMap:  map[int]string{},
	}
}

func (m *Manager) AddNodes(nodes ...string) {
	for _, name := range nodes {
		for i := 0; i < m.replicas; i++ {
			hash := m.hashFn([]byte(strconv.Itoa(i) + name))
			m.nodes = append(m.nodes, int(hash))
			m.hashMap[int(hash)] = name
		}
	}
	sort.Ints(m.nodes)
}

func (m *Manager) GetNode(key string) string {
	if len(m.nodes) == 0 {
		return ""
	}
	hash := m.hashFn([]byte(key))
	nodeHash := sort.Search(len(m.nodes), func(i int) bool {
		return m.nodes[i] >= int(hash)
	})

	return m.hashMap[m.nodes[nodeHash%len(m.nodes)]]
}
