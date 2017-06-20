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

func Benchmark_Byte1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Byte1()
	}
}

func Benchmark_Uint(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Uint()
	}
}

func Benchmark_Uint1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Uint1()
	}
}

func Benchmark_Int(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Int()
	}
}

func Benchmark_Int1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Int1()
	}
}

func Benchmark_String(b *testing.B) {
	for n := 0; n < b.N; n++ {
		String()
	}
}

func Benchmark_String1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		String1()
	}
}

func Benchmark_BoolF(b *testing.B) {
	for n := 0; n < b.N; n++ {
		BoolF()
	}
}

func Benchmark_BoolT(b *testing.B) {
	for n := 0; n < b.N; n++ {
		BoolT()
	}
}
