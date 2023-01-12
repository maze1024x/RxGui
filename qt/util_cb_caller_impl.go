package qt

/*
#include <stdint.h>
*/
import "C"
import "rxgui/qt/cgohelper"

//export cgo_callback_caller_impl
func cgo_callback_caller_impl(id C.uint64_t) {
    var k = cgohelper.GetCallback(uint64(id))
    k()
}


