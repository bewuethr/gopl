package ch11ex06

import (
	"testing"

	"github.com/bewuethr/gopl/chapter02/ch02ex04"
	"github.com/bewuethr/gopl/chapter02/ch02ex05"
)

func benchmarkPopCountTableLookup(b *testing.B, n int) {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			PopCount(uint64(j))
		}
	}
}

func BenchmarkPopCountTableLookup1(b *testing.B)     { benchmarkPopCountTableLookup(b, 1) }
func BenchmarkPopCountTableLookup10(b *testing.B)    { benchmarkPopCountTableLookup(b, 10) }
func BenchmarkPopCountTableLookup100(b *testing.B)   { benchmarkPopCountTableLookup(b, 100) }
func BenchmarkPopCountTableLookup1000(b *testing.B)  { benchmarkPopCountTableLookup(b, 1000) }
func BenchmarkPopCountTableLookup10000(b *testing.B) { benchmarkPopCountTableLookup(b, 10000) }

func benchmarkPopCountBitShift(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			ch02ex04.PopCount(uint64(j))
		}
	}
}

func BenchmarkPopCountBitShift1(b *testing.B)     { benchmarkPopCountBitShift(b, 1) }
func BenchmarkPopCountBitShift10(b *testing.B)    { benchmarkPopCountBitShift(b, 10) }
func BenchmarkPopCountBitShift100(b *testing.B)   { benchmarkPopCountBitShift(b, 100) }
func BenchmarkPopCountBitShift1000(b *testing.B)  { benchmarkPopCountBitShift(b, 1000) }
func BenchmarkPopCountBitShift10000(b *testing.B) { benchmarkPopCountBitShift(b, 10000) }

func benchmarkPopCountBitClear(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			ch02ex05.PopCount(uint64(j))
		}
	}
}

func BenchmarkPopCountBitClear1(b *testing.B)     { benchmarkPopCountBitClear(b, 1) }
func BenchmarkPopCountBitClear10(b *testing.B)    { benchmarkPopCountBitClear(b, 10) }
func BenchmarkPopCountBitClear100(b *testing.B)   { benchmarkPopCountBitClear(b, 100) }
func BenchmarkPopCountBitClear1000(b *testing.B)  { benchmarkPopCountBitClear(b, 1000) }
func BenchmarkPopCountBitClear10000(b *testing.B) { benchmarkPopCountBitClear(b, 10000) }
