package core


type Function interface {
    Call(args ([] Object), ctx ([] Object), h RuntimeHandle) Object
}

type NativeFunction func(args ([] Object), ctx ([] Object), h RuntimeHandle) Object
func (f NativeFunction) Call(args ([] Object), ctx ([] Object), h RuntimeHandle) Object {
    return f(args, ctx, h)
}

type FieldValueGetterFunction struct { Index int }
func (f FieldValueGetterFunction) Call(args ([] Object), _ ([] Object), _ RuntimeHandle) Object {
    var arg = args[0]
    var record = (*arg).(Record)
    return record.Objects[f.Index]
}

func FunctionToLambda(op Function, unpack bool, ctx ([] Object), h RuntimeHandle) Lambda {
    if unpack {
        return Lambda {
            Call: func(arg Object) Object {
                var args = (*arg).(Record).Objects
                return op.Call(args, ctx, h)
            },
        }
    } else {
        return Lambda {
            Call: func(arg Object) Object {
                var args = [] Object { arg }
                return op.Call(args, ctx, h)
            },
        }
    }
}
func FunctionToLambdaObject(op Function, unpack bool, ctx ([] Object), h RuntimeHandle) Object {
    var o = ObjectImpl(FunctionToLambda(op, unpack, ctx, h))
    return &o
}


