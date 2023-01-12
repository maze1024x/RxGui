package ctn


func MapEach[A any, B any] (sa ([] A), f func(A)(B)) ([] B) {
    var sb = make([] B, len(sa))
    for i := range sa {
        sb[i] = f(sa[i])
    }
    return sb
}

func MapEachDeflate[A any, B any] (sa ([] A), f func(A)(B,bool)) ([] B) {
    var sb = make([] B, 0)
    for _, a := range sa {
        if b, ok := f(a); ok {
            sb = append(sb, b)
        }
    }
    return sb
}

func Filter[T any] (s ([] T), p func(T)(bool)) ([] T) {
    var result = make([] T, 0)
    for _, t := range s {
        if p(t) {
            result = append(result, t)
        }
    }
    return result
}

func RemoveFrom[T comparable] (s ([] T), target T) ([] T) {
    var result = make([] T, 0)
    for _, t := range s {
        if t != target {
            result = append(result, t)
        }
    }
    return result
}

func Reduce[A any, B any] (sa ([] A), b0 B, f func(b B, a A)(B)) B {
    var b = b0
    for _, a := range sa {
        b = f(b, a)
    }
    return b
}

func Reverse[T any] (s ([] T)) ([] T) {
    var L = len(s)
    var result = make([] T, L)
    for i := range s {
        result[((L-1)-i)] = s[i]
    }
    return result
}

func StableSorted[T any] (s ([] T), lt Less[T]) ([] T, error) {
    if len(s) == 0 {
        return nil, nil
    } else if len(s) == 1 {
        return [] T { s[0] }, nil
    } else {
        var middle = (len(s) / 2)
        var u, _ = StableSorted(s[:middle], lt)
        var v, _ = StableSorted(s[middle:], lt)
        return mergeSorted(u, v, lt), nil
    }
}
func mergeSorted[T any] (s1 ([] T), s2 ([] T), lt Less[T]) ([] T) {
    var p1 = 0
    var p2 = 0
    var l1 = len(s1)
    var l2 = len(s2)
    var r = make([] T, (l1 + l2))
    var q = 0
    for {
        if p1 < l1 && p2 < l2 {
            var v1 = s1[p1]
            var v2 = s2[p2]
            if !(lt(v2, v1)) {
                r[q] = v1; p1++; q++
            } else {
                r[q] = v2; p2++; q++
            }
        } else if p1 < l1 {
            var v1 = s1[p1]
            r[q] = v1; p1++; q++
        } else if p2 < l2 {
            var v2 = s2[p2]
            r[q] = v2; p2++; q++
        } else {
            return r
        }
    }
}


