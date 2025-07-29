package algorithms

import (
	"crypto/md5"
)

// RendezvousHash Rendezvous哈希算法实现（最高随机权重哈希）
type RendezvousHash struct {
	nodes []string // 节点列表
}

// NewRendezvousHash 创建一个新的Rendezvous哈希
func NewRendezvousHash() *RendezvousHash {
	return &RendezvousHash{
		nodes: make([]string, 0),
	}
}

// hash 计算key的hash值
func (rh *RendezvousHash) hash(key string) uint64 {
	h := md5.Sum([]byte(key))
	return (uint64(h[0]) << 56) | (uint64(h[1]) << 48) | (uint64(h[2]) << 40) | (uint64(h[3]) << 32) |
		(uint64(h[4]) << 24) | (uint64(h[5]) << 16) | (uint64(h[6]) << 8) | uint64(h[7])
}

// computeWeight 计算key-node对的权重
func (rh *RendezvousHash) computeWeight(key string, node string) uint64 {
	// 使用hash(key + node)作为权重
	combined := key + "#" + node
	return rh.hash(combined)
}

// AddNode 添加节点
func (rh *RendezvousHash) AddNode(node string) {
	for _, n := range rh.nodes {
		if n == node {
			return
		}
	}
	rh.nodes = append(rh.nodes, node)
}

// RemoveNode 移除节点
func (rh *RendezvousHash) RemoveNode(node string) {
	index := -1
	for i, n := range rh.nodes {
		if n == node {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	rh.nodes = append(rh.nodes[:index], rh.nodes[index+1:]...)
}

// GetNode 根据key获取对应的节点
func (rh *RendezvousHash) GetNode(key string) string {
	if len(rh.nodes) == 0 {
		return ""
	}

	maxWeight := uint64(0)
	selectedNode := ""

	for _, node := range rh.nodes {
		weight := rh.computeWeight(key, node)
		if weight > maxWeight {
			maxWeight = weight
			selectedNode = node
		}
	}

	return selectedNode
}

// GetNodesCount 获取当前节点数量
func (rh *RendezvousHash) GetNodesCount() int {
	return len(rh.nodes)
}