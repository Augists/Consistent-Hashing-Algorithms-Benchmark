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
	hr40 := algorithms.NewHashRing(40) // 添加40个虚拟节点的测试
	jh := algorithms.NewJumpHash()
	mh := algorithms.NewMaglevHash(65537)
	mh2039 := algorithms.NewMaglevHash(2039) // 添加2039表长的测试
	rh := algorithms.NewRendezvousHash()
	ah := algorithms.NewAnchorHash(2000)
	dx := algorithms.NewDxHashWithParams(nodeCount) // 使用预设节点数
	mpch5 := algorithms.NewMPCH(5)                  // 使用5个探针
	mpch21 := algorithms.NewMPCH(21)                // 使用21个探针

	// 添加节点
	for i := 0; i < nodeCount; i++ {
		node := fmt.Sprintf("node_%d", i)
		mod.AddNode(node)
		hr.AddNode(node)
		hr40.AddNode(node)
		jh.AddNode(node)
		mh.AddNode(node)
		mh2039.AddNode(node)
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

	// 测试哈希环分布(160个虚拟节点)
	hrDistribution := make(map[string]int)
	for _, key := range keys {
		node := hr.GetNode(key)
		hrDistribution[node]++
	}

	// 测试哈希环分布(40个虚拟节点)
	hr40Distribution := make(map[string]int)
	for _, key := range keys {
		node := hr40.GetNode(key)
		hr40Distribution[node]++
	}

	// 测试跳跃哈希分布
	jhDistribution := make(map[string]int)
	for _, key := range keys {
		node := jh.GetNode(key)
		jhDistribution[node]++
	}

	// 测试Maglev哈希分布(65537表长)
	mhDistribution := make(map[string]int)
	for _, key := range keys {
		node := mh.GetNode(key)
		mhDistribution[node]++
	}

	// 测试Maglev哈希分布(2039表长)
	mh2039Distribution := make(map[string]int)
	for _, key := range keys {
		node := mh2039.GetNode(key)
		mh2039Distribution[node]++
	}

	// 测试Rendezvous哈希分布
	rhDistribution := make(map[string]int)
	for _, key := range keys {
		node := rh.GetNode(key)
		rhDistribution[node]++
	}

	// 测试AnchorHash分布(1000长度)
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
	hr40StdDev := calculateStdDev(hr40Distribution, avg, nodeCount)
	jhStdDev := calculateStdDev(jhDistribution, avg, nodeCount)
	mhStdDev := calculateStdDev(mhDistribution, avg, nodeCount)
	mh2039StdDev := calculateStdDev(mh2039Distribution, avg, nodeCount)
	rhStdDev := calculateStdDev(rhDistribution, avg, nodeCount)
	ahStdDev := calculateStdDev(ahDistribution, avg, nodeCount)
	dxStdDev := calculateStdDev(dxDistribution, avg, nodeCount)
	mpch5StdDev := calculateStdDev(mpch5Distribution, avg, nodeCount)
	mpch21StdDev := calculateStdDev(mpch21Distribution, avg, nodeCount)

	fmt.Println("=== 分布均匀性测试 ===")
	fmt.Printf("使用 %d 个节点和 %d 个键进行测试\n", nodeCount, keyCount)
	fmt.Printf("每个节点平均分配键数: %d\n", avg)
	fmt.Printf("  直接哈希取模标准差: %.2f\n", modStdDev)
	fmt.Printf("  哈希环(160个虚拟节点)标准差: %.2f\n", hrStdDev)
	fmt.Printf("  哈希环(40个虚拟节点)标准差: %.2f\n", hr40StdDev)
	fmt.Printf("  跳跃哈希标准差: %.2f\n", jhStdDev)
	fmt.Printf("  Maglev哈希(65537表长)标准差: %.2f\n", mhStdDev)
	fmt.Printf("  Maglev哈希(2039表长)标准差: %.2f\n", mh2039StdDev)
	fmt.Printf("  Rendezvous哈希标准差: %.2f\n", rhStdDev)
	fmt.Printf("  AnchorHash(2000长度)标准差: %.2f\n", ahStdDev)
	fmt.Printf("  DxHash标准差: %.2f\n", dxStdDev)
	fmt.Printf("  Multi-Probe一致性哈希(k=5)标准差: %.2f\n", mpch5StdDev)
	fmt.Printf("  Multi-Probe一致性哈希(k=21)标准差: %.2f\n", mpch21StdDev)
}

