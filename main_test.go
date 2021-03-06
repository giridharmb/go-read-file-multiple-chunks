package main

import "testing"

func BenchmarkPerformanceRead(b *testing.B) {
	for n := 0; n < b.N; n++ {
		performanceRead("large_file.bin", true)
	}
}

func BenchmarkNormaleRead(b *testing.B) {
	for n := 0; n < b.N; n++ {
		normalRead("large_file.bin", true)
	}
}
