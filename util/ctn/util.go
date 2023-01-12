package ctn

func assert(ok bool, msg string) {
    if !ok { panic(msg) }
}

func max(a uint64, b uint64) uint64 {
    if a >= b { return a } else { return b }
}

func zero[T any] () T {
    var t T
    return t
}


