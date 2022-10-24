package source

import (
	"fmt"
	"rxgui/standalone/ctn"
)


type Ref struct {
	Namespace  string
	ItemName   string
}
func MakeRef(ns string, item string) Ref {
	return Ref {
		Namespace: ns,
		ItemName:  item,
	}
}
func (r Ref) String() string {
	if r.Namespace == "" {
		return r.ItemName
	} else {
		return fmt.Sprintf("%s::%s", r.Namespace, r.ItemName)
	}
}
func RefCompare(a Ref, b Ref) ctn.Ordering {
	var o = ctn.StringCompare(a.Namespace, b.Namespace)
	if o != ctn.Equal {
		return o
	} else {
		return ctn.StringCompare(a.ItemName, b.ItemName)
	}
}


