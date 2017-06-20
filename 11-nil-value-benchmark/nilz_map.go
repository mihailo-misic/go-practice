package nilzmap

func Struct() {
	m := make(map[int]struct{}, 100000)
	for i := 0; i < 100000; i++ {
		m[i] = struct{}{}
	}
}

func Byte() {
	m := make(map[int]byte, 100000)
	for i := 0; i < 100000; i++ {
		m[i] = 0
	}
}

func Int() {
	m := make(map[int]int8, 100000)
	for i := 0; i < 100000; i++ {
		m[i] = 0
	}
}

func Bool() {
	m := make(map[int]bool, 100000)
	for i := 0; i < 100000; i++ {
		m[i] = false
	}
}
