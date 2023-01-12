package qt

/*
#include <stdlib.h>
#include <stdint.h>
*/
import "C"
import "unsafe"


type Pkg struct { *pkgImpl }
type pkgImpl struct {
    items  [] func()
}
func CreatePkg() (Pkg, func()) {
    var pkg = Pkg { &pkgImpl { make([] func(), 0) } }
    return pkg, pkg.dispose
}
func (pkg Pkg) push(item func()) {
    pkg.items = append(pkg.items, item)
}
func (pkg Pkg) pop() (func(), bool) {
    if len(pkg.items) > 0 {
        var last = (len(pkg.items) - 1)
        var last_item = pkg.items[last]
        pkg.items[last] = nil
        pkg.items = pkg.items[:last]
        return last_item, true
    } else {
        return nil, false
    }
}
func (pkg Pkg) dispose() {
    if pkg.items == nil {
        return
    }
    if len(pkg.items) == 0 {
        // in case we forget `defer` keyword
        panic("something went wrong")
    }
    for item, ok := pkg.pop(); ok; item, ok = pkg.pop() {
        item()
    }
    pkg.items = nil
}

func str(s string, ctx Pkg) *C.char {
    var ptr = C.CString(s)
    var del = func() { C.free(unsafe.Pointer(ptr)) }
    ctx.push(del)
    return ptr
}

func addrlen(buf ([] byte)) (*C.uint8_t, C.size_t) {
    if len(buf) == 0 {
        return nil, C.size_t(0)
    }
    return (*C.uint8_t)(unsafe.Pointer(&(buf[0]))), C.size_t(uint(len(buf)))
}

func ptrlen(widgets ([] Widget)) (*unsafe.Pointer, C.size_t) {
    if len(widgets) == 0 {
        return nil, C.size_t(0)
    }
    var ptr = make([] unsafe.Pointer, len(widgets))
    for i := range widgets {
        ptr[i] = widgets[i].ptr
    }
    return &(ptr[0]), C.size_t(len(ptr))
}


