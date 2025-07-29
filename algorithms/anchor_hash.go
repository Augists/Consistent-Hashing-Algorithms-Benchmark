package algorithms

import (
	"crypto/md5"
)

// AnchorHash AnchorHash一致性哈希算法实现
type AnchorHash struct {
	workers []bool   // 工作节点标记数组
	removed []bool   // 删除节点标记数组
	next    []int    // next指针数组
	nodes   []string // 节点名称数组
	size    int      // 预期最大节点数
	active  int      // 当前活跃节点数
}

// NewAnchorHash 创建一个新的AnchorHash实例
func NewAnchorHash(size int) *AnchorHash {
	ah := &AnchorHash{
		workers: make([]bool, size),
		removed: make([]bool, size),
		next:    make([]int, size),
		nodes:   make([]string, size),
		size:    size,
		active:  0,
	}

	// 初始化next数组
	for i := 0; i < size; i++ {
		ah.next[i] = i + 1
	}
	ah.next[size-1] = 0 // 最后一个指向第一个形成环

	return ah
}

// hash 计算key的hash值
func (ah *AnchorHash) hash(key string, seed int) uint64 {
	h := md5.Sum([]byte(key + string(rune(seed))))
	// 确保结果为正数
	result := (uint64(h[0]) << 56) | (uint64(h[1]) << 48) | (uint64(h[2]) << 40) | (uint64(h[3]) << 32) |
		(uint64(h[4]) << 24) | (uint64(h[5]) << 16) | (uint64(h[6]) << 8) | uint64(h[7])
	return result & 0x7fffffffffffffff // 清除符号位确保为正数
}

// AddNode 添加节点
func (ah *AnchorHash) AddNode(node string) {
	// 查找一个空闲位置
	index := -1
	for i := 0; i < ah.size; i++ {
		if ah.nodes[i] == "" {
			index = i
			break
		}
	}

	if index == -1 {
		// 没有空闲位置
		return
	}

	ah.nodes[index] = node
	ah.workers[index] = true
	ah.removed[index] = false
	ah.active++
}

// RemoveNode 移除节点
func (ah *AnchorHash) RemoveNode(node string) {
	index := -1
	for i := 0; i < ah.size; i++ {
		if ah.nodes[i] == node {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	ah.workers[index] = false
	ah.removed[index] = true
	ah.nodes[index] = ""
	ah.active--
}

// GetNode 根据key获取对应的节点
func (ah *AnchorHash) GetNode(key string) string {
	if ah.active == 0 {
		return ""
	}

	b := int(ah.hash(key, 0)) % ah.size // 初始桶

	// 如果桶是工作节点，直接返回
	if ah.workers[b] {
		return ah.nodes[b]
	}

	// 否则启动回填过程
	start := b
	b = ah.next[b] // 跟随next指针

	// 查找下一个工作节点
	for b != start {
		if ah.workers[b] {
			return ah.nodes[b]
		}
		b = ah.next[b]
	}

	return ""
}

// GetNodesCount 获取当前节点数量
func (ah *AnchorHash) GetNodesCount() int {
	return ah.active
}

// GetSize 获取哈希空间大小
func (ah *AnchorHash) GetSize() int {
	return ah.size
}