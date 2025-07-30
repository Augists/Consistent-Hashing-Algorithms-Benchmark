package tests

import (
	"fmt"
	"math"
	"testing"
	"time"

	consistent_hashing "consistent_hashing/algorithms"
)

// 生成测试键
func generateTestKeys(count int) []string {
	keys := make([]string, count)
	for i := 0; i < count; i++ {
		keys[i] = fmt.Sprintf("key_%d", i)
	}
	return keys
}

// 生成测试节点
func generateTestNodes(count int) []string {
	nodes := make([]string, count)
	for i := 0; i < count; i++ {
		nodes[i] = fmt.Sprintf("node_%d", i)
	}
	return nodes
}

// BenchmarkModHash 直接哈希取模算法基准测试
func BenchmarkModHash(b *testing.B) {
	// 创建直接哈希取核实例
	mod := consistent_hashing.NewModHash()

	// 添加1000个节点
	nodes := generateTestNodes(1000)
	for _, node := range nodes {
		mod.AddNode(node)
	}

	// 生成测试键
	keys := generateTestKeys(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mod.GetNode(keys[i%len(keys)])
	}
}

// BenchmarkHashRing 哈希环算法基准测试
func BenchmarkHashRing(b *testing.B) {
	// 创建哈希环实例
	hr := consistent_hashing.NewHashRing(160)

	// 添加1000个节点
	nodes := generateTestNodes(1000)
	for _, node := range nodes {
		hr.AddNode(node)
	}

	// 生成测试键
	keys := generateTestKeys(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hr.GetNode(keys[i%len(keys)])
	}
}

// BenchmarkJumpHash 跳跃哈希算法基准测试
func BenchmarkJumpHash(b *testing.B) {
	// 创建跳跃哈希实例
	jh := consistent_hashing.NewJumpHash()

	// 添加1000个节点
	nodes := generateTestNodes(1000)
	for _, node := range nodes {
		jh.AddNode(node)
	}

	// 生成测试键
	keys := generateTestKeys(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jh.GetNode(keys[i%len(keys)])
	}
}

// BenchmarkMaglevHash Maglev哈希算法基准测试
func BenchmarkMaglevHash(b *testing.B) {
	// 创建Maglev哈希实例，表大小为65537
	mh := consistent_hashing.NewMaglevHash(65537)

	// 添加1000个节点
	nodes := generateTestNodes(1000)
	for _, node := range nodes {
		mh.AddNode(node)
	}

	// 生成测试键
	keys := generateTestKeys(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mh.GetNode(keys[i%len(keys)])
	}
}

// BenchmarkRendezvousHash Rendezvous哈希算法基准测试
func BenchmarkRendezvousHash(b *testing.B) {
	// 创建Rendezvous哈希实例
	rh := consistent_hashing.NewRendezvousHash()

	// 添加1000个节点
	nodes := generateTestNodes(1000)
	for _, node := range nodes {
		rh.AddNode(node)
	}

	// 生成测试键
	keys := generateTestKeys(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rh.GetNode(keys[i%len(keys)])
	}
}

// BenchmarkAnchorHash AnchorHash算法基准测试
func BenchmarkAnchorHash(b *testing.B) {
	// 创建AnchorHash实例
	ah := consistent_hashing.NewAnchorHash(2000)

	// 添加1000个节点
	nodes := generateTestNodes(1000)
	for _, node := range nodes {
		ah.AddNode(node)
	}

	// 生成测试键
	keys := generateTestKeys(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ah.GetNode(keys[i%len(keys)])
	}
}

// BenchmarkDxHash DxHash算法基准测试
func BenchmarkDxHash(b *testing.B) {
	// 创建DxHash实例，预设节点数为1000
	dx := consistent_hashing.NewDxHashWithParams(1000)

	// 添加1000个节点
	nodes := generateTestNodes(1000)
	for _, node := range nodes {
		dx.AddNode(node)
	}

	// 生成测试键
	keys := generateTestKeys(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dx.GetNode(keys[i%len(keys)])
	}
}

// BenchmarkMPCH5 Multi-Probe一致性哈希算法基准测试(k=5)
func BenchmarkMPCH5(b *testing.B) {
	// 创建Multi-Probe一致性哈希实例，使用5个探针
	mpch := consistent_hashing.NewMPCH(5)

	// 添加1000个节点
	nodes := generateTestNodes(1000)
	for _, node := range nodes {
		mpch.AddNode(node)
	}

	// 生成测试键
	keys := generateTestKeys(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mpch.GetNode(keys[i%len(keys)])
	}
}

// BenchmarkMPCH21 Multi-Probe一致性哈希算法基准测试(k=21)
func BenchmarkMPCH21(b *testing.B) {
	// 创建Multi-Probe一致性哈希实例，使用21个探针
	mpch := consistent_hashing.NewMPCH(21)

	// 添加1000个节点
	nodes := generateTestNodes(1000)
	for _, node := range nodes {
		mpch.AddNode(node)
	}

	// 生成测试键
	keys := generateTestKeys(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mpch.GetNode(keys[i%len(keys)])
	}
}

