package nilzmap

import "testing"

func Benchmark_Struct(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Struct()
	}
}

func Benchmark_Byte(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Byte()
	}
}

func Benchmark_Int(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Int()
	}
}

func Benchmark_Bool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Bool()
	}
}
