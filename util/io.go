package util

import (
    "io"
)

// A well-behaved substitution for fmt.Fscanln.
// The reader is recommended to be a buffered reader.
// Note that
//   1. trailing \r is ignored
//   2. ...[\n][EOF] and ...[EOF] are not distinguished.
//
func WellBehavedFscanln(f io.Reader) (string, int, error) {
    var buf = make([] byte, 0)
    var collect = func() string {
        var line = WellBehavedDecodeUtf8(buf)
        if len(line) > 0 && line[len(line)-1] == '\r' {
            line = line[:len(line)-1]
        }
        return line
    }
    var total = 0
    var one_byte_ [1] byte
    var one_byte = one_byte_[:]
    for {
        var n, err = f.Read(one_byte)
        total += n
        if err != nil {
            if err == io.EOF && len(buf) > 0 {
                return collect(), total, nil
            } else {
                return "", total, err
            }
        }
        if rune(one_byte[0]) != '\n' {
            buf = append(buf, one_byte[0])
        } else {
            return collect(), total, nil
        }
    }
}


