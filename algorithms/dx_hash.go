package algorithms

import (
	"crypto/md5"
)

// DxHash DxHash一致性哈希算法实现
type DxHash struct {
	nodes    []string // 节点列表
	nsTable  []int    // NSArray查找表
	nsSize   int      // 当前NSArray大小
	maxRetry int      // 最大重试次数
}

// NewDxHash 创建一个新的DxHash实例
func NewDxHash() *DxHash {
	initialSize := 8 // 初始大小为2^n
	return &DxHash{
		nodes:    make([]string, 0),
		nsTable:  make([]int, initialSize),
		nsSize:   initialSize,
		maxRetry: 8, // 默认最大重试次数
	}
}

// hash 计算key的hash值
func (dx *DxHash) hash(key string, seed int) uint64 {
	h := md5.Sum([]byte(key + string(rune(seed))))
	// 确保结果为正数
	result := (uint64(h[0]) << 56) | (uint64(h[1]) << 48) | (uint64(h[2]) << 40) | (uint64(h[3]) << 32) |
		(uint64(h[4]) << 24) | (uint64(h[5]) << 16) | (uint64(h[6]) << 8) | uint64(h[7])
	return result & 0x7fffffffffffffff // 清除符号位确保为正数
}

// resize 扩容NSArray
func (dx *DxHash) resize() {
	newSize := dx.nsSize * 2
	newTable := make([]int, newSize)
	
	// 复制旧数据
	copy(newTable, dx.nsTable)
	
	// 初始化新部分
	for i := dx.nsSize; i < newSize; i++ {
		newTable[i] = -1
	}
	
	dx.nsTable = newTable
	dx.nsSize = newSize
}

// populateTable 填充查找表
func (dx *DxHash) populateTable() {
	// 初始化表
	for i := range dx.nsTable {
		dx.nsTable[i] = -1
	}
	
	// 为每个节点在表中分配位置
	for i, node := range dx.nodes {
		// 使用节点索引作为种子计算位置
		pos := int(dx.hash(node, 0)) % dx.nsSize
		dx.nsTable[pos] = i
	}
}

// AddNode 添加节点
func (dx *DxHash) AddNode(node string) {
	for _, n := range dx.nodes {
		if n == node {
			return
		}
	}
	
	dx.nodes = append(dx.nodes, node)
	
	// 检查是否需要扩容
	if len(dx.nodes) > dx.nsSize/2 { // 负载因子超过0.5时扩容
		dx.resize()
	}
	
	// 重新填充查找表
	dx.populateTable()
}

// RemoveNode 移除节点
func (dx *DxHash) RemoveNode(node string) {
	index := -1
	for i, n := range dx.nodes {
		if n == node {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	// 从节点列表中移除
	dx.nodes = append(dx.nodes[:index], dx.nodes[index+1:]...)
	
	// 重新填充查找表
	dx.populateTable()
}

// GetNode 根据key获取对应的节点
func (dx *DxHash) GetNode(key string) string {
	if len(dx.nodes) == 0 {
		return ""
	}
	
	// 生成伪随机序列
	positions := make([]int, 0, dx.maxRetry)
	for i := 0; i < dx.maxRetry && i < len(dx.nodes)*8; i++ {
		pos := int(dx.hash(key, i)) % dx.nsSize
		positions = append(positions, pos)
	}
	
	// 按顺序查找活跃节点
	for _, pos := range positions {
		if dx.nsTable[pos] != -1 && dx.nsTable[pos] < len(dx.nodes) {
			return dx.nodes[dx.nsTable[pos]]
		}
	}
	
	// 如果都没找到，返回第一个节点
	return dx.nodes[0]
}

// GetNodesCount 获取当前节点数量
func (dx *DxHash) GetNodesCount() int {
	return len(dx.nodes)
}

// GetTableSize 获取查找表大小
func (dx *DxHash) GetTableSize() int {
	return dx.nsSize
}