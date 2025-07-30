package algorithms

import (
	"crypto/md5"
)

// JumpHash 跳跃一致性哈希算法实现
type JumpHash struct {
	nodes []string // 节点列表
	nodeMap map[string]bool // 节点映射，用于快速检查节点是否已存在
}

// NewJumpHash 创建一个新的跳跃哈希
func NewJumpHash() *JumpHash {
	return &JumpHash{
		nodes: make([]string, 0),
		nodeMap: make(map[string]bool),
	}
}

// hash 计算key的hash值
func (jh *JumpHash) hash(key string) uint64 {
	h := md5.Sum([]byte(key))
	return (uint64(h[0]) << 56) | (uint64(h[1]) << 48) | (uint64(h[2]) << 40) | (uint64(h[3]) << 32) |
		(uint64(h[4]) << 24) | (uint64(h[5]) << 16) | (uint64(h[6]) << 8) | uint64(h[7])
}

// jumpConsistentHash 跳跃一致性哈希算法核心实现
func (jh *JumpHash) jumpConsistentHash(keyHash uint64, numBuckets int) int {
	if numBuckets <= 0 {
		return -1
	}

	b := int64(-1)
	j := int64(0)
	for j < int64(numBuckets) {
		b = j
		keyHash = keyHash*2862933555777941757 + 1
		j = int64((float64(b+1) * float64(1<<31)) / float64((keyHash>>33)+1))
	}
	return int(b)
}

// AddNode 添加节点
func (jh *JumpHash) AddNode(node string) {
	if _, exists := jh.nodeMap[node]; exists {
		return
	}
	
	jh.nodeMap[node] = true
	jh.nodes = append(jh.nodes, node)
}

// RemoveNode 移除节点
func (jh *JumpHash) RemoveNode(node string) {
	if _, exists := jh.nodeMap[node]; !exists {
		return
	}
	
	delete(jh.nodeMap, node)
	
	for i, n := range jh.nodes {
		if n == node {
			jh.nodes = append(jh.nodes[:i], jh.nodes[i+1:]...)
			return
		}
	}
}

// GetNode 根据key获取对应的节点
func (jh *JumpHash) GetNode(key string) string {
	if len(jh.nodes) == 0 {
		return ""
	}

	keyHash := jh.hash(key)
	idx := jh.jumpConsistentHash(keyHash, len(jh.nodes))
	return jh.nodes[idx]
}

// GetNodesCount 获取当前节点数量
func (jh *JumpHash) GetNodesCount() int {
	return len(jh.nodes)
}