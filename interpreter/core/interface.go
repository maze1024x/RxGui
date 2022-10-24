package core


type DispatchTable struct {
	Methods    [] *Function
	Children   [] *DispatchTable
	Parent     *DispatchTable
	Interface  string
}

func CallFirstMethod(I Interface, h RuntimeHandle) Object {
	var this = I.UnderlyingObject
	var f = *(I.DispatchTable.Methods[0])
	return f.Call([] Object { this }, nil, h)
}

func CraftSamInterface(o Object) Interface {
	var method = Function(NativeFunction(func(_ ([] Object), _ ([] Object), _ RuntimeHandle) Object {
		return o
	}))
	var table = &DispatchTable {
		Methods:  [] *Function { &method },
		Children: nil,
		Parent:   nil,
	}
	return Interface {
		UnderlyingObject: nil,
		DispatchTable:    table,
	}
}


