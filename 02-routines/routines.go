package _2_routines

import (
    "fmt"
    "time"
    "strconv"
)

func counting(name string) {
    for i := 0; i < 10; i++ {
        fmt.Println(name, strconv.Itoa(i))
        time.Sleep(time.Millisecond * 500)
    }
}

func main() {
    go counting("first")
    time.Sleep(time.Millisecond * 250)
    counting("second")
}