func testNodeAdditionRemapping() {
	fmt.Println("\n2. 添加节点时的重映射测试:")

	initialNodes := 1000
	addCount := 10
	keyCount := 100000

	// 创建各种算法实例
	mod := algorithms.NewModHash()
	hr := algorithms.NewHashRing(160)
	hr40 := algorithms.NewHashRing(40) // 添加40个虚拟节点的测试
	jh := algorithms.NewJumpHash()
	mh := algorithms.NewMaglevHash(65537)
	mh2039 := algorithms.NewMaglevHash(2039) // 添加2039表长的测试
	rh := algorithms.NewRendezvousHash()
	ah := algorithms.NewAnchorHash(2000)
	dx := algorithms.NewDxHashWithParams(initialNodes + addCount) // 使用预设节点数
	mpch5 := algorithms.NewMPCH(5)                                // 使用5个探针
	mpch21 := algorithms.NewMPCH(21)                              // 使用21个探针

	// 添加初始节点
	for i := 0; i < initialNodes; i++ {
		node := fmt.Sprintf("node_%d", i)
		mod.AddNode(node)
		hr.AddNode(node)
		hr40.AddNode(node)
		jh.AddNode(node)
		mh.AddNode(node)
		mh2039.AddNode(node)
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

	// 测试哈希环(160个虚拟节点)
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

	// 测试哈希环(40个虚拟节点)
	hr40Before := make([]string, keyCount)
	for i, key := range keys {
		hr40Before[i] = hr40.GetNode(key)
	}

	start = time.Now()
	for i := initialNodes; i < initialNodes+addCount; i++ {
		hr40.AddNode(newNodes[i])
	}
	hr40Elapsed := time.Since(start)

	hr40After := make([]string, keyCount)
	hr40Changed := 0
	for i, key := range keys {
		hr40After[i] = hr40.GetNode(key)
		if hr40Before[i] != hr40After[i] {
			hr40Changed++
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

	// 测试Maglev哈希(65537表长)
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

	// 测试Maglev哈希(2039表长)
	mh2039Before := make([]string, keyCount)
	for i, key := range keys {
		mh2039Before[i] = mh2039.GetNode(key)
	}

	start = time.Now()
	for i := initialNodes; i < initialNodes+addCount; i++ {
		mh2039.AddNode(newNodes[i])
	}
	mh2039Elapsed := time.Since(start)

	mh2039After := make([]string, keyCount)
	mh2039Changed := 0
	for i, key := range keys {
		mh2039After[i] = mh2039.GetNode(key)
		if mh2039Before[i] != mh2039After[i] {
			mh2039Changed++
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

	// 测试AnchorHash(2000长度)
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

	fmt.Println("=== 添加节点时的重映射测试 ===")
	fmt.Printf("添加 %d 个节点到 %d 个初始节点:\n", addCount, initialNodes)
	fmt.Printf("  直接哈希取模: 耗时 %v, 重映射键数 %d (%.2f%%)\n", modElapsed, modChanged, float64(modChanged)*100/float64(keyCount))
	fmt.Printf("  哈希环(160个虚拟节点): 耗时 %v, 重映射键数 %d (%.2f%%)\n", hrElapsed, hrChanged, float64(hrChanged)*100/float64(keyCount))
	fmt.Printf("  哈希环(40个虚拟节点): 耗时 %v, 重映射键数 %d (%.2f%%)\n", hr40Elapsed, hr40Changed, float64(hr40Changed)*100/float64(keyCount))
	fmt.Printf("  跳跃哈希: 耗时 %v, 重映射键数 %d (%.2f%%)\n", jhElapsed, jhChanged, float64(jhChanged)*100/float64(keyCount))
	fmt.Printf("  Maglev哈希(65537表长): 耗时 %v, 重映射键数 %d (%.2f%%)\n", mhElapsed, mhChanged, float64(mhChanged)*100/float64(keyCount))
	fmt.Printf("  Maglev哈希(2039表长): 耗时 %v, 重映射键数 %d (%.2f%%)\n", mh2039Elapsed, mh2039Changed, float64(mh2039Changed)*100/float64(keyCount))
	fmt.Printf("  Rendezvous哈希: 耗时 %v, 重映射键数 %d (%.2f%%)\n", rhElapsed, rhChanged, float64(rhChanged)*100/float64(keyCount))
	fmt.Printf("  AnchorHash(2000长度): 耗时 %v, 重映射键数 %d (%.2f%%)\n", ahElapsed, ahChanged, float64(ahChanged)*100/float64(keyCount))
	fmt.Printf("  DxHash: 耗时 %v, 重映射键数 %d (%.2f%%)\n", dxElapsed, dxChanged, float64(dxChanged)*100/float64(keyCount))
	fmt.Printf("  Multi-Probe CH(k=5): 耗时 %v, 重映射键数 %d (%.2f%%)\n", mpch5Elapsed, mpch5Changed, float64(mpch5Changed)*100/float64(keyCount))
	fmt.Printf("  Multi-Probe CH(k=21): 耗时 %v, 重映射键数 %d (%.2f%%)\n", mpch21Elapsed, mpch21Changed, float64(mpch21Changed)*100/float64(keyCount))
}

func testPerformance() {
	fmt.Println("\n3. 查询性能测试:")

	nodeCount := 1000
	opCount := 100000

	// 创建各种算法实例
	mod := algorithms.NewModHash()
	hr := algorithms.NewHashRing(160)
	hr40 := algorithms.NewHashRing(40) // 添加40个虚拟节点的测试
	jh := algorithms.NewJumpHash()
	mh := algorithms.NewMaglevHash(65537)
	mh2039 := algorithms.NewMaglevHash(2039) // 添加2039表长的测试
	rh := algorithms.NewRendezvousHash()
	ah := algorithms.NewAnchorHash(2000)
	dx := algorithms.NewDxHashWithParams(nodeCount) // 使用预设节点数
	mpch5 := algorithms.NewMPCH(5)                  // 使用5个探针
	mpch21 := algorithms.NewMPCH(21)                // 使用21个探针

	// 添加节点
	for i := 0; i < nodeCount; i++ {
		node := fmt.Sprintf("node_%d", i)
		mod.AddNode(node)
		hr.AddNode(node)
		hr40.AddNode(node)
		jh.AddNode(node)
		mh.AddNode(node)
		mh2039.AddNode(node)
		rh.AddNode(node)
		ah.AddNode(node)
		dx.AddNode(node)
		mpch5.AddNode(node)
		mpch21.AddNode(node)
	}

	// 生成测试键
	keys := make([]string, opCount)
	for i := 0; i < opCount; i++ {
		keys[i] = fmt.Sprintf("key_%d", i)
	}

	// 测试直接哈希取模性能
	start := time.Now()
	for _, key := range keys {
		mod.GetNode(key)
	}
	modElapsed := time.Since(start)

	// 测试哈希环性能(160个虚拟节点)
	start = time.Now()
	for _, key := range keys {
		hr.GetNode(key)
	}
	hrElapsed := time.Since(start)

	// 测试哈希环性能(40个虚拟节点)
	start = time.Now()
	for _, key := range keys {
		hr40.GetNode(key)
	}
	hr40Elapsed := time.Since(start)

	// 测试跳跃哈希性能
	start = time.Now()
	for _, key := range keys {
		jh.GetNode(key)
	}
	jhElapsed := time.Since(start)

	// 测试Maglev哈希性能(65537表长)
	start = time.Now()
	for _, key := range keys {
		mh.GetNode(key)
	}
	mhElapsed := time.Since(start)

	// 测试Maglev哈希性能(2039表长)
	start = time.Now()
	for _, key := range keys {
		mh2039.GetNode(key)
	}
	mh2039Elapsed := time.Since(start)

	// 测试Rendezvous哈希性能
	start = time.Now()
	for _, key := range keys {
		rh.GetNode(key)
	}
	rhElapsed := time.Since(start)

	// 测试AnchorHash性能(2000长度)
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

	fmt.Println("=== 查询性能测试 ===")
	fmt.Printf("执行 %d 次查询操作:\n", opCount)
	fmt.Printf("  直接哈希取模: %v\n", modElapsed)
	fmt.Printf("  哈希环(160个虚拟节点): %v\n", hrElapsed)
	fmt.Printf("  哈希环(40个虚拟节点): %v\n", hr40Elapsed)
	fmt.Printf("  跳跃哈希: %v\n", jhElapsed)
	fmt.Printf("  Maglev哈希(65537表长): %v\n", mhElapsed)
	fmt.Printf("  Maglev哈希(2039表长): %v\n", mh2039Elapsed)
	fmt.Printf("  Rendezvous哈希: %v\n", rhElapsed)
	fmt.Printf("  AnchorHash(2000长度): %v\n", ahElapsed)
	fmt.Printf("  DxHash: %v\n", dxElapsed)
	fmt.Printf("  Multi-Probe一致性哈希(k=5): %v\n", mpch5Elapsed)
	fmt.Printf("  Multi-Probe一致性哈希(k=21): %v\n", mpch21Elapsed)
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
