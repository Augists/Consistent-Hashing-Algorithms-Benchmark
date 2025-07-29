package algorithms

import (
	"crypto/md5"
)

// MaglevHash Maglev一致性哈希算法实现
type MaglevHash struct {
	nodes       []string // 节点列表
	tableSize   int      // 查找表大小
	lookupTable []string // 查找表
}

// NewMaglevHash 创建一个新的Maglev哈希
func NewMaglevHash(tableSize int) *MaglevHash {
	return &MaglevHash{
		nodes:       make([]string, 0),
		tableSize:   tableSize,
		lookupTable: make([]string, tableSize),
	}
}

// hash 计算key的hash值
func (mh *MaglevHash) hash(key string, seed int) uint64 {
	h := md5.Sum([]byte(key + string(rune(seed))))
	// 确保结果为正数
	result := (uint64(h[0]) << 56) | (uint64(h[1]) << 48) | (uint64(h[2]) << 40) | (uint64(h[3]) << 32) |
		(uint64(h[4]) << 24) | (uint64(h[5]) << 16) | (uint64(h[6]) << 8) | uint64(h[7])
	return result & 0x7fffffffffffffff // 清除符号位确保为正数
}

// generatePermutation 为节点生成偏好序列
func (mh *MaglevHash) generatePermutation(node string, tableSize int) []int {
	offset := int(mh.hash(node, 0)) % tableSize
	skip := int(mh.hash(node, 1))%(tableSize-1) + 1

	permutation := make([]int, tableSize)
	for j := 0; j < tableSize; j++ {
		permutation[j] = (offset + j*skip) % tableSize
	}
	return permutation
}

// populateLookupTable 构建查找表
func (mh *MaglevHash) populateLookupTable() {
	if len(mh.nodes) == 0 {
		for i := range mh.lookupTable {
			mh.lookupTable[i] = ""
		}
		return
	}

	// 为每个节点生成偏好序列
	permutations := make([][]int, len(mh.nodes))
	for i, node := range mh.nodes {
		permutations[i] = mh.generatePermutation(node, mh.tableSize)
	}

	// 初始化查找表
	for i := range mh.lookupTable {
		mh.lookupTable[i] = ""
	}
	nextIdx := make([]int, len(mh.nodes)) // 每个节点下一个填充位置

	// 按轮次填充查找表
	filledCount := 0
	round := 0
	maxRounds := mh.tableSize * 10 // 防止无限循环
	for filledCount < mh.tableSize && round < maxRounds {
		for i, node := range mh.nodes {
			if filledCount >= mh.tableSize {
				break
			}

			// 找到下一个可填充的位置
			for nextIdx[i] < mh.tableSize && mh.lookupTable[permutations[i][nextIdx[i]]] != "" {
				nextIdx[i]++
			}

			// 如果还有可填充的位置，则填充
			if nextIdx[i] < mh.tableSize {
				// 直接填充，不需要检查是否为空，因为上面的循环已经确保找到空位
				mh.lookupTable[permutations[i][nextIdx[i]]] = node
				nextIdx[i]++
				filledCount++
			}
		}
		round++
	}
}

// AddNode 添加节点
func (mh *MaglevHash) AddNode(node string) {
	for _, n := range mh.nodes {
		if n == node {
			return
		}
	}
	mh.nodes = append(mh.nodes, node)
	mh.populateLookupTable()
}

// RemoveNode 移除节点
func (mh *MaglevHash) RemoveNode(node string) {
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
	mh.populateLookupTable()
}

// GetNode 根据key获取对应的节点
func (mh *MaglevHash) GetNode(key string) string {
	if len(mh.lookupTable) == 0 {
		return ""
	}

	keyHash := int(mh.hash(key, 0)) % mh.tableSize
	return mh.lookupTable[keyHash]
}

// GetNodesCount 获取当前节点数量
func (mh *MaglevHash) GetNodesCount() int {
	return len(mh.nodes)
}

// GetLookupTableSize 获取查找表大小
func (mh *MaglevHash) GetLookupTableSize() int {
	return len(mh.lookupTable)
}