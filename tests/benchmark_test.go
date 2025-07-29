package tests

import (
	"fmt"
	"math"
	"testing"
	"time"

	"consistent_hashing/algorithms"
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
	mod := algorithms.NewModHash()

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
	hr := algorithms.NewHashRing(160)

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
	jh := algorithms.NewJumpHash()

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
	mh := algorithms.NewMaglevHash(65537)

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
	rh := algorithms.NewRendezvousHash()

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
	ah := algorithms.NewAnchorHash(2000)

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
	// 创建DxHash实例
	dx := algorithms.NewDxHash()

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
	mpch := algorithms.NewMPCH(5)

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
	mpch := algorithms.NewMPCH(21)

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
	mod := algorithms.NewModHash()
	hr := algorithms.NewHashRing(160)
	jh := algorithms.NewJumpHash()
	mh := algorithms.NewMaglevHash(65537)
	rh := algorithms.NewRendezvousHash()
	ah := algorithms.NewAnchorHash(1000)
	dx := algorithms.NewDxHash()
	mpch5 := algorithms.NewMPCH(5)
	mpch21 := algorithms.NewMPCH(21)

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
	keys := generateTestKeys(keyCount)

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
	initialNodes := 1000
	addNodes := []int{1, 5, 10, 50, 100}
	keyCount := 100000

	// 创建各种算法实例
	mod := algorithms.NewModHash()
	hr := algorithms.NewHashRing(160)
	jh := algorithms.NewJumpHash()
	mh := algorithms.NewMaglevHash(65537)
	rh := algorithms.NewRendezvousHash()
	ah := algorithms.NewAnchorHash(2000)
	dx := algorithms.NewDxHash()
	mpch5 := algorithms.NewMPCH(5)
	mpch21 := algorithms.NewMPCH(21)

	// 添加初始节点
	nodes := generateTestNodes(initialNodes)
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
	keys := generateTestKeys(keyCount)

	for _, addCount := range addNodes {
		// 测试直接哈希取模
		modBefore := make([]string, keyCount)
		for i, key := range keys {
			modBefore[i] = mod.GetNode(key)
		}

		newNodes := generateTestNodes(initialNodes + addCount)
		start := time.Now()
		for i := initialNodes; i < initialNodes+addCount; i++ {
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

		t.Logf("Mod Hash: Adding %d nodes to %d initial nodes took %v, %d keys remapped (%.2f%%)",
			addCount, initialNodes, elapsed, changed, float64(changed)*100/float64(keyCount))

		// 测试哈希环
		hrBefore := make([]string, keyCount)
		for i, key := range keys {
			hrBefore[i] = hr.GetNode(key)
		}

		start = time.Now()
		for i := initialNodes; i < initialNodes+addCount; i++ {
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

		t.Logf("Hash Ring: Adding %d nodes to %d initial nodes took %v, %d keys remapped (%.2f%%)",
			addCount, initialNodes, elapsed, changed, float64(changed)*100/float64(keyCount))

		// 测试跳跃哈希
		jhBefore := make([]string, keyCount)
		for i, key := range keys {
			jhBefore[i] = jh.GetNode(key)
		}

		start = time.Now()
		for i := initialNodes; i < initialNodes+addCount; i++ {
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

		t.Logf("Jump Hash: Adding %d nodes to %d initial nodes took %v, %d keys remapped (%.2f%%)",
			addCount, initialNodes, elapsed, changed, float64(changed)*100/float64(keyCount))

		// 测试Maglev哈希
		mhBefore := make([]string, keyCount)
		for i, key := range keys {
			mhBefore[i] = mh.GetNode(key)
		}

		start = time.Now()
		for i := initialNodes; i < initialNodes+addCount; i++ {
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

		t.Logf("Maglev Hash: Adding %d nodes to %d initial nodes took %v, %d keys remapped (%.2f%%)",
			addCount, initialNodes, elapsed, changed, float64(changed)*100/float64(keyCount))

		// 测试Rendezvous哈希
		rhBefore := make([]string, keyCount)
		for i, key := range keys {
			rhBefore[i] = rh.GetNode(key)
		}

		start = time.Now()
		for i := initialNodes; i < initialNodes+addCount; i++ {
			rh.AddNode(newNodes[i])
		}
		elapsed = time.Since(start)

		rhAfter := make([]string, keyCount)
		changed = 0
		for i, key := range keys {
			rhAfter[i] = rh.GetNode(key)
			if rhBefore[i] != rhAfter[i] {
				changed++
			}
		}

		t.Logf("Rendezvous Hash: Adding %d nodes to %d initial nodes took %v, %d keys remapped (%.2f%%)",
			addCount, initialNodes, elapsed, changed, float64(changed)*100/float64(keyCount))

		// 测试AnchorHash
		ahBefore := make([]string, keyCount)
		for i, key := range keys {
			ahBefore[i] = ah.GetNode(key)
		}

		start = time.Now()
		for i := initialNodes; i < initialNodes+addCount; i++ {
			ah.AddNode(newNodes[i])
		}
		elapsed = time.Since(start)

		ahAfter := make([]string, keyCount)
		changed = 0
		for i, key := range keys {
			ahAfter[i] = ah.GetNode(key)
			if ahBefore[i] != ahAfter[i] {
				changed++
			}
		}

		t.Logf("AnchorHash: Adding %d nodes to %d initial nodes took %v, %d keys remapped (%.2f%%)",
			addCount, initialNodes, elapsed, changed, float64(changed)*100/float64(keyCount))

		// 测试DxHash
		dxBefore := make([]string, keyCount)
		for i, key := range keys {
			dxBefore[i] = dx.GetNode(key)
		}

		start = time.Now()
		for i := initialNodes; i < initialNodes+addCount; i++ {
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

		// 测试Multi-Probe一致性哈希(k=5)
		mpch5Before := make([]string, keyCount)
		for i, key := range keys {
			mpch5Before[i] = mpch5.GetNode(key)
		}

		start = time.Now()
		for i := initialNodes; i < initialNodes+addCount; i++ {
			mpch5.AddNode(newNodes[i])
		}
		elapsed = time.Since(start)

		mpch5After := make([]string, keyCount)
		changed = 0
		for i, key := range keys {
			mpch5After[i] = mpch5.GetNode(key)
			if mpch5Before[i] != mpch5After[i] {
				changed++
			}
		}

		t.Logf("Multi-Probe CH(k=5): Adding %d nodes to %d initial nodes took %v, %d keys remapped (%.2f%%)",
			addCount, initialNodes, elapsed, changed, float64(changed)*100/float64(keyCount))

		// 测试Multi-Probe一致性哈希(k=21)
		mpch21Before := make([]string, keyCount)
		for i, key := range keys {
			mpch21Before[i] = mpch21.GetNode(key)
		}

		start = time.Now()
		for i := initialNodes; i < initialNodes+addCount; i++ {
			mpch21.AddNode(newNodes[i])
		}
		elapsed = time.Since(start)

		mpch21After := make([]string, keyCount)
		changed = 0
		for i, key := range keys {
			mpch21After[i] = mpch21.GetNode(key)
			if mpch21Before[i] != mpch21After[i] {
				changed++
			}
		}

		t.Logf("Multi-Probe CH(k=21): Adding %d nodes to %d initial nodes took %v, %d keys remapped (%.2f%%)",
			addCount, initialNodes, elapsed, changed, float64(changed)*100/float64(keyCount))
	}
}