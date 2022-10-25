package qt

/*
#include <stdlib.h>
#include <stdint.h>
typedef void (*CgoCallbackCaller_t)(uint64_t);
void CgoCallbackCaller(uint64_t id) {
	void cgo_callback_caller_impl(uint64_t);
	cgo_callback_caller_impl(id);
}
*/
import "C"
import "rxgui/qt/cgohelper"

var cgo_callback_caller = C.CgoCallbackCaller_t(C.CgoCallbackCaller)

func cgo_callback(k func()) (C.uint64_t, func()) {
	var id, del = cgohelper.NewCallback(k)
	return C.uint64_t(id), del
}