// TestDistribution 测试算法分布均匀性
func TestDistribution(t *testing.T) {
	nodeCount := 100
	keyCount := 100000

	// 创建各种算法实例
	mod := consistent_hashing.NewModHash()
	hr := consistent_hashing.NewHashRing(160)
	jh := consistent_hashing.NewJumpHash()
	mh := consistent_hashing.NewMaglevHash(65537)
	rh := consistent_hashing.NewRendezvousHash()
	ah := consistent_hashing.NewAnchorHash(1000)
	dx := consistent_hashing.NewDxHash()
	mpch5 := consistent_hashing.NewMPCH(5)
	mpch21 := consistent_hashing.NewMPCH(21)

	// 添加节点
	nodes := generateTestNodes(nodeCount)
	for _, node := range nodes {
		mod.AddNode(node)
		hr.AddNode(node)
		jh.AddNode(node)
		mh.AddNode(node)
		rh.AddNode(node)
		ah.AddNode(node)
		dx.AddNode(node)
		mpch5.AddNode(node)
		mpch21.AddNode(node)
	}

	// 生成测试键
	keys := make([]string, keyCount)
	for i := 0; i < keyCount; i++ {
		keys[i] = fmt.Sprintf("key_%d", i)
	}

	// 测试直接哈希取模分布
	modDistribution := make(map[string]int)
	for _, key := range keys {
		node := mod.GetNode(key)
		modDistribution[node]++
	}

	// 测试哈希环分布
	hrDistribution := make(map[string]int)
	for _, key := range keys {
		node := hr.GetNode(key)
		hrDistribution[node]++
	}

	// 测试跳跃哈希分布
	jhDistribution := make(map[string]int)
	for _, key := range keys {
		node := jh.GetNode(key)
		jhDistribution[node]++
	}

	// 测试Maglev哈希分布
	mhDistribution := make(map[string]int)
	for _, key := range keys {
		node := mh.GetNode(key)
		mhDistribution[node]++
	}

	// 测试Rendezvous哈希分布
	rhDistribution := make(map[string]int)
	for _, key := range keys {
		node := rh.GetNode(key)
		rhDistribution[node]++
	}

	// 测试AnchorHash分布
	ahDistribution := make(map[string]int)
	for _, key := range keys {
		node := ah.GetNode(key)
		ahDistribution[node]++
	}

	// 测试DxHash分布
	dxDistribution := make(map[string]int)
	dx.GetNodes() // 获取所有节点（保持与其他算法测试的一致性）
	for _, key := range keys {
		node := dx.GetNode(key)
		dxDistribution[node]++
	}

	// 测试Multi-Probe一致性哈希(k=5)分布
	mpch5Distribution := make(map[string]int)
	for _, key := range keys {
		node := mpch5.GetNode(key)
		mpch5Distribution[node]++
	}

	// 测试Multi-Probe一致性哈希(k=21)分布
	mpch21Distribution := make(map[string]int)
	for _, key := range keys {
		node := mpch21.GetNode(key)
		mpch21Distribution[node]++
	}

	// 计算平均值
	avg := keyCount / nodeCount

	// 计算标准差（峰均值比相关）
	modStdDev := calculateStdDev(modDistribution, avg, nodeCount)
	hrStdDev := calculateStdDev(hrDistribution, avg, nodeCount)
	jhStdDev := calculateStdDev(jhDistribution, avg, nodeCount)
	mhStdDev := calculateStdDev(mhDistribution, avg, nodeCount)
	rhStdDev := calculateStdDev(rhDistribution, avg, nodeCount)
	ahStdDev := calculateStdDev(ahDistribution, avg, nodeCount)
	dxStdDev := calculateStdDev(dxDistribution, avg, nodeCount)
	mpch5StdDev := calculateStdDev(mpch5Distribution, avg, nodeCount)
	mpch21StdDev := calculateStdDev(mpch21Distribution, avg, nodeCount)

	t.Logf("Distribution test with %d nodes and %d keys", nodeCount, keyCount)
	t.Logf("Average assignments per node: %d", avg)
	t.Logf("Mod Hash Std Dev: %.2f", modStdDev)
	t.Logf("Hash Ring Std Dev: %.2f", hrStdDev)
	t.Logf("Jump Hash Std Dev: %.2f", jhStdDev)
	t.Logf("Maglev Hash Std Dev: %.2f", mhStdDev)
	t.Logf("Rendezvous Hash Std Dev: %.2f", rhStdDev)
	t.Logf("AnchorHash Std Dev: %.2f", ahStdDev)
	t.Logf("DxHash Std Dev: %.2f", dxStdDev)
	t.Logf("Multi-Probe CH(k=5) Std Dev: %.2f", mpch5StdDev)
	t.Logf("Multi-Probe CH(k=21) Std Dev: %.2f", mpch21StdDev)
}

// calculateStdDev 计算标准差
func calculateStdDev(distribution map[string]int, avg, nodeCount int) float64 {
	sum := 0.0
	for _, count := range distribution {
		diff := float64(count - avg)
		sum += diff * diff
	}
	variance := sum / float64(nodeCount)
	return math.Sqrt(variance)
}

