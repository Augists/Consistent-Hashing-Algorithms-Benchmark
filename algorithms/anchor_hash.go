package algorithms

import (
	"crypto/md5"
)

// AnchorHash AnchorHash一致性哈希算法实现
type AnchorHash struct {
	a       []int          // anchor array A (锚点数组)
	w       []int          // working set W (工作集)
	l       []int          // location array L (位置数组)
	k       []int          // kick array K (踢出数组)
	r       []int          // removed stack R (移除栈)
	n       int            // current working set size (当前工作集大小)
	maxSize int            // max bucket capacity |A| (最大桶容量)
	nodes   []string       // node names (节点名称)
	nodeMap map[string]int // node name to bucket mapping (节点名称到桶的映射)
}

// NewAnchorHash 创建一个新的AnchorHash实例
func NewAnchorHash(size int) *AnchorHash {
	ah := &AnchorHash{
		a:       make([]int, size),
		w:       make([]int, size),
		l:       make([]int, size),
		k:       make([]int, size),
		r:       make([]int, 0, size),
		n:       0,
		maxSize: size,
		nodes:   make([]string, size),
		nodeMap: make(map[string]int),
	}

	// Initialize arrays
	for i := 0; i < size; i++ {
		ah.k[i], ah.w[i], ah.l[i] = i, i, i
		ah.a[i] = size // Initially all buckets are removed
	}

	// Initially all buckets are in removed stack
	for i := size - 1; i >= 0; i-- {
		ah.r = append(ah.r, i)
	}

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
	// 检查节点是否已存在
	if _, exists := ah.nodeMap[node]; exists {
		return
	}

	if ah.n >= ah.maxSize {
		// 达到最大节点数
		return
	}

	// 从移除栈中取出一个桶
	b := ah.r[len(ah.r)-1]
	ah.r = ah.r[:len(ah.r)-1]

	// 添加节点到工作集
	ah.a[b] = 0
	ah.l[ah.w[ah.n]] = ah.n
	ah.w[ah.l[b]], ah.k[b] = b, b
	ah.nodes[b] = node
	ah.nodeMap[node] = b
	ah.n++
}

// RemoveNode 移除节点
func (ah *AnchorHash) RemoveNode(node string) {
	// 检查节点是否存在
	b, exists := ah.nodeMap[node]
	if !exists {
		return
	}

	// 将桶放回移除栈
	ah.r = append(ah.r, b)
	ah.n--
	ah.a[b] = ah.n
	ah.w[ah.l[b]], ah.k[b] = ah.w[ah.n], ah.w[ah.n]
	ah.l[ah.w[ah.n]] = ah.l[b]

	// 从节点映射中删除
	delete(ah.nodeMap, node)
	ah.nodes[b] = ""
}

// GetNode 根据key获取对应的节点
func (ah *AnchorHash) GetNode(key string) string {
	if ah.n == 0 {
		return ""
	}

	b := int(ah.hash(key, 0)) % ah.maxSize
	iterations := 0
	maxIterations := ah.maxSize * 2 // 防止无限循环
	
	for ah.a[b] > 0 && iterations < maxIterations {
		h := int(ah.hash(key, b+1)) % ah.a[b]
		steps := 0
		maxSteps := ah.maxSize // 防止内部循环无限
		
		for ah.a[h] >= ah.a[b] && steps < maxSteps {
			h = ah.k[h]
			steps++
			if h < 0 || h >= len(ah.k) {
				// 防止数组越界
				return ah.nodes[b]
			}
		}
		
		if steps >= maxSteps {
			// 内部循环超时，返回当前桶
			return ah.nodes[b]
		}
		
		b = h
		iterations++
	}
	
	if iterations >= maxIterations {
		// 外部循环超时，返回默认节点
		if ah.n > 0 {
			return ah.nodes[ah.w[0]]
		}
		return ""
	}
	
	// 确保返回的索引有效
	if b >= 0 && b < len(ah.nodes) {
		return ah.nodes[b]
	}
	return ""
}

// GetNodesCount 获取当前节点数量
func (ah *AnchorHash) GetNodesCount() int {
	return ah.n
}

// GetSize 获取哈希空间大小
func (ah *AnchorHash) GetSize() int {
	return ah.maxSize
}
