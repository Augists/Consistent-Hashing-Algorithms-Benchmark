package algorithms

import (
	"crypto/md5"
)

// DxHash DxHash一致性哈希算法实现
type DxHash struct {
	nodes    []string // 节点列表
	nodeMap  map[string]bool // 节点映射，用于快速检查节点是否已存在
	nsTable  []int    // NSArray查找表
	nsSize   int      // 当前NSArray大小
	nodeCount int     // 当前节点数量
	availableStack []int // 可用位置的栈，用于优化节点添加/删除
}

// NewDxHash 创建一个新的DxHash实例，使用默认参数
func NewDxHash() *DxHash {
	return NewDxHashWithParams(8) // 默认初始大小为8
}

// NewDxHashWithParams 创建一个新的DxHash实例，可以指定初始大小
func NewDxHashWithParams(initialSize int) *DxHash {
	// 确保初始大小是2的幂次
	size := 1
	for size < initialSize {
		size <<= 1
	}
	
	dx := &DxHash{
		nodes:    make([]string, size), // 预分配节点数组
		nodeMap:  make(map[string]bool),
		nsTable:  make([]int, size),
		nsSize:   size,
		nodeCount: 0, // 当前节点数量
		availableStack: make([]int, 0, size),
	}
	
	// 初始化可用位置栈
	for i := 0; i < size; i++ {
		dx.availableStack = append(dx.availableStack, i)
	}
	
	return dx
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
	newNodes := make([]string, newSize) // 新的节点数组

	// 复制旧数据
	copy(newTable, dx.nsTable)
	copy(newNodes, dx.nodes)

	// 初始化新部分
	for i := dx.nsSize; i < newSize; i++ {
		newTable[i] = -1
		// 将新位置添加到可用栈
		dx.availableStack = append(dx.availableStack, i)
	}

	dx.nsTable = newTable
	dx.nodes = newNodes
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
	// 检查节点是否已存在
	if _, exists := dx.nodeMap[node]; exists {
		return
	}
	
	// 检查是否需要扩容
	if dx.nodeCount >= dx.nsSize/2 { // 负载因子超过0.5时扩容
		dx.resize()
	}
	
	// 从栈中获取一个可用位置
	if len(dx.availableStack) == 0 {
		return
	}
	
	// 从栈顶取出一个可用位置
	index := dx.availableStack[len(dx.availableStack)-1]
	dx.availableStack = dx.availableStack[:len(dx.availableStack)-1]
	
	// 在获取的位置添加节点
	dx.nodes[index] = node
	dx.nodeMap[node] = true
	
	// 在NSArray中为节点分配位置
	pos := int(dx.hash(node, 0)) % dx.nsSize
	dx.nsTable[pos] = index
	
	dx.nodeCount++
}

// RemoveNode 移除节点
func (dx *DxHash) RemoveNode(node string) {
	if _, exists := dx.nodeMap[node]; !exists {
		return
	}
	
	delete(dx.nodeMap, node)
	
	index := -1
	for i := 0; i < dx.nodeCount; i++ {
		if dx.nodes[i] == node {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	// 从NSArray中移除节点位置
	pos := int(dx.hash(node, 0)) % dx.nsSize
	dx.nsTable[pos] = -1
	
	// 将位置放回可用栈
	dx.availableStack = append(dx.availableStack, index)
	
	// 减少节点计数
	dx.nodeCount--
}

// GetNode 根据key获取对应的节点
func (dx *DxHash) GetNode(key string) string {
	if dx.nodeCount == 0 {
		return ""
	}
	
	// 不预先生成伪随机序列，而是按需计算
	// 最大重试次数为8*当前节点数，按照论文建议
	maxRetries := 8 * dx.nodeCount
	for i := 0; i < maxRetries; i++ {
		pos := int(dx.hash(key, i)) % dx.nsSize
		if dx.nsTable[pos] != -1 && dx.nsTable[pos] < len(dx.nodes) {
			// 检查该位置是否有效
			nodeIndex := dx.nsTable[pos]
			if nodeIndex >= 0 && nodeIndex < len(dx.nodes) && dx.nodes[nodeIndex] != "" {
				return dx.nodes[nodeIndex]
			}
		}
	}
	
	// 如果都没找到，返回第一个有效节点
	for i := 0; i < dx.nodeCount; i++ {
		if dx.nodes[i] != "" {
			return dx.nodes[i]
		}
	}
	
	return ""
}

// GetNodesCount 获取当前节点数量
func (dx *DxHash) GetNodesCount() int {
	return dx.nodeCount
}

// GetNodes 获取所有节点列表
func (dx *DxHash) GetNodes() []string {
	nodes := make([]string, dx.nodeCount)
	for i := 0; i < dx.nodeCount; i++ {
		nodes[i] = dx.nodes[i]
	}
	return nodes
}

// GetTableSize 获取查找表大小
func (dx *DxHash) GetTableSize() int {
	return dx.nsSize
}
