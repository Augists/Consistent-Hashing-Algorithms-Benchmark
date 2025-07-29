package algorithms

import (
	"crypto/md5"
)

// MPCH Multi-Probe一致性哈希算法实现
type MPCH struct {
	nodes  []string // 节点列表
	k      int      // 探针数量
}

// DefaultK 默认探针数量
const DefaultK int = 21

// NewMPCH 创建一个新的Multi-Probe一致性哈希实例
func NewMPCH(k int) *MPCH {
	if k <= 0 {
		k = DefaultK
	}
	return &MPCH{
		nodes: make([]string, 0),
		k:     k, // 探针数量
	}
}

// hash 计算key的hash值
func (mpch *MPCH) hash(key string, seed int) uint64 {
	h := md5.Sum([]byte(key + string(rune(seed))))
	// 确保结果为正数
	result := (uint64(h[0]) << 56) | (uint64(h[1]) << 48) | (uint64(h[2]) << 40) | (uint64(h[3]) << 32) |
		(uint64(h[4]) << 24) | (uint64(h[5]) << 16) | (uint64(h[6]) << 8) | uint64(h[7])
	return result & 0x7fffffffffffffff // 清除符号位确保为正数
}

// AddNode 添加节点
func (mpch *MPCH) AddNode(node string) {
	for _, n := range mpch.nodes {
		if n == node {
			return
		}
	}
	mpch.nodes = append(mpch.nodes, node)
}

// RemoveNode 移除节点
func (mpch *MPCH) RemoveNode(node string) {
	index := -1
	for i, n := range mpch.nodes {
		if n == node {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	mpch.nodes = append(mpch.nodes[:index], mpch.nodes[index+1:]...)
}

// GetNode 根据key获取对应的节点
func (mpch *MPCH) GetNode(key string) string {
	if len(mpch.nodes) == 0 {
		return ""
	}

	if len(mpch.nodes) == 1 {
		return mpch.nodes[0]
	}

	// 使用k个探针对key进行哈希，找到最匹配的节点
	bestNode := ""
	bestScore := uint64(0)

	// 对key进行k次哈希（使用不同的种子）
	for i := 0; i < mpch.k; i++ {
		hash := mpch.hash(key, i)
		
		// 对每个节点计算得分
		for _, node := range mpch.nodes {
			nodeHash := mpch.hash(node, 0)
			
			// 计算匹配得分（较高的哈希值表示更好的匹配）
			score := uint64(0)
			if hash >= nodeHash {
				score = hash - nodeHash
			} else {
				// 环形距离
				score = (^uint64(0)) - nodeHash + hash
			}
			
			// 选择得分最高的节点
			if score > bestScore {
				bestScore = score
				bestNode = node
			}
		}
	}

	return bestNode
}

// GetNodesCount 获取当前节点数量
func (mpch *MPCH) GetNodesCount() int {
	return len(mpch.nodes)
}

// GetProbeCount 获取探针数量
func (mpch *MPCH) GetProbeCount() int {
	if mpch.k <= 0 {
		return DefaultK
	}
	return mpch.k
}