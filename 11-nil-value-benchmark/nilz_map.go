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

func Byte1() {
	m := make(map[int]byte, 100000)
	for i := 0; i < 100000; i++ {
		m[i] = 0
	}
}

func Uint() {
	m := make(map[int]uint8, 100000)
	for i := 0; i < 100000; i++ {
		m[i] = 0
	}
}

func Uint1() {
	m := make(map[int]uint8, 100000)
	for i := 0; i < 100000; i++ {
		m[i] = 1
	}
}

func Int() {
	m := make(map[int]int8, 100000)
	for i := 0; i < 100000; i++ {
		m[i] = 0
	}
}

func Int1() {
	m := make(map[int]int8, 100000)
	for i := 0; i < 100000; i++ {
		m[i] = 0
	}
}

func String() {
	m := make(map[int]string, 100000)
	for i := 0; i < 100000; i++ {
		m[i] = ""
	}
}

func String1() {
	m := make(map[int]string, 100000)
	for i := 0; i < 100000; i++ {
		m[i] = "1"
	}
}

func BoolF() {
	m := make(map[int]bool, 100000)
	for i := 0; i < 100000; i++ {
		m[i] = false
	}
}

func BoolT() {
	m := make(map[int]bool, 100000)
	for i := 0; i < 100000; i++ {
		m[i] = true
	}
}
