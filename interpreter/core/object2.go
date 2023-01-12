package core

import (
    "time"
    "math"
    "regexp"
    "reflect"
    "math/big"
    "rxgui/qt"
    "rxgui/util/ctn"
    "rxgui/util/pseudounion"
)


var pseudoUnionTagReflectType = reflect.TypeOf(pseudounion.Tag(0))

func ToObject[T any] (v T) Object {
    return MakeObject(v, nil)
}
func MakeObject(v interface{}, h RuntimeHandle) Object {
    if ((v == nil) || (v == struct{}{})) {
        v = Object(nil)
    }
    if o, is_object := v.(Object); is_object {
        return o
    }
    var o = (func() ObjectImpl {
    switch V := v.(type) {
    // Primitive
        // Bool
            case bool:
                return Bool(V)
        // Int
            case *big.Int:
                return Int { V }
            case int:
                return Int { big.NewInt(int64(V)) }
        // Float
            case float64:
                return Float(V)
        // Bytes
            case [] byte:
                return Bytes(V)
        // String
            case string:
                return String(V)
        // Char
            case int32: // = rune
                return Char(V)
        // RegExp
            case *regexp.Regexp:
                return RegExp { V }
        // Time
            case time.Time:
                return Time(V)
        // File:
            case File:
                return V
        // Error
            case error:
                return Error { V }
    // Reflect
        // ReflectType
            case ReflectType:
                return V
        // ReflectValue
            case ReflectValue:
                return V
    // Rx
        // Observable
            case Observable:
                return V
        // Subject
            case Subject:
                return V
    // Interface and Lambda (Object)
        // Interface
            case Interface:
                return V
        // Lambda (Object)
            case Lambda:
                return V
    // Container (Object)
        // List (Object)
            case List:
                return V
        // Seq
            case Seq:
                return V
        // Queue
            case Queue:
                return V
        // Heap
            case Heap:
                return V
        // Set
            case Set:
                return V
        // Map
            case Map:
                return V
    // GUI
        // Action
            case Action:
                return V
        // Widget
            case Widget:
                return V
        // Signal
            case Signal:
                return V
        // Events
            case Events:
                return V
        // Prop
            case Prop:
                return V
    // Union (Object)
        case Union:
            return V
    // check for mistake
    case reflect.Value:
        panic("invalid argument")
    default:
        var rv = reflect.ValueOf(v)
        var rt = rv.Type()
        // Enum (int)
            if rt.Kind() == reflect.Int {
                var v = rv.Convert(reflect.TypeOf(int(0))).Interface().(int)
                return Enum(v)
            }
        // Union (ctn.Maybe)
            if _, ok := ctn.ReflectTypeMatchMaybe(rt); ok {
                if inner_rv, ok := ctn.ReflectMaybeValue(rv); ok {
                    var inner_v = inner_rv.Interface()
                    return Union {
                        Index:  OkIndex,
                        Object: MakeObject(inner_v, h),
                    }
                } else {
                    return Union { Index: NgIndex }
                }
            }
        // Record (ctn.Pair)
            if _, _, ok := ctn.ReflectTypeMatchPair(rt); ok {
                var a_rv, b_rv = ctn.ReflectPairUnpack(rv)
                var a_v = a_rv.Interface()
                var b_v = b_rv.Interface()
                var a = MakeObject(a_v, h)
                var b = MakeObject(b_v, h)
                return Record {
                    Objects: [] Object { a, b },
                }
            }
        // Union (pseudo-union)
            if rt.Kind() == reflect.Struct {
                if rt.NumField() >= 1 {
                if rt.Field(0).Type == pseudoUnionTagReflectType {
                    var index = int(rv.Field(0).Int())
                    var object_rv = rv.Field(index)
                    var object_v = object_rv.Interface()
                    var object = MakeObject(object_v, h)
                    return Union {
                        Index:  (index - 1),
                        Object: object,
                    }
                }}
        // Record (struct)
                var n = rt.NumField()
                var objects = make([] Object, n)
                for i := 0; i < n; i += 1 {
                    var field_rv = rv.Field(i)
                    var field_v = field_rv.Interface()
                    objects[i] = MakeObject(field_v, h)
                }
                return Record { objects }
            }
    // Container (generic)
        // List (generic)
            if rt.Kind() == reflect.Slice {
                var n = rv.Len()
                var nodes = make([] ListNode, n)
                for i := 0; i < n; i += 1 {
                    var item_rv = rv.Index(i)
                    var item_v = item_rv.Interface()
                    nodes[i].Value = MakeObject(item_v, h)
                }
                return NodesToList(nodes)
            }
    // Interface and Lambda (generic)
        // Lambda (generic)
            if isLambdaFuncReflectType(rt) {
                var get_in_rv func(Object)([] reflect.Value)
                if rt.NumIn() == 0 {
                    get_in_rv = func(_ Object) ([] reflect.Value) {
                        return nil
                    }
                } else if rt.NumIn() == 1 {
                    get_in_rv = func(arg Object) ([] reflect.Value) {
                        var arg_ptr, arg_rv = reflectNew(rt.In(0))
                        ConvertObject(arg, arg_ptr, h)
                        return [] reflect.Value { arg_rv }
                    }
                } else if rt.NumIn() == 2 {
                    get_in_rv = func(arg Object) ([] reflect.Value) {
                        var a_ptr, a_rv = reflectNew(rt.In(0))
                        var b_ptr, b_rv = reflectNew(rt.In(1))
                        var r = (*arg).(Record) // ctn.Pair
                        ConvertObject(r.Objects[0], a_ptr, h)
                        ConvertObject(r.Objects[1], b_ptr, h)
                        return [] reflect.Value { a_rv, b_rv }
                    }
                } else {
                    panic("something went wrong")
                }
                var call = func(arg Object) Object {
                    var in_rv = get_in_rv(arg)
                    var out_rv = rv.Call(in_rv)
                    var out = out_rv[0].Interface()
                    return MakeObject(out, h)
                }
                var l = Lambda { call }
                if rt.Name() == "" {
                    return l
                } else {
                    if rt.NumIn() == 0 {
                        return CraftSamInterface(l.Call(nil))
                    } else {
                        return CraftSamInterface(Obj(l))
                    }
                }
            }
        panic("unsupported value type: " + rt.String())
    } })()
    return &o
}

