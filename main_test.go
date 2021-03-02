package main

import "testing"

func BenchmarkPerformanceRead(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = performanceRead("large_file.bin")
	}
}

func BenchmarkNormaleRead(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = normalRead("large_file.bin")
	}
}
