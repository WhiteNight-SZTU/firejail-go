package main

import (
	"runtime"
)

func busyLoop() {
	for {
		// 空循环，持续占用 CPU
	}
}

func main() {
	// 获取 CPU 核心数
	cpuCount := runtime.NumCPU()

	// 设置最大并发线程数与核心数一致
	runtime.GOMAXPROCS(cpuCount)

	// 为每个核心创建一个 Goroutine
	for i := 0; i < cpuCount; i++ {
		go busyLoop()
	}

	// 主程序保持运行
	select {}
}
