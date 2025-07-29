package algorithms

import (
	"crypto/md5"
)

// ModHash 直接哈希取模算法实现
type ModHash struct {
	nodes []string // 节点列表
}

// NewModHash 创建一个新的直接哈希取模实例
func NewModHash() *ModHash {
	return &ModHash{
		nodes: make([]string, 0),
	}
}

// hash 计算key的hash值
func (mh *ModHash) hash(key string) uint64 {
	h := md5.Sum([]byte(key))
	// 确保结果为正数
	result := (uint64(h[0]) << 56) | (uint64(h[1]) << 48) | (uint64(h[2]) << 40) | (uint64(h[3]) << 32) |
		(uint64(h[4]) << 24) | (uint64(h[5]) << 16) | (uint64(h[6]) << 8) | uint64(h[7])
	return result & 0x7fffffffffffffff // 清除符号位确保为正数
}

// AddNode 添加节点
func (mh *ModHash) AddNode(node string) {
	for _, n := range mh.nodes {
		if n == node {
			return
		}
	}
	mh.nodes = append(mh.nodes, node)
}

// RemoveNode 移除节点
func (mh *ModHash) RemoveNode(node string) {
	index := -1
	for i, n := range mh.nodes {
		if n == node {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	mh.nodes = append(mh.nodes[:index], mh.nodes[index+1:]...)
}

// GetNode 根据key获取对应的节点
func (mh *ModHash) GetNode(key string) string {
	if len(mh.nodes) == 0 {
		return ""
	}

	keyHash := mh.hash(key)
	idx := keyHash % uint64(len(mh.nodes))
	return mh.nodes[idx]
}

// GetNodesCount 获取当前节点数量
func (mh *ModHash) GetNodesCount() int {
	return len(mh.nodes)
}