func FromObject[T any] (o Object) T {
    var t T
    var p = &t
    ConvertObject(o, p, nil)
    return t
}
func ConvertObject(o Object, p interface{}, h RuntimeHandle) {
    if o == nil {
        return
    }
    if _, is_empty_struct := p.(*struct{}); is_empty_struct {
        panic("cannot assign non-nil object to struct{}")
    }
    if o_ptr, is_object_ptr := p.(*Object); is_object_ptr {
        *o_ptr = o
        return
    }
    { var o = *o
    switch P := p.(type) {
    // Primitive
        // Bool
            case *(bool):
                *P = bool(o.(Bool))
        // Int
            case *(*big.Int):
                *P = o.(Int).Value
            case *(int):
                *P = clampTo32Int(o.(Int).Value)
        // Float
            case *(float64):
                *P = float64(o.(Float))
        // Bytes
            case *([] byte):
                *P = o.(Bytes)
        // String
            case *(string):
                *P = string(o.(String))
        // Char
            case *(int32): // = rune
                *P = int32(o.(Char))
        // RegExp
            case *(*regexp.Regexp):
                *P = o.(RegExp).Value
        // Time
            case *(time.Time):
                *P = time.Time(o.(Time))
        // File
            case *(File):
                *P = o.(File)
        // Error
            case *(error):
                *P = o.(Error).Value
    // Reflect
        // ReflectType
            case *(ReflectType):
                *P = o.(ReflectType)
        // ReflectValue
            case *(ReflectValue):
                *P = o.(ReflectValue)
    // Rx
        // Observable
            case *(Observable):
                *P = o.(Observable)
        // Subject
            case *(Subject):
                *P = o.(Subject)
    // Interface and Lambda (Object)
        // Interface
            case *(Interface):
                *P = o.(Interface)
        // Lambda (Object)
            case *(Lambda):
                *P = o.(Lambda)
    // Container (Object)
        // List (Object)
            case *(List):
                *P = o.(List)
        // Seq
            case *(Seq):
                *P = o.(Seq)
        // Queue
            case *(Queue):
                *P = o.(Queue)
        // Heap
            case *(Heap):
                *P = o.(Heap)
        // Set
            case *(Set):
                *P = o.(Set)
        // Map
            case *(Map):
                *P = o.(Map)
    // GUI
        // Action
            case *(Action):
                *P = o.(Action)
        // Widget
            case *(Widget):
                *P = o.(Widget)
        // Signal
            case *(Signal):
                *P = o.(Signal)
        // Events
            case *(Events):
                *P = o.(Events)
        // Prop
            case *(Prop):
                *P = o.(Prop)
    // Union (Object)
        case *(Union):
            *P = o.(Union)
    default:
        var ptr_rv = reflect.ValueOf(p)
        // check for mistake
        if ptr_rv.Kind() != reflect.Ptr {
            panic("invalid argument")
        }
        var rv = ptr_rv.Elem()
        var rt = rv.Type()
        // Enum (int)
            if rt.Kind() == reflect.Int {
                var v = int(o.(Enum))
                rv.Set(reflect.ValueOf(v).Convert(rt))
                return
            }
        // Union (ctn.Maybe)
            if inner_rt, ok := ctn.ReflectTypeMatchMaybe(rt); ok {
                var u = o.(Union)
                if u.Index == OkIndex {
                    var inner_ptr, inner_rv = reflectNew(inner_rt)
                    ConvertObject(u.Object, inner_ptr, h)
                    rv.Set(ctn.ReflectJust(inner_rv))
                } else if u.Index == NgIndex {
                    rv.Set(ctn.ReflectNothing(inner_rt))
                } else {
                    panic("something went wrong")
                }
                return
            }
        // Record (ctn.Pair)
            if a_rt, b_rt, ok := ctn.ReflectTypeMatchPair(rt); ok {
                var r = o.(Record)
                var a_ptr, a_rv = reflectNew(a_rt)
                var b_ptr, b_rv = reflectNew(b_rt)
                ConvertObject(r.Objects[0], a_ptr, h)
                ConvertObject(r.Objects[1], b_ptr, h)
                rv.Set(ctn.ReflectMakePair(a_rv, b_rv))
                return
            }
        // Union (pseudo-union)
            if rt.Kind() == reflect.Struct {
                if rt.NumField() >= 1 {
                if rt.Field(0).Type == pseudoUnionTagReflectType {
                    var u = o.(Union)
                    var i = (u.Index + 1)
                    var object_ptr_rv = rv.Field(i).Addr()
                    var object_ptr = object_ptr_rv.Interface()
                    rv.Field(0).SetInt(int64(i))
                    ConvertObject(u.Object, object_ptr, h)
                    return
                }}
        // Record (struct)
                var r = o.(Record)
                var n = len(r.Objects)
                for i := 0; i < n; i += 1 {
                    var field_ptr_rv = rv.Field(i).Addr()
                    var field_ptr = field_ptr_rv.Interface()
                    ConvertObject(r.Objects[i], field_ptr, h)
                }
                return
            }
    // Container (generic)
        // List (generic)
            if rt.Kind() == reflect.Slice {
                var l = o.(List)
                var slice_rv = reflect.MakeSlice(rt, 0, 0)
                l.ForEach(func(el_obj Object) {
                    var el_ptr_rv = reflect.New(rt.Elem())
                    var el_ptr = el_ptr_rv.Interface()
                    ConvertObject(el_obj, el_ptr, h)
                    var el_rv = el_ptr_rv.Elem()
                    slice_rv = reflect.Append(slice_rv, el_rv)
                })
                rv.Set(slice_rv)
                return
            }
    // Interface and Lambda (generic)
        // Lambda (generic)
            if isLambdaFuncReflectType(rt) {
                var l = (func() Lambda {
                    if rt.Name() == "" {
                        return o.(Lambda)
                    } else {
                        if h == nil {
                            panic("SAM conversion without RuntimeHandle")
                        }
                        var I = o.(Interface)
                        if rt.NumIn() == 0 {
                            return getSamInterfaceValueAsLambda(I, h)
                        } else {
                            return getSamInterfaceLambda(I, h)
                        }
                    }
                })()
                var get_arg_obj func([] reflect.Value) Object
                if rt.NumIn() == 0 {
                    get_arg_obj = func(_ ([] reflect.Value)) Object {
                        return nil
                    }
                } else if rt.NumIn() == 1 {
                    get_arg_obj = func(in ([] reflect.Value)) Object {
                        var in0_v = in[0].Interface()
                        return MakeObject(in0_v, h)
                    }
                } else if rt.NumIn() == 2 {
                    get_arg_obj = func(in ([] reflect.Value)) Object {
                        var in0_v = in[0].Interface()
                        var in1_v = in[1].Interface()
                        var a = MakeObject(in0_v, h)
                        var b = MakeObject(in1_v, h)
                        var pair = [] Object { a, b }
                        var o = ObjectImpl(Record { pair })
                        return &o
                    }
                } else {
                    panic("something went wrong")
                }
                var out_rt = rt.Out(0)
                var f_rv = reflect.MakeFunc(rt, func(in ([] reflect.Value)) ([] reflect.Value) {
                    var arg = get_arg_obj(in)
                    var ret = l.Call(arg)
                    var out_ptr, out_rv = reflectNew(out_rt)
                    ConvertObject(ret, out_ptr, h)
                    var out = [] reflect.Value { out_rv }
                    return out
                })
                rv.Set(f_rv)
                return
            }
        panic("unsupported pointer type: " + rt.String())
    } }
}

