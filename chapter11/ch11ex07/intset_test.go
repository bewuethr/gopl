package ch11ex07

import (
	"math/rand"
	"testing"
	"time"

	intset "github.com/bewuethr/gopl/chapter11/ch11ex02"
)

func getRNG() *rand.Rand {
	seed := time.Now().UTC().UnixNano()
	return rand.New(rand.NewSource(seed))
}

func BenchmarkAdd(b *testing.B) {
	rng := getRNG()
	is := &IntSet{}
	for i := 0; i < b.N; i++ {
		is.Add(rng.Int() & 0xffff)
	}
}

func benchmarkUnionWith(b *testing.B, n int) {
	rng := getRNG()
	is1, is2 := &IntSet{}, &IntSet{}
	is1.Add(rng.Int() & n)
	is2.Add(rng.Int() & n)
	for i := 0; i < b.N; i++ {
		is1.UnionWith(is2)
	}
}

func BenchmarkUnionWith1(b *testing.B) { benchmarkUnionWith(b, 0xf) }
func BenchmarkUnionWith2(b *testing.B) { benchmarkUnionWith(b, 0xff) }
func BenchmarkUnionWith4(b *testing.B) { benchmarkUnionWith(b, 0xffff) }
func BenchmarkUnionWith8(b *testing.B) { benchmarkUnionWith(b, 0xffffffff) }

func Benchmark32BitAdd(b *testing.B) {
	rng := getRNG()
	is := &IntSet32{}
	for i := 0; i < b.N; i++ {
		is.Add(rng.Int() & 0xffff)
	}
}

func benchmark32BitUnionWith(b *testing.B, n int) {
	rng := getRNG()
	is1, is2 := &IntSet32{}, &IntSet32{}
	is1.Add(rng.Int() & n)
	is2.Add(rng.Int() & n)
	for i := 0; i < b.N; i++ {
		is1.UnionWith(is2)
	}
}

func Benchmark32BitUnionWith1(b *testing.B) { benchmark32BitUnionWith(b, 0xf) }
func Benchmark32BitUnionWith2(b *testing.B) { benchmark32BitUnionWith(b, 0xff) }
func Benchmark32BitUnionWith4(b *testing.B) { benchmark32BitUnionWith(b, 0xffff) }
func Benchmark32BitUnionWith8(b *testing.B) { benchmark32BitUnionWith(b, 0xffffffff) }

func BenchmarkRefSetAdd(b *testing.B) {
	rng := getRNG()
	rs := make(intset.RefSet)
	for i := 0; i < b.N; i++ {
		rs[rng.Int()&0xffff] = true
	}
}

func benchmarkRefSetUnionWith(b *testing.B, n int) {
	rng := getRNG()
	rs1, rs2 := make(intset.RefSet), make(intset.RefSet)
	rs1[rng.Int()&n] = true
	rs2[rng.Int()&n] = true
	for i := 0; i < b.N; i++ {
		rs1.UnionWith(rs2)
	}
}

func BenchmarkRefSetUnionWith1(b *testing.B) { benchmarkRefSetUnionWith(b, 0xf) }
func BenchmarkRefSetUnionWith2(b *testing.B) { benchmarkRefSetUnionWith(b, 0xff) }
func BenchmarkRefSetUnionWith4(b *testing.B) { benchmarkRefSetUnionWith(b, 0xffff) }
func BenchmarkRefSetUnionWith8(b *testing.B) { benchmarkRefSetUnionWith(b, 0xffffffff) }
