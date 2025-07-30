package main

import (
	"fmt"
	"math"
	"time"

	"consistent_hashing/algorithms"
)

func main() {
	fmt.Println("一致性哈希算法对比测试")
	fmt.Println("========================")

	// 测试分布均匀性
	testDistribution()

	// 测试添加节点时的重映射
	testNodeAdditionRemapping()

	// 测试查询性能
	testPerformance()
}

func testDistribution() {
	fmt.Println("\n1. 分布均匀性测试:")

	nodeCount := 100
	keyCount := 100000

	// 创建各种算法实例
	mod := algorithms.NewModHash()
	hr := algorithms.NewHashRing(160)
	jh := algorithms.NewJumpHash()
	mh := algorithms.NewMaglevHash(65537)
	rh := algorithms.NewRendezvousHash()
	ah := algorithms.NewAnchorHash(1000)
	dx := algorithms.NewDxHashWithParams(nodeCount) // 使用预设节点数
	mpch5 := algorithms.NewMPCH(5)   // 使用5个探针
	mpch21 := algorithms.NewMPCH(21) // 使用21个探针

	// 添加节点
	for i := 0; i < nodeCount; i++ {
		node := fmt.Sprintf("node_%d", i)
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

	fmt.Printf("  节点数: %d, 键数: %d\n", nodeCount, keyCount)
	fmt.Printf("  平均分配数: %d\n", avg)
	fmt.Printf("  直接哈希取模标准差: %.2f\n", modStdDev)
	fmt.Printf("  哈希环标准差: %.2f\n", hrStdDev)
	fmt.Printf("  跳跃哈希标准差: %.2f\n", jhStdDev)
	fmt.Printf("  Maglev哈希标准差: %.2f\n", mhStdDev)
	fmt.Printf("  Rendezvous哈希标准差: %.2f\n", rhStdDev)
	fmt.Printf("  AnchorHash标准差: %.2f\n", ahStdDev)
	fmt.Printf("  DxHash标准差: %.2f\n", dxStdDev)
	fmt.Printf("  Multi-Probe CH(k=5)标准差: %.2f\n", mpch5StdDev)
	fmt.Printf("  Multi-Probe CH(k=21)标准差: %.2f\n", mpch21StdDev)
}

func testNodeAdditionRemapping() {
	fmt.Println("\n2. 添加节点时的重映射测试:")

	initialNodes := 1000
	addCount := 10
	keyCount := 100000

	// 创建各种算法实例
	mod := algorithms.NewModHash()
	hr := algorithms.NewHashRing(160)
	jh := algorithms.NewJumpHash()
	mh := algorithms.NewMaglevHash(65537)
	rh := algorithms.NewRendezvousHash()
	ah := algorithms.NewAnchorHash(2000)
	dx := algorithms.NewDxHashWithParams(initialNodes+addCount) // 使用预设节点数
	mpch5 := algorithms.NewMPCH(5)   // 使用5个探针
	mpch21 := algorithms.NewMPCH(21) // 使用21个探针

	// 添加初始节点
	for i := 0; i < initialNodes; i++ {
		node := fmt.Sprintf("node_%d", i)
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

	// 测试直接哈希取模
	modBefore := make([]string, keyCount)
	for i, key := range keys {
		modBefore[i] = mod.GetNode(key)
	}

	newNodes := make([]string, initialNodes+addCount)
	for i := 0; i < initialNodes+addCount; i++ {
		newNodes[i] = fmt.Sprintf("node_%d", i)
	}

	start := time.Now()
	for i := initialNodes; i < initialNodes+addCount; i++ {
		mod.AddNode(newNodes[i])
	}
	modElapsed := time.Since(start)

	modAfter := make([]string, keyCount)
	modChanged := 0
	for i, key := range keys {
		modAfter[i] = mod.GetNode(key)
		if modBefore[i] != modAfter[i] {
			modChanged++
		}
	}

	// 测试哈希环
	hrBefore := make([]string, keyCount)
	for i, key := range keys {
		hrBefore[i] = hr.GetNode(key)
	}

	start = time.Now()
	for i := initialNodes; i < initialNodes+addCount; i++ {
		hr.AddNode(newNodes[i])
	}
	hrElapsed := time.Since(start)

	hrAfter := make([]string, keyCount)
	hrChanged := 0
	for i, key := range keys {
		hrAfter[i] = hr.GetNode(key)
		if hrBefore[i] != hrAfter[i] {
			hrChanged++
		}
	}

	// 测试跳跃哈希
	jhBefore := make([]string, keyCount)
	for i, key := range keys {
		jhBefore[i] = jh.GetNode(key)
	}

	start = time.Now()
	for i := initialNodes; i < initialNodes+addCount; i++ {
		jh.AddNode(newNodes[i])
	}
	jhElapsed := time.Since(start)

	jhAfter := make([]string, keyCount)
	jhChanged := 0
	for i, key := range keys {
		jhAfter[i] = jh.GetNode(key)
		if jhBefore[i] != jhAfter[i] {
			jhChanged++
		}
	}

	// 测试Maglev哈希
	mhBefore := make([]string, keyCount)
	for i, key := range keys {
		mhBefore[i] = mh.GetNode(key)
	}

	start = time.Now()
	for i := initialNodes; i < initialNodes+addCount; i++ {
		mh.AddNode(newNodes[i])
	}
	mhElapsed := time.Since(start)

	mhAfter := make([]string, keyCount)
	mhChanged := 0
	for i, key := range keys {
		mhAfter[i] = mh.GetNode(key)
		if mhBefore[i] != mhAfter[i] {
			mhChanged++
		}
	}

	// 测试Rendezvous哈希
	rhBefore := make([]string, keyCount)
	for i, key := range keys {
		rhBefore[i] = rh.GetNode(key)
	}

	start = time.Now()
	for i := initialNodes; i < initialNodes+addCount; i++ {
		rh.AddNode(newNodes[i])
	}
	rhElapsed := time.Since(start)

	rhAfter := make([]string, keyCount)
	rhChanged := 0
	for i, key := range keys {
		rhAfter[i] = rh.GetNode(key)
		if rhBefore[i] != rhAfter[i] {
			rhChanged++
		}
	}

	// 测试AnchorHash
	ahBefore := make([]string, keyCount)
	for i, key := range keys {
		ahBefore[i] = ah.GetNode(key)
	}

	start = time.Now()
	for i := initialNodes; i < initialNodes+addCount; i++ {
		ah.AddNode(newNodes[i])
	}
	ahElapsed := time.Since(start)

	ahAfter := make([]string, keyCount)
	ahChanged := 0
	for i, key := range keys {
		ahAfter[i] = ah.GetNode(key)
		if ahBefore[i] != ahAfter[i] {
			ahChanged++
		}
	}

	// 测试DxHash
	dxBefore := make([]string, keyCount)
	for i, key := range keys {
		dxBefore[i] = dx.GetNode(key)
	}

	start = time.Now()
	for i := initialNodes; i < initialNodes+addCount; i++ {
		dx.AddNode(newNodes[i])
	}
	dxElapsed := time.Since(start)

	dxAfter := make([]string, keyCount)
	dxChanged := 0
	for i, key := range keys {
		dxAfter[i] = dx.GetNode(key)
		if dxBefore[i] != dxAfter[i] {
			dxChanged++
		}
	}

	// 测试Multi-Probe一致性哈希(k=5)
	mpch5Before := make([]string, keyCount)
	for i, key := range keys {
		mpch5Before[i] = mpch5.GetNode(key)
	}

	start = time.Now()
	for i := initialNodes; i < initialNodes+addCount; i++ {
		mpch5.AddNode(newNodes[i])
	}
	mpch5Elapsed := time.Since(start)

	mpch5After := make([]string, keyCount)
	mpch5Changed := 0
	for i, key := range keys {
		mpch5After[i] = mpch5.GetNode(key)
		if mpch5Before[i] != mpch5After[i] {
			mpch5Changed++
		}
	}

	// 测试Multi-Probe一致性哈希(k=21)
	mpch21Before := make([]string, keyCount)
	for i, key := range keys {
		mpch21Before[i] = mpch21.GetNode(key)
	}

	start = time.Now()
	for i := initialNodes; i < initialNodes+addCount; i++ {
		mpch21.AddNode(newNodes[i])
	}
	mpch21Elapsed := time.Since(start)

	mpch21After := make([]string, keyCount)
	mpch21Changed := 0
	for i, key := range keys {
		mpch21After[i] = mpch21.GetNode(key)
		if mpch21Before[i] != mpch21After[i] {
			mpch21Changed++
		}
	}

	fmt.Printf("  初始节点数: %d, 新增节点数: %d, 测试键数: %d\n", initialNodes, addCount, keyCount)
	fmt.Printf("  直接哈希取模: 耗时 %v, 重映射键数 %d (%.2f%%)\n", modElapsed, modChanged, float64(modChanged)*100/float64(keyCount))
	fmt.Printf("  哈希环: 耗时 %v, 重映射键数 %d (%.2f%%)\n", hrElapsed, hrChanged, float64(hrChanged)*100/float64(keyCount))
	fmt.Printf("  跳跃哈希: 耗时 %v, 重映射键数 %d (%.2f%%)\n", jhElapsed, jhChanged, float64(jhChanged)*100/float64(keyCount))
	fmt.Printf("  Maglev哈希: 耗时 %v, 重映射键数 %d (%.2f%%)\n", mhElapsed, mhChanged, float64(mhChanged)*100/float64(keyCount))
	fmt.Printf("  Rendezvous哈希: 耗时 %v, 重映射键数 %d (%.2f%%)\n", rhElapsed, rhChanged, float64(rhChanged)*100/float64(keyCount))
	fmt.Printf("  AnchorHash: 耗时 %v, 重映射键数 %d (%.2f%%)\n", ahElapsed, ahChanged, float64(ahChanged)*100/float64(keyCount))
	fmt.Printf("  DxHash: 耗时 %v, 重映射键数 %d (%.2f%%)\n", dxElapsed, dxChanged, float64(dxChanged)*100/float64(keyCount))
	fmt.Printf("  Multi-Probe CH(k=5): 耗时 %v, 重映射键数 %d (%.2f%%)\n", mpch5Elapsed, mpch5Changed, float64(mpch5Changed)*100/float64(keyCount))
	fmt.Printf("  Multi-Probe CH(k=21): 耗时 %v, 重映射键数 %d (%.2f%%)\n", mpch21Elapsed, mpch21Changed, float64(mpch21Changed)*100/float64(keyCount))
}

func testPerformance() {
	fmt.Println("\n3. 查询性能测试:")

	nodeCount := 1000
	testKeys := 100000

	// 创建各种算法实例
	mod := algorithms.NewModHash()
	hr := algorithms.NewHashRing(160)
	jh := algorithms.NewJumpHash()
	mh := algorithms.NewMaglevHash(65537)
	rh := algorithms.NewRendezvousHash()
	ah := algorithms.NewAnchorHash(2000)
	dx := algorithms.NewDxHashWithParams(nodeCount) // 使用预设节点数
	mpch5 := algorithms.NewMPCH(5)   // 使用5个探针
	mpch21 := algorithms.NewMPCH(21) // 使用21个探针

	// 添加节点
	for i := 0; i < nodeCount; i++ {
		node := fmt.Sprintf("node_%d", i)
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
	keys := make([]string, testKeys)
	for i := 0; i < testKeys; i++ {
		keys[i] = fmt.Sprintf("key_%d", i)
	}

	// 测试直接哈希取模性能
	start := time.Now()
	for _, key := range keys {
		mod.GetNode(key)
	}
	modElapsed := time.Since(start)

	// 测试哈希环性能
	start = time.Now()
	for _, key := range keys {
		hr.GetNode(key)
	}
	hrElapsed := time.Since(start)

	// 测试跳跃哈希性能
	start = time.Now()
	for _, key := range keys {
		jh.GetNode(key)
	}
	jhElapsed := time.Since(start)

	// 测试Maglev哈希性能
	start = time.Now()
	for _, key := range keys {
		mh.GetNode(key)
	}
	mhElapsed := time.Since(start)

	// 测试Rendezvous哈希性能
	start = time.Now()
	for _, key := range keys {
		rh.GetNode(key)
	}
	rhElapsed := time.Since(start)

	// 测试AnchorHash性能
	start = time.Now()
	for _, key := range keys {
		ah.GetNode(key)
	}
	ahElapsed := time.Since(start)

	// 测试DxHash性能
	start = time.Now()
	for _, key := range keys {
		dx.GetNode(key)
	}
	dxElapsed := time.Since(start)

	// 测试Multi-Probe一致性哈希(k=5)性能
	start = time.Now()
	for _, key := range keys {
		mpch5.GetNode(key)
	}
	mpch5Elapsed := time.Since(start)

	// 测试Multi-Probe一致性哈希(k=21)性能
	start = time.Now()
	for _, key := range keys {
		mpch21.GetNode(key)
	}
	mpch21Elapsed := time.Since(start)

	fmt.Printf("  节点数: %d, 测试键数: %d\n", nodeCount, testKeys)
	fmt.Printf("  直接哈希取模: %v\n", modElapsed)
	fmt.Printf("  哈希环: %v\n", hrElapsed)
	fmt.Printf("  跳跃哈希: %v\n", jhElapsed)
	fmt.Printf("  Maglev哈希: %v\n", mhElapsed)
	fmt.Printf("  Rendezvous哈希: %v\n", rhElapsed)
	fmt.Printf("  AnchorHash: %v\n", ahElapsed)
	fmt.Printf("  DxHash: %v\n", dxElapsed)
	fmt.Printf("  Multi-Probe CH(k=5): %v\n", mpch5Elapsed)
	fmt.Printf("  Multi-Probe CH(k=21): %v\n", mpch21Elapsed)
	fmt.Println("\n测试完成!")
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