func MakeNativeFunction(v interface{}) NativeFunction {
    var rv = reflect.ValueOf(v)
    var rt = rv.Type()
    if !((rt.Kind() == reflect.Func) && !(rt.IsVariadic())) {
        panic("invalid argument")
    }
    return NativeFunction(func(args ([] Object), ctx ([] Object), h RuntimeHandle) Object {
        var num_in = rt.NumIn()
        var num_out = rt.NumOut()
        var in = make([] reflect.Value, num_in)
        for i := 0; i < num_in; i += 1 {
            var in_t = rt.In(i)
            if i < len(args) {
                var arg = args[i]
                var arg_ptr, arg_rv = reflectNew(in_t)
                ConvertObject(arg, arg_ptr, h)
                in[i] = arg_rv
            } else {
                var j = (i - len(args))
                if j < len(ctx) {
                    var item = ctx[j]
                    var item_ptr, item_rv = reflectNew(in_t)
                    ConvertObject(item, item_ptr, h)
                    in[i] = item_rv
                } else {
                    in[i] = reflect.ValueOf(h)
                }
            }
        }
        var out = rv.Call(in)
        if num_out == 0 {
            return nil
        } else if num_out == 1 {
            var out_v = out[0].Interface()
            return MakeObject(out_v, h)
        } else {
            panic("invalid argument")
        }
    })
}

