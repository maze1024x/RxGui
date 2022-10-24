package program

import (
	"unsafe"
	"rxgui/standalone/ctn"
	"rxgui/lang/source"
	"rxgui/lang/typsys"
)


type Binding struct {
	Name      string
	Type      typsys.CertainType
	Location  source.Location
	Constant  bool
}
func (b *Binding) PointerNumber() uintptr {
	return uintptr(unsafe.Pointer(b))
}
func BindingCompare(a *Binding, b *Binding) ctn.Ordering {
	if a == b {
		return ctn.Equal
	} else if a.PointerNumber() < b.PointerNumber() {
		return ctn.Smaller
	} else {
		return ctn.Bigger
	}
}


