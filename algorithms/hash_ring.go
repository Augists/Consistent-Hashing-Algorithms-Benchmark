package algorithms

import (
	"crypto/md5"
	"sort"
)

// HashRing 经典哈希环一致性哈希算法实现
type HashRing struct {
	nodes        map[string]bool     // 节点集合
	ring         map[uint32]string   // hash -> node 映射
	sortedKeys   []uint32            // 排序后的hash值
	virtualNodes int                 // 每个真实节点对应的虚拟节点数
}

// NewHashRing 创建一个新的哈希环
func NewHashRing(virtualNodes int) *HashRing {
	return &HashRing{
		nodes:        make(map[string]bool),
		ring:         make(map[uint32]string),
		sortedKeys:   make([]uint32, 0),
		virtualNodes: virtualNodes,
	}
}

// hash 计算key的hash值
func (hr *HashRing) hash(key string) uint32 {
	h := md5.Sum([]byte(key))
	return (uint32(h[0]) << 24) | (uint32(h[1]) << 16) | (uint32(h[2]) << 8) | uint32(h[3])
}

// generateVirtualNodes 为节点生成虚拟节点hash值
func (hr *HashRing) generateVirtualNodes(node string) []uint32 {
	hashes := make([]uint32, 0, hr.virtualNodes)
	for i := 0; i < hr.virtualNodes; i++ {
		key := node + "#" + string(rune(i))
		hashes = append(hashes, hr.hash(key))
	}
	return hashes
}

// AddNode 添加节点到哈希环
func (hr *HashRing) AddNode(node string) {
	if _, exists := hr.nodes[node]; exists {
		return
	}

	hr.nodes[node] = true
	virtualHashes := hr.generateVirtualNodes(node)

	for _, hash := range virtualHashes {
		hr.ring[hash] = node
	}

	// 重新排序hash值
	hr.sortedKeys = make([]uint32, 0, len(hr.ring))
	for k := range hr.ring {
		hr.sortedKeys = append(hr.sortedKeys, k)
	}
	sort.Slice(hr.sortedKeys, func(i, j int) bool {
		return hr.sortedKeys[i] < hr.sortedKeys[j]
	})
}

// RemoveNode 从哈希环中移除节点
func (hr *HashRing) RemoveNode(node string) {
	if _, exists := hr.nodes[node]; !exists {
		return
	}

	delete(hr.nodes, node)
	virtualHashes := hr.generateVirtualNodes(node)

	for _, hash := range virtualHashes {
		delete(hr.ring, hash)
	}

	// 重新排序hash值
	hr.sortedKeys = make([]uint32, 0, len(hr.ring))
	for k := range hr.ring {
		hr.sortedKeys = append(hr.sortedKeys, k)
	}
	sort.Slice(hr.sortedKeys, func(i, j int) bool {
		return hr.sortedKeys[i] < hr.sortedKeys[j]
	})
}

// GetNode 根据key获取对应的节点
func (hr *HashRing) GetNode(key string) string {
	if len(hr.ring) == 0 {
		return ""
	}

	hash := hr.hash(key)
	// 使用二分查找找到对应的节点
	idx := sort.Search(len(hr.sortedKeys), func(i int) bool {
		return hr.sortedKeys[i] >= hash
	})

	if idx == len(hr.sortedKeys) {
		idx = 0
	}

	return hr.ring[hr.sortedKeys[idx]]
}

// GetNodesCount 获取当前节点数量
func (hr *HashRing) GetNodesCount() int {
	return len(hr.nodes)
}

// GetSortedKeysLength 获取排序键的数量（用于测试）
func (hr *HashRing) GetSortedKeysLength() int {
	return len(hr.sortedKeys)
}