func retrieveObject[T any] (o Observable, h RuntimeHandle, k func(T)(Observable)) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        pub.run(o, ctx, awaitNoexceptObserver(ob, h, func(obj Object) {
            pub.run(k(FromObject[T](obj)), ctx, ob)
        }))
    })
}
func retrieveObjectInChildContext[T any] (parent *context, o Observable, h RuntimeHandle, k func(T,*context,func())(Observable)) Observable {
    return Observable(func(pub DataPublisher) {
        var _, ob = pub.useInheritedContext()
        var ctx, dispose = parent.createChild()
        pub.run(o, ctx, awaitNoexceptObserver(ob, h, func(obj Object) {
            pub.run(k(FromObject[T](obj), ctx, dispose), ctx, ob)
        }))
    })
}

func doSync(k func()) Observable {
    return Observable(func(pub DataPublisher) {
        pub.SyncReturn(func() (Object, error) {
            k()
            return nil, nil
        })
    })
}
func doSync1[T any] (k func()(T)) Observable {
    return Observable(func(pub DataPublisher) {
        pub.SyncReturn(func() (Object, error) {
            var t = k()
            return ToObject(t), nil
        })
    })
}
func doSync2[T any] (k func()(T,func())) Observable {
    return Observable(func(pub DataPublisher) {
        pub.SyncReturn(func() (Object, error) {
            var t, c = k()
            pub.context.registerCleaner(c)
            return ToObject(t), nil
        })
    })
}
func onSync(k func()(func(qt.Pkg,func()))) Observable {
    return Observable(func(pub DataPublisher) {
        var pkg, dispose = qt.CreatePkg()
        pub.context.registerCleaner(dispose)
        k()(pkg, func() {
            pub.observer.value(nil)
        })
    })
}
func onSync1[T any] (k func()(func(qt.Pkg,func(T)))) Observable {
    return Observable(func(pub DataPublisher) {
        var pkg, dispose = qt.CreatePkg()
        pub.context.registerCleaner(dispose)
        k()(pkg, func(v T) {
            pub.observer.value(ToObject(v))
        })
    })
}

func reflectNew(t reflect.Type) (interface{}, reflect.Value) {
    var ptr_rv = reflect.New(t)
    var ptr = ptr_rv.Interface()
    var rv = ptr_rv.Elem()
    return ptr, rv
}
func isLambdaFuncReflectType(t reflect.Type) bool {
    return (t.Kind() == reflect.Func) &&
        !(t.IsVariadic()) &&
        (t.NumIn() == 0 || t.NumIn() == 1 || t.NumIn() == 2) &&
        (t.NumOut() == 1)
}
func getSamInterfaceLambda(I Interface, h RuntimeHandle) Lambda {
    if len(I.DispatchTable.Methods) == 1 && len(I.DispatchTable.Children) == 0 {
        return (*(CallFirstMethod(I, h))).(Lambda)
    } else {
        panic("expect SAM interface but got non-SAM interface")
    }
}
func getSamInterfaceValueAsLambda(I Interface, h RuntimeHandle) Lambda {
    if len(I.DispatchTable.Methods) == 1 && len(I.DispatchTable.Children) == 0 {
        var value = CallFirstMethod(I, h)
        return Lambda { func(_ Object) Object {
            return value
        }}
    } else {
        panic("expect SAM interface but got non-SAM interface")
    }
}