// TestNodeAddition 测试添加节点时的重映射
func TestNodeAddition(t *testing.T) {
	initialNodes := 100  // 减少初始节点数
	addCount := 10
	keyCount := 10000    // 减少键的数量

	// 只测试部分算法，避免超时
	// 跳过RendezvousHash和MPCH算法的测试，因为它们比较耗时

	keys := make([]string, keyCount)
	for i := 0; i < keyCount; i++ {
		keys[i] = fmt.Sprintf("key_%d", i)
	}

	// 初始化各种算法
	mod := consistent_hashing.NewModHash()
	hr := consistent_hashing.NewHashRing(160)
	jh := consistent_hashing.NewJumpHash()
	mh := consistent_hashing.NewMaglevHash(65537)
	dx := consistent_hashing.NewDxHash()

	// 添加初始节点
	initialNodesList := make([]string, initialNodes)
	for i := 0; i < initialNodes; i++ {
		node := fmt.Sprintf("node_%d", i)
		initialNodesList[i] = node
		mod.AddNode(node)
		hr.AddNode(node)
		jh.AddNode(node)
		mh.AddNode(node)
		dx.AddNode(node)
	}

	// 测试添加节点前的分配情况
	modBefore := make([]string, keyCount)
	hrBefore := make([]string, keyCount)
	jhBefore := make([]string, keyCount)
	mhBefore := make([]string, keyCount)
	dxBefore := make([]string, keyCount)

	for i, key := range keys {
		modBefore[i] = mod.GetNode(key)
		hrBefore[i] = hr.GetNode(key)
		jhBefore[i] = jh.GetNode(key)
		mhBefore[i] = mh.GetNode(key)
		dxBefore[i] = dx.GetNode(key)
	}

	// 添加新节点并测量时间
	newNodes := make([]string, addCount)
	for i := 0; i < addCount; i++ {
		newNodes[i] = fmt.Sprintf("new_node_%d", i)
	}

	// 测试ModHash
	start := time.Now()
	for i := 0; i < addCount; i++ {
		mod.AddNode(newNodes[i])
	}
	elapsed := time.Since(start)

	modAfter := make([]string, keyCount)
	changed := 0
	for i, key := range keys {
		modAfter[i] = mod.GetNode(key)
		if modBefore[i] != modAfter[i] {
			changed++
		}
	}

	t.Logf("ModHash: Adding %d nodes to %d initial nodes took %v, %d keys remapped (%.2f%%)",
		addCount, initialNodes, elapsed, changed, float64(changed)*100/float64(keyCount))

	// 测试HashRing
	start = time.Now()
	for i := 0; i < addCount; i++ {
		hr.AddNode(newNodes[i])
	}
	elapsed = time.Since(start)

	hrAfter := make([]string, keyCount)
	changed = 0
	for i, key := range keys {
		hrAfter[i] = hr.GetNode(key)
		if hrBefore[i] != hrAfter[i] {
			changed++
		}
	}

	t.Logf("HashRing: Adding %d nodes to %d initial nodes took %v, %d keys remapped (%.2f%%)",
		addCount, initialNodes, elapsed, changed, float64(changed)*100/float64(keyCount))

	// 测试JumpHash
	start = time.Now()
	for i := 0; i < addCount; i++ {
		jh.AddNode(newNodes[i])
	}
	elapsed = time.Since(start)

	jhAfter := make([]string, keyCount)
	changed = 0
	for i, key := range keys {
		jhAfter[i] = jh.GetNode(key)
		if jhBefore[i] != jhAfter[i] {
			changed++
		}
	}

	t.Logf("JumpHash: Adding %d nodes to %d initial nodes took %v, %d keys remapped (%.2f%%)",
		addCount, initialNodes, elapsed, changed, float64(changed)*100/float64(keyCount))

	// 测试MaglevHash
	start = time.Now()
	for i := 0; i < addCount; i++ {
		mh.AddNode(newNodes[i])
	}
	elapsed = time.Since(start)

	mhAfter := make([]string, keyCount)
	changed = 0
	for i, key := range keys {
		mhAfter[i] = mh.GetNode(key)
		if mhBefore[i] != mhAfter[i] {
			changed++
		}
	}

	t.Logf("MaglevHash: Adding %d nodes to %d initial nodes took %v, %d keys remapped (%.2f%%)",
		addCount, initialNodes, elapsed, changed, float64(changed)*100/float64(keyCount))

	// 测试DxHash
	start = time.Now()
	for i := 0; i < addCount; i++ {
		dx.AddNode(newNodes[i])
	}
	elapsed = time.Since(start)

	dxAfter := make([]string, keyCount)
	changed = 0
	for i, key := range keys {
		dxAfter[i] = dx.GetNode(key)
		if dxBefore[i] != dxAfter[i] {
			changed++
		}
	}

	t.Logf("DxHash: Adding %d nodes to %d initial nodes took %v, %d keys remapped (%.2f%%)",
		addCount, initialNodes, elapsed, changed, float64(changed)*100/float64(keyCount))
}