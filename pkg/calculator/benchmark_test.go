package calculator

import (
	"fmt"
	"testing"
)

// Input sizes for benchmarks
var benchmarkSizes = []int{1, 250, 501, 1000, 5000, 10000, 50000}

// Standard pack sizes for benchmarks
var standardPackSizes = []int{250, 500, 1000, 2000, 5000}

// BenchmarkCalculatePacks benchmarks the original CalculatePacks function
func BenchmarkCalculatePacks(b *testing.B) {
	for _, size := range benchmarkSizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = CalculatePacks(size, standardPackSizes)
			}
		})
	}
}

// BenchmarkCalculatePacksOptimized benchmarks the optimized version
func BenchmarkCalculatePacksOptimized(b *testing.B) {
	for _, size := range benchmarkSizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = CalculatePacksOptimized(size, standardPackSizes)
			}
		})
	}
}

// BenchmarkCalculatePacksMemory benchmarks memory usage of the original function
func BenchmarkCalculatePacksMemory(b *testing.B) {
	for _, size := range benchmarkSizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = CalculatePacks(size, standardPackSizes)
			}
		})
	}
}

// BenchmarkCalculatePacksOptimizedMemory benchmarks memory usage of the optimized function
func BenchmarkCalculatePacksOptimizedMemory(b *testing.B) {
	for _, size := range benchmarkSizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = CalculatePacksOptimized(size, standardPackSizes)
			}
		})
	}
}

// BenchmarkCalculatePacksLarge benchmarks the original function with large inputs
func BenchmarkCalculatePacksLarge(b *testing.B) {
	var largeSize int = 100000
	b.Run(fmt.Sprintf("Size_%d", largeSize), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = CalculatePacks(largeSize, standardPackSizes)
		}
	})
}

// BenchmarkCalculatePacksOptimizedLarge benchmarks the optimized function with large inputs
func BenchmarkCalculatePacksOptimizedLarge(b *testing.B) {
	var largeSize int = 100000
	b.Run(fmt.Sprintf("Size_%d", largeSize), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = CalculatePacksOptimized(largeSize, standardPackSizes)
		}
	})
}

// BenchmarkCalculatePacksCustom benchmarks the original function with custom pack sizes
func BenchmarkCalculatePacksCustom(b *testing.B) {
	customPackSizes := []int{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000}
	for _, size := range []int{501, 5000, 10000} {
		b.Run(fmt.Sprintf("Size_%d_CustomPacks", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = CalculatePacks(size, customPackSizes)
			}
		})
	}
}

// BenchmarkCalculatePacksOptimizedCustom benchmarks the optimized function with custom pack sizes
func BenchmarkCalculatePacksOptimizedCustom(b *testing.B) {
	customPackSizes := []int{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000}
	for _, size := range []int{501, 5000, 10000} {
		b.Run(fmt.Sprintf("Size_%d_CustomPacks", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = CalculatePacksOptimized(size, customPackSizes)
			}
		})
	}
}