var maxInt32 = big.NewInt(math.MaxInt32)
var minInt32 = big.NewInt(math.MinInt32)
func clampTo32Int(n *big.Int) int {
    if n.Cmp(maxInt32) > 0 {
        return math.MaxInt32
    } else if n.Cmp(minInt32) < 0 {
        return math.MinInt32
    } else {
        return int(n.Int64())
    }
}

func ObjInt(n int) Object {
    return ObjIntFromInt64(int64(n))
}
func ObjIntFromBigInt(n *big.Int) Object {
    return Obj(Int { n })
}
func ObjIntFromInt64(n int64) Object {
    return Obj(Int { big.NewInt(n) })
}
func ObjBool(p bool) Object {
    return Obj(Bool(p))
}
func ObjFloat(x float64) Object {
    return Obj(Float(x))
}
func ObjString(s string) Object {
    return Obj(String(s))
}
func ObjBytes(b ([] byte)) Object {
    return Obj(Bytes(b))
}
func ObjPair(a Object, b Object) Object {
    return Obj(Record { Objects: [] Object { a, b } })
}
func ObjList(l ([] Object)) Object {
    var buf ListBuilder
    for _, item := range l {
        buf.Append(item)
    }
    return Obj(buf.Collect())
}
func ObjQueue(q ctn.Queue[Object]) Object {
    return Obj(Queue(q))
}
func ObjTime(t time.Time) Object {
    return Obj(Time(t))
}
func ObjTimeNow() Object {
    return ObjTime(time.Now())
}
func ObjFile(path string) Object {
    return Obj(File { path })
}
func GetBool(o Object) bool {
    return bool((*o).(Bool))
}
func GetInt(o Object) int {
    return clampTo32Int((*o).(Int).Value)
}
func GetIntAsRawBigInt(o Object) *big.Int {
    return (*o).(Int).Value
}
func GetFloat(o Object) float64 {
    return float64((*o).(Float))
}
func GetChar(o Object) rune {
    return int32((*o).(Char))
}
func GetString(o Object) string {
    return string((*o).(String))
}
func GetBytes(o Object) ([] byte) {
    return ([] byte)((*o).(Bytes))
}
func GetError(o Object) error {
    return (*o).(Error).Value
}
func GetPair(o Object) (Object, Object) {
    return FromObject[ctn.Pair[Object,Object]](o)()
}
func GetList(o Object) List {
    return (*o).(List)
}
func GetSeq(o Object) Seq {
    return (*o).(Seq)
}
func GetQueue(o Object) ctn.Queue[Object] {
    return ctn.Queue[Object]((*o).(Queue))
}
func GetHeap(o Object) ctn.Heap[Object] {
    return ctn.Heap[Object]((*o).(Heap))
}
func GetSet(o Object) ctn.Set[Object] {
    return ctn.Set[Object]((*o).(Set))
}
func GetMap(o Object) ctn.Map[Object,Object] {
    return ctn.Map[Object,Object]((*o).(Map))
}
func GetTime(o Object) time.Time {
    return time.Time((*o).(Time))
}
func GetFile(o Object) File {
    return (*o).(File)
}
func GetObservable(o Object) Observable {
    return (*o).(Observable)
}
func GetWidget(o Object) Widget {
    return (*o).(Widget)
}
func YieldObservables(o ...Observable) func(func(Observable)) {
    return func(yield func(Observable)) {
        for _, o := range o {
            yield(o)
        }
    }
}
func ListToPair(o Object) Object {
    var l = GetList(o)
    a, l, ok1 := l.Shifted()
    b, l, ok2 := l.Shifted()
    if !(ok1 && ok2) { panic("invalid argument") }
    return ObjPair(a, b)
}
func QueueToPair(o Object) Object {
    var q = GetQueue(o)
    a, q, ok1 := q.Shifted()
    b, q, ok2 := q.Shifted()
    if !(ok1 && ok2) { panic("invalid argument") }
    return ObjPair(a, b)
}


