package core

import (
    "math/big"
    "math/rand"
    cryptorand "crypto/rand"
)


func Rand(h RuntimeHandle, k func(*rand.Rand)(Object)) Observable {
    return Observable(func(pub DataPublisher) {
        randWorker <- func() {
            var r = rand.New(randSource { h })
            pub.AsyncReturn(k(r))
        }
    })
}
func RandBigInt(supremum *big.Int, h RuntimeHandle) Observable {
    if supremum.Sign() <= 0 {
        Crash(h, InvalidArgument, "non-positive random supremum")
    }
    return Observable(func(pub DataPublisher) {
        randWorker <- func() {
            pub.AsyncReturn(ObjIntFromBigInt(randBigInt(supremum, h)))
        }
    })
}

var randWorker = make(chan func(), 256)
var _ = (func() struct{} {
    go (func() {
        for k := range randWorker {
            k()
        }
    })()
    return struct{}{}
})()
type randSource struct {
    h RuntimeHandle
}
func (src randSource) Seed(_ int64) {
    panic("dummy method")
}
func (src randSource) Int63() int64 {
    var n big.Int
    return randBigInt(n.Lsh(big.NewInt(1), 63), src.h).Int64()
}
func (src randSource) Uint64() uint64 {
    var n big.Int
    return randBigInt(n.Lsh(big.NewInt(1), 64), src.h).Uint64()
}
func randBigInt(supremum *big.Int, h RuntimeHandle) *big.Int {
    if supremum.Sign() <= 0 {
        panic("invalid argument")
    }
    var n, err = cryptorand.Int(cryptorand.Reader, supremum)
    if err != nil {
        Crash(h, FailedToGenerateRandomNumber, err.Error())
    }
    return n
}


