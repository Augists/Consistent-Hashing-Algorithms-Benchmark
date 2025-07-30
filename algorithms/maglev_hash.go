package algorithms

import (
	"crypto/md5"
)

// NodePreference 节点偏好信息
type NodePreference struct {
	offset  int // 偏移量
	skip    int // 跳跃步长
	nextIdx int // 下一个排列序号
}

// MaglevHash Maglev一致性哈希算法实现
type MaglevHash struct {
	nodes       []string         // 节点列表
	nodeMap     map[string]bool  // 节点映射，用于快速检查节点是否已存在
	preferences []NodePreference // 节点偏好信息
	tableSize   int              // 查找表大小
	lookupTable []string         // 查找表
}

// NewMaglevHash 创建一个新的Maglev哈希
func NewMaglevHash(tableSize int) *MaglevHash {
	return &MaglevHash{
		nodes:       make([]string, 0),
		nodeMap:     make(map[string]bool),
		preferences: make([]NodePreference, 0),
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

// calculatePreference 计算节点的偏好参数
func (mh *MaglevHash) calculatePreference(node string) NodePreference {
	offset := int(mh.hash(node, 0)) % mh.tableSize
	skip := int(mh.hash(node, 1))%(mh.tableSize-1) + 1
	return NodePreference{
		offset:  offset,
		skip:    skip,
		nextIdx: 0,
	}
}


// getPermutationItem 按需计算排列序列中的指定项
func (mh *MaglevHash) getPermutationItem(preference NodePreference, index int) int {
	return (preference.offset + index*preference.skip) % mh.tableSize
}

// populateLookupTable 构建查找表
func (mh *MaglevHash) populateLookupTable() {
	if len(mh.nodes) == 0 {
		for i := range mh.lookupTable {
			mh.lookupTable[i] = ""
		}
		return
	}

	// 初始化节点偏好信息
	mh.preferences = make([]NodePreference, len(mh.nodes))
	for i, node := range mh.nodes {
		mh.preferences[i] = mh.calculatePreference(node)
	}

	// 初始化查找表
	for i := range mh.lookupTable {
		mh.lookupTable[i] = ""
	}
	
	// 重置每个节点的下一个填充位置
	for i := range mh.preferences {
		mh.preferences[i].nextIdx = 0
	}

	// 按轮次填充查找表
	filledCount := 0
	round := 0
	maxRounds := mh.tableSize * 10 // 防止无限循环
	for filledCount < mh.tableSize && round < maxRounds {
		for i := range mh.nodes {
			if filledCount >= mh.tableSize {
				break
			}

			// 找到下一个可填充的位置
			for mh.preferences[i].nextIdx < mh.tableSize && mh.lookupTable[mh.getPermutationItem(mh.preferences[i], mh.preferences[i].nextIdx)] != "" {
				mh.preferences[i].nextIdx++
			}

			// 如果还有可填充的位置，则填充
			if mh.preferences[i].nextIdx < mh.tableSize {
				pos := mh.getPermutationItem(mh.preferences[i], mh.preferences[i].nextIdx)
				mh.lookupTable[pos] = mh.nodes[i]
				mh.preferences[i].nextIdx++
				filledCount++
			}
		}
		round++
	}
}

// AddNode 添加节点
func (mh *MaglevHash) AddNode(node string) {
	if _, exists := mh.nodeMap[node]; exists {
		return
	}
	
	mh.nodeMap[node] = true
	mh.nodes = append(mh.nodes, node)
	mh.populateLookupTable()
}

// RemoveNode 移除节点
func (mh *MaglevHash) RemoveNode(node string) {
	if _, exists := mh.nodeMap[node]; !exists {
		return
	}
	
	delete(mh.nodeMap, node)
	
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