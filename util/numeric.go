package util

import (
    "math"
    "math/big"
    "math/cmplx"
    "strconv"
    "unicode"
    "unicode/utf8"
)


const MaxSafeIntegerToDouble = 9007199254740991
const MinSafeIntegerToDouble = -9007199254740991

func GetNumberUint(small uint) *big.Int {
    var n big.Int
    n.SetUint64(uint64(small))
    return &n
}

func GetNumberUint64(small uint64) *big.Int {
    var n big.Int
    n.SetUint64(small)
    return &n
}

func GetUintNumber(n *big.Int) uint {
    if n.Cmp(big.NewInt(0)) < 0 { panic("something went wrong") }
    var limit big.Int
    limit.SetUint64(uint64(^uint(0)))
    if n.Cmp(&limit) <= 0 {
        return uint(n.Uint64())
    } else {
        panic("given number too big")
    }
}

func GetInt64Integer(n *big.Int) int64 {
    if n.IsInt64() {
        return n.Int64()
    } else {
        panic("given number too big")
    }
}

func IsNonNegative(n *big.Int) bool {
    return (n.Cmp(big.NewInt(0)) >= 0)
}

func IsNormalFloat(x float64) bool {
    return !(math.IsNaN(x) || math.IsInf(x, 0))
}

func IsNormalComplex(z complex128) bool {
    return !(cmplx.IsNaN(z) || cmplx.IsInf(z))
}

func WellBehavedParseInteger(chars ([] rune)) (*big.Int, bool) {
    var abs_chars ([] rune)
    if chars[0] == '-' {
        abs_chars = chars[1:]
    } else {
        abs_chars = chars
    }
    var has_base_prefix = false
    if len(abs_chars) >= 2 {
        var c1 = abs_chars[0]
        var c2 = abs_chars[1]
        if c1 == '0' &&
        (c2 == 'x' || c2 == 'o' || c2 == 'b' ||
        c2 == 'X' || c2 == 'O' || c2 == 'B') {
            has_base_prefix = true
        }
    }
    var str = string(chars)
    // note: forbid "0" as octal prefix
    if has_base_prefix {
        return big.NewInt(0).SetString(str, 0)
    } else {
        if len(str) >= 2 && str[0] == '0' {
            return nil, false
        }
        return big.NewInt(0).SetString(str, 10)
    }
}

func DoubleToInteger(value float64) (*big.Int, bool) {
    if IsNormalFloat(value) {
        var n, _ = big.NewFloat(value).Int(nil)
        return n, true
    } else {
        return nil, false
    }
}

func IntegerToDouble(value *big.Int) (float64, bool) {
    if value.IsInt64() {
        var i64 = value.Int64()
        var ok = MinSafeIntegerToDouble <= i64 && i64 <= MaxSafeIntegerToDouble
        if ok {
            return float64(i64), true
        } else {
            return math.NaN(), false
        }
    } else {
        return math.NaN(), false
    }
}

func ParseDouble(chars ([] rune)) (float64, bool) {
    var value, err = strconv.ParseFloat(string(chars), 64)
    if err != nil { return math.NaN(), false }
    return value, true
}

func ParseRune(chars ([] rune)) (rune, bool) {
    var invalid = func() (rune, bool) {
        return utf8.RuneError, false
    }
    var got = func(r rune) (rune, bool) {
        if 0xD800 <= r && r <= 0xDFFF {
            return invalid()
        } else {
            return r, true
        }
    }
    if len(chars) == 0 {
        return invalid()
    } else if len(chars) == 1 {
        return got(chars[0])
    } else {
        var c0 = chars[0]
        if c0 == '`' {
            return got(chars[1])
        } else if c0 == '\\' {
            var c1 = chars[1]
            switch c1 {
            case 'n':
                return got('\n')
            case 'r':
                return got('\r')
            case 't':
                return got('\t')
            case 'e':
                return got('\033')
            case 'b':
                return got('\b')
            case 'a':
                return got('\a')
            case 'f':
                return got('\f')
            case 'u':
                var code_point_raw = string(chars[2:])
                var n, ok1 = big.NewInt(0).SetString(code_point_raw, 16)
                if !ok1 { return invalid() }
                var min = big.NewInt(0)
                var max = big.NewInt(unicode.MaxRune)
                var ok2 = ((min.Cmp(n) <= 0) && (n.Cmp(max) <= 0))
                if !ok2 { return invalid() }
                return got(rune(n.Int64()))
            default:
                return invalid()
            }
        } else {
            return invalid()
        }
    }
}


