package core

import ( gctx "context"; "rxgui/util/ctn" )


type Observable func(DataPublisher)

type DataPublisher struct {
    eventloop  *EventLoop
    context    *context
    observer   *observer
}
type DataSubscriber struct {
    Values     chan <- Object
    Error      chan <- error
    Terminate  chan <- bool
}
type observer struct {
    value     func(Object)
    error     func(error)
    complete  func()
}
func run(run Observable, eventloop *EventLoop, ctx *context, ob *observer) {
    if ctx.isDisposed() {
        return
    }
    var terminated = false
    var terminate = func() { terminated = true }
    var active = func() bool { return !(terminated || ctx.isDisposed()) }
    var proxy = &observer {
        value:    func(v Object) { if active() { ob.value(v) } },
        error:    func(e error)  { if active() { terminate(); ob.error(e) } },
        complete: func()         { if active() { terminate(); ob.complete() } },
    }
    run(DataPublisher {
        eventloop: eventloop,
        context:   ctx,
        observer:  proxy,
    })
}
func Run(o Observable, eventloop *EventLoop, sub DataSubscriber) {
    var V, E, T = sub.Values, sub.Error, sub.Terminate
    eventloop.addTask(func() {
        var ob = &observer {
            value: func(v Object) {
                if V != nil { V <- v }
            },
            error: func(e error) {
                if E != nil { E <- e; close(E) }
                if T != nil { T <- false }
            },
            complete: func() {
                if V != nil { close(V) }
                if T != nil { T <- true }
            },
        }
        run(o, eventloop, nil, ob)
    })
}
func (pub DataPublisher) useInheritedContext() (*context, *observer) {
    return pub.context, pub.observer
}
func (pub DataPublisher) useNewChildContext() (*context, func(), *observer) {
    var ctx, dispose = pub.context.createChild()
    return ctx, dispose, pub.observer
}
func (pub DataPublisher) run(o Observable, ctx *context, ob *observer) {
    run(o, pub.eventloop, ctx, ob)
}
func (pub DataPublisher) addTask(k func()) {
    pub.eventloop.addTask(k)
}
func (pub DataPublisher) addTimer(ms int, n int, ctx *context, k func()) {
    pub.eventloop.addTimer(ms, n, ctx, k)
}
func (pub DataPublisher) SyncReturn(k func()(Object,error)) {
    var v, e = k()
    if e == nil {
        pub.observer.value(v)
        pub.observer.complete()
    } else {
        pub.observer.error(e)
    }
}
func (pub DataPublisher) SyncGenerate(k func(yield func(Object))(error)) {
    var e = k(pub.observer.value)
    if e == nil {
        pub.observer.complete()
    } else {
        pub.observer.error(e)
    }
}
func (pub DataPublisher) AsyncContext() gctx.Context {
    return pub.context.goContext()
}
func (pub DataPublisher) AsyncThrow(e error) {
    pub.eventloop.addTask(func() {
        pub.observer.error(e)
    })
}
func (pub DataPublisher) AsyncReturn(v Object) {
    pub.eventloop.addTask(func() {
        pub.observer.value(v)
        pub.observer.complete()
    })
}
func (pub DataPublisher) AsyncGenerate() (func(Object),func()) {
    var yield = func(v Object) {
        pub.eventloop.addTask(func() {
            pub.observer.value(v)
        })
    }
    var complete = func() {
        pub.eventloop.addTask(func() {
            pub.observer.complete()
        })
    }
    return yield, complete
}
func ObservableSyncValue(v Object) Observable {
    return Observable(func(pub DataPublisher) {
        pub.observer.value(v)
        pub.observer.complete()
    })
}
func ObservableFlattenLast(o Observable) Observable {
    return o.Await(func(obj Object) Observable {
        return GetObservable(obj)
    })
}

type context struct {
    parent     *context
    children   [] *context
    disposed   bool
    cleaners   [] cleaner
    number     uint64
    g_context  gctx.Context
    g_dispose  func()
}
type cleaner struct {
    clean   func()
    number  uint64
}
func (ctx *context) createChild() (*context, func()) {
    var g_ctx, g_dispose = gctx.WithCancel(gctx.Background())
    var parent = ctx
    var child = &context {
        parent:    parent,
        children:  make([] *context, 0),
        disposed:  false,
        cleaners:  make([] cleaner, 0),
        number:    getNumber(),
        g_context: g_ctx,
        g_dispose: g_dispose,
    }
    if parent != nil {
        if parent.disposed { panic("something went wrong") }
        parent.children = append(parent.children, child)
    }
    return child, child.__dispose
}
func (ctx *context) __dispose() {
    var parent = ctx.parent
    var self = ctx
    if !(self.disposed) {
        if parent != nil {
            self.parent = nil
            parent.children = ctn.RemoveFrom(parent.children, self)
        }
        var contexts = make([] *context, 0)
        var cleaners = make([] cleaner, 0)
        var q = [] *context { self }
        for len(q) > 0 {
            var c = q[0]; q = q[1:]
            contexts = append(contexts, c)
            cleaners = append(cleaners, c.cleaners...)
            q = append(q, c.children...)
        }
        contexts, _ = ctn.StableSorted(contexts, contextSorter)
        cleaners, _ = ctn.StableSorted(cleaners, cleanerSorter)
        for _, c := range contexts {
            c.disposed = true
        }
        for _, c := range contexts {
            c.g_dispose()
        }
        for _, c := range cleaners {
            c.clean()
        }
    }
}
func contextSorter(a *context, b *context) bool { return (a.number > b.number) }
func cleanerSorter(a cleaner, b cleaner) bool { return (a.number > b.number) }
var numberCounter = uint64(0)
func getNumber() uint64 { var n = numberCounter; numberCounter++; return n }
// registerCleaner registers release operations of persistent resources
func (ctx *context) registerCleaner(c func()) {
    if c == nil {
        panic("invalid argument")
    }
    if ctx != nil {
        if !(ctx.disposed) {
            ctx.cleaners = append(ctx.cleaners, cleaner {
                clean:  c,
                number: getNumber(),
            })
        }
    }
}
func (ctx *context) isDisposed() bool {
    if ctx != nil {
        return ctx.disposed
    } else {
        return false
    }
}
func (ctx *context) goContext() gctx.Context {
    if ctx != nil {
        return ctx.g_context
    } else {
        return gctx.Background()
    }
}
func WithChildContext(o Observable) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, dispose, ob = pub.useNewChildContext()
        pub.run(o, ctx, &observer {
            value: ob.value,
            error: func(err error) {
                dispose()
                ob.error(err)
            },
            complete: func() {
                dispose()
                ob.complete()
            },
        })
    })
}
type CancelError struct {}
func (CancelError) Error() string { return "cancelled" }
func WithCancelTrigger(sig Observable, o Observable) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, dispose, ob = pub.useNewChildContext()
        pub.run(o, ctx, &observer {
            value: ob.value,
            error: func(e error) {
                dispose()
                ob.error(e)
            },
            complete: func() {
                dispose()
                ob.complete()
            },
        })
        pub.run(sig, ctx, &observer {
            value: func(_ Object) {
                dispose()
                ob.error(CancelError {})
            },
            error: func(e error) {
                dispose()
                ob.error(e)
            },
            complete: func() {
                dispose()
                ob.error(CancelError {})
            },
        })
    })
}
func WithCancelTimeout(ms int, o Observable) Observable {
    return WithCancelTrigger(SetTimeout(ms), o)
}

type Subject struct { *subjectImpl }
type subjectImpl struct {
    observerNextId  uint64
    observerList    [] *observer
    observerIndex   map[uint64] uint
    terminated      bool
    maybeError      error
    notifyingFlag   bool
    recentValues    [] Object
    runtimeHandle   RuntimeHandle
}
func CreateSubject(h RuntimeHandle, replay int, items ...Object) Subject {
    if replay < 0 { replay = 0 }
    var b = Subject { &subjectImpl {
        observerNextId: 0,
        observerList:   make([] *observer, 0),
        observerIndex:  make(map[uint64] uint),
        terminated:     false,
        maybeError:     nil,
        notifyingFlag:  false,
        recentValues:   make([] Object, 0, replay),
        runtimeHandle:  h,
    } }
    for _, item := range items { b.appendRecentValue(item) }
    return b
}
func (b Subject) Observe() Observable {
    return Observable(func(pub DataPublisher) {
        if b.terminated {
            if b.maybeError != nil {
                var err = b.maybeError
                pub.observer.error(err)
            } else {
                b.iterateRecentValues(pub.observer.value)
                pub.observer.complete()
            }
            return
        }
        var ctx, ob = pub.useInheritedContext()
        var id = b.appendObserver(ob)
        ctx.registerCleaner(func() { b.deleteObserver(id) })
        b.iterateRecentValues(pub.observer.value)
    })
}
func (b Subject) Plug(o Observable) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        pub.run(o, ctx, &observer {
            value:    b.value,
            error:    ob.error,
            complete: ob.complete,
        })
    })
}
func (b Subject) value(v Object) {
    if b.terminated {
        return
    }
    b.appendRecentValue(v)
    b.iterateObservers(func(ob *observer) {
        ob.value(v)
    })
}
func (b Subject) error(err error) {
    if b.terminated {
        return
    }
    b.terminated, b.maybeError = true, err
    b.iterateObservers(func(ob *observer) {
        ob.error(err)
    })
}
func (b Subject) complete() {
    if b.terminated {
        return
    }
    b.terminated, b.maybeError = true, nil
    b.iterateObservers(func(ob *observer) {
        ob.complete()
    })
}
func (b Subject) multicastInput() *observer {
    return &observer {
        value:    b.value,
        error:    b.error,
        complete: b.complete,
    }
}
func (b Subject) iterateObservers(k func(*observer)) {
    var L = len(b.observerList)
    if L > 0 {
        if b.notifyingFlag {
            var h = b.runtimeHandle
            Crash(h, InvariantViolation, "synchronous feedback")
        }
        b.notifyingFlag = true
        var snapshot = make([] *observer, L)
        copy(snapshot, b.observerList)
        for _, ob := range snapshot {
            k(ob)
        }
        b.notifyingFlag = false
    }
}
func (b Subject) appendObserver(ob *observer) uint64 {
    var id = b.observerNextId
    var pos = uint(len(b.observerList))
    b.observerList = append(b.observerList, ob)
    b.observerIndex[id] = pos
    b.observerNextId = (id + 1)
    return id
}
func (b Subject) deleteObserver(id uint64) {
    var pos, exists = b.observerIndex[id]
    if !(exists) { panic("invalid argument") }
    // update index
    delete(b.observerIndex, id)
    for current, _ := range b.observerIndex {
        if current > id {
            // position left shifted
            b.observerIndex[current] -= 1
        }
    }
    // update queue
    b.observerList[pos] = nil
    var L = uint(len(b.observerList))
    if !(L >= 1) { panic("something went wrong") }
    for i := pos; i < (L-1); i += 1 {
        b.observerList[i] = b.observerList[i + 1]
    }
    b.observerList[L-1] = nil
    b.observerList = b.observerList[:L-1]
}
func (b Subject) iterateRecentValues(f func(Object)) {
    var L = len(b.recentValues)
    if L > 0 {
        var snapshot = make([] Object, L)
        copy(snapshot, b.recentValues)
        for _, item := range snapshot {
            f(item)
        }
    }
}
func (b Subject) appendRecentValue(v Object) {
    var L = len(b.recentValues)
    if L < cap(b.recentValues) {
        b.recentValues = append(b.recentValues, v)
    } else if L > 0 {
        for i := 0; i < (L - 1); i += 1 {
            b.recentValues[i] = b.recentValues[i+1]
        }
        b.recentValues[L-1] = v
    }
}
func Multicast(o Observable, h RuntimeHandle) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        var bus = CreateSubject(h, 0)
        ob.value(Obj(bus.Observe()))
        ob.complete()
        pub.run(o, ctx, bus.multicastInput())
    })
}
func Loopback(k func(Observable)(Observable), h RuntimeHandle) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        var bus = CreateSubject(h, 0)
        pub.run(bus.Observe(), ctx, ob)
        pub.run(k(bus.Observe()), ctx, bus.multicastInput())
    })
}
func SkipSync(o Observable) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        var sync = true
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                if sync {
                    return
                }
                ob.value(v)
            },
            error:    ob.error,
            complete: ob.complete,
        })
        sync = false
    })
}


func Now() Observable {
    return Observable(func(pub DataPublisher) {
        pub.observer.value(ObjTimeNow())
        pub.observer.complete()
    })
}
func SetTimeout(ms int) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        pub.addTimer(ms, 1, ctx, func() {
            ob.value(nil)
            ob.complete()
        })
    })
}
func SetInterval(ms int, n int) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        if n == 0 {
            ob.complete()
            return
        }
        var current = 0
        pub.addTimer(ms, n, ctx, func() {
            if (n < 0) || (current < n) {
                ob.value(ObjInt(current + 1))
                current += 1
            }
            if (n >= 0) && (current == n) {
                ob.complete()
            }
        })
    })
}
func ObservableSequence(forEach func(func(Observable))) Observable {
    return Observable(func(pub DataPublisher) {
        forEach(func(item Observable) {
            pub.observer.value(Obj(item))
        })
        pub.observer.complete()
    })
}

func (o Observable) Catch(f func(error,Observable)(Observable)) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        pub.run(o, ctx, &observer {
            error: func(err error) {
                pub.run(f(err,o.Catch(f)), ctx, ob)
            },
            value:    ob.value,
            complete: ob.complete,
        })
    })
}
func (o Observable) Retry(limit int) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        var retrial = 0
        var proxy *observer
        proxy = &observer {
            error: func(err error) {
                if retrial == limit {
                    ob.error(err)
                    return
                }
                retrial++
                pub.run(o, ctx, proxy)
            },
            value:    ob.value,
            complete: ob.complete,
        }
        pub.run(o, ctx, proxy)
    })
}
func (o Observable) LogError(h RuntimeHandle) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        pub.run(o, ctx, &observer {
            error: func(err error) {
                MakeLogger(h).LogError(err)
                ob.complete()
            },
            value:    ob.value,
            complete: ob.complete,
        })
    })
}

func (o Observable) DistinctUntilChanged(equal func(Object,Object)(bool)) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        var previous Object
        var available = false
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                if available {
                    if equal(v, previous) {
                        return
                    }
                }
                previous = v
                available = true
                ob.value(v)
            },
            error:    ob.error,
            complete: ob.complete,
        })
    })
}
func (o Observable) DistinctUntilObjectChanged() Observable {
    return o.DistinctUntilChanged(func(a Object, b Object) bool {
        return (a == b)
    })
}
func (o Observable) ObjectPairEqualities() Observable {
    return o.Map(func(obj Object) Object {
        var a, b = GetPair(obj)
        return ObjBool(a == b)
    }).DistinctUntilChanged(func(p Object, q Object) bool {
        return (GetBool(p) == GetBool(q))
    })
}

func (o Observable) WithLatestFrom(another Observable) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, dispose, ob = pub.useNewChildContext()
        var attached Object
        var available = false
        pub.run(another, ctx, &observer {
            value: func(v Object) {
                attached = v
                available = true
            },
            error: func(err error) {
                dispose()
                ob.error(err)
            },
            complete: func() {},
        })
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                if available {
                    ob.value(ObjPair(v, attached))
                }
            },
            error: func(err error) {
                dispose()
                ob.error(err)
            },
            complete: func() {
                dispose()
                ob.complete()
            },
        })
    })
}
func (o Observable) MapToLatestFrom(another Observable) Observable {
    return o.WithLatestFrom(another).Map(func(obj Object) Object {
        var _, latest = GetPair(obj)
        return latest
    })
}
func (o Observable) WithCycle(l List) Observable {
    return Observable(func(pub DataPublisher) {
        if l.Empty() {
            pub.observer.complete()
            return
        }
        var ctx, ob = pub.useInheritedContext()
        var node = l.Head
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                ob.value(ObjPair(v, node.Value))
                node = node.Next
                if node == nil { node = l.Head }
            },
            error:    ob.error,
            complete: ob.complete,
        })
    })
}
func (o Observable) WithIndex() Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        var i = 0
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                ob.value(ObjPair(v, ObjInt(i)))
                i += 1
            },
            error:    ob.error,
            complete: ob.complete,
        })
    })
}
func (o Observable) WithTime() Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                ob.value(ObjPair(v, ObjTimeNow()))
            },
            error:    ob.error,
            complete: ob.complete,
        })
    })
}

func (o Observable) DelayRun(ms int) Observable {
    return SetTimeout(ms).Then(o)
}
func (o Observable) DelayValues(ms int) Observable {
    return o.MergeMap(func(v Object) Observable {
        return SetTimeout(ms).MapTo(v)
    })
}
func (o Observable) StartWith(first Object) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        ob.value(first)
        pub.run(o, ctx, ob)
    })
}
func (o Observable) EndWith(last Object) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        pub.run(o, ctx, &observer {
            value: ob.value,
            error: ob.error,
            complete: func() {
                ob.value(last)
                ob.complete()
            },
        })
    })
}
func (o Observable) Throttle(f func(Object)(Observable)) Observable {
    return o.ExhaustMap(func(v Object) Observable {
        return f(v).CompleteOnEmit().StartWith(v)
    })
}
func (o Observable) Debounce(f func(Object)(Observable)) Observable {
    return o.SwitchMap(func(v Object) Observable {
        return f(v).CompleteOnEmit().EndWith(v)
    })
}
func (o Observable) ThrottleTime(ms int) Observable {
    return o.Throttle(func(_ Object) Observable {
        return SetTimeout(ms)
    })
}
func (o Observable) DebounceTime(ms int) Observable {
    return o.Debounce(func(_ Object) Observable {
        return SetTimeout(ms)
    })
}
func (o Observable) CompleteOnEmit() Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, dispose, ob = pub.useNewChildContext()
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                dispose()
                ob.complete()
            },
            error: func(e error) {
                dispose()
                ob.error(e)
            },
            complete: func() {
                dispose()
                ob.complete()
            },
        })
    })
}

func (o Observable) Skip(n int) Observable {
    if n <= 0 {
        return o
    }
    return Observable(func(pub DataPublisher) {
        var ctx, dispose, ob = pub.useNewChildContext()
        var i = 0
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                if i >= n {
                    ob.value(v)
                }
                i += 1
            },
            error: func(err error) {
                dispose()
                ob.error(err)
            },
            complete: func() {
                dispose()
                ob.complete()
            },
        })
    })
}

func (o Observable) Take(limit int) Observable {
    if limit <= 0 {
        return Observable(func(pub DataPublisher) {
            pub.observer.complete()
        })
    }
    return Observable(func(pub DataPublisher) {
        var ctx, dispose, ob = pub.useNewChildContext()
        var i = 0
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                ob.value(v)
                i += 1
                if i == limit {
                    dispose()
                    ob.complete()
                }
            },
            error: func(err error) {
                dispose()
                ob.error(err)
            },
            complete: func() {
                dispose()
                ob.complete()
            },
        })
    })
}
func (o Observable) TakeLast() Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        var last Object
        var available = false
        pub.run(o, ctx, &observer {
            error: ob.error,
            value: func(v Object) {
                last = v
                available = true
            },
            complete: func() {
                if available {
                    ob.value(last)
                }
                ob.complete()
            },
        })
    })
}
func (o Observable) TakeWhile(f func(Object)(bool)) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, dispose, ob = pub.useNewChildContext()
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                if f(v) {
                    ob.value(v)
                } else {
                    dispose()
                    ob.complete()
                }
            },
            error: func(err error) {
                dispose()
                ob.error(err)
            },
            complete: func() {
                dispose()
                ob.complete()
            },
        })
    })
}
func (o Observable) TakeUntil(stop Observable) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, dispose, ob = pub.useNewChildContext()
        pub.run(o, ctx, &observer {
            value: ob.value,
            error: func(err error) {
                dispose()
                ob.error(err)
            },
            complete: func() {
                dispose()
                ob.complete()
            },
        })
        var stop_ctx, stop_dispose = ctx.createChild()
        pub.run(stop, stop_ctx, &observer {
            value: func(_ Object) {
                dispose()
                ob.complete()
            },
            error: func(err error) {
                dispose()
                ob.error(err)
            },
            complete: func() {
                stop_dispose()
            },
        })
    })
}

func (o Observable) Count() Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        var count = 0
        pub.run(o, ctx, &observer {
            error: ob.error,
            value: func(_ Object) {
                count++
            },
            complete: func() {
                ob.value(ObjInt(count))
                ob.complete()
            },
        })
    })
}
func (o Observable) Collect() Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        var buf ListBuilder
        pub.run(o, ctx, &observer {
            error: ob.error,
            value: func(v Object) {
                buf.Append(v)
            },
            complete: func() {
                ob.value(Obj(buf.Collect()))
                ob.complete()
            },
        })
    })
}
func (o Observable) BufferTime(ms int) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, dispose, ob = pub.useNewChildContext()
        var buf ListBuilder
        var renew = func() Object {
            var l = buf.Collect()
            buf = ListBuilder {}
            return Obj(l)
        }
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                buf.Append(v)
            },
            error: func(err error) {
                dispose()
                ob.error(err)
            },
            complete: func() {
                dispose()
                ob.value(renew())
                ob.complete()
            },
        })
        pub.run(SetInterval(ms, -1), ctx, &observer {
            value: func(_ Object) {
                ob.value(renew())
            },
            error:    func(_ error) { panic("something went wrong") },
            complete: func() { panic("something went wrong") },
        })
    })
}

func (o Observable) Pairwise() Observable {
    return o.BufferCount(2).Map(QueueToPair)
}
func (o Observable) BufferCount(n int) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        var sb = createSlidingBuffer(n)
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                var buf, ok = sb.append(v)
                if ok {
                    ob.value(buf)
                }
            },
            error:    ob.error,
            complete: ob.complete,
        })
    })
}
type slidingBuffer struct {
    size  int
    mq    ctn.MutQueue[Object]
}
func createSlidingBuffer(size int) *slidingBuffer {
    if size < 1 {
        size = 1
    }
    return &slidingBuffer {
        size: size,
        mq:   ctn.MakeMutQueue[Object](),
    }
}
func (sb *slidingBuffer) append(v Object) (Object, bool) {
    sb.mq.Append(v)
    var diff = (sb.mq.Size() - sb.size)
    if diff < 0 {
        return nil, false
    } else if diff == 0 {
        var buf = ObjQueue(sb.mq.Queue())
        return buf, true
    } else {
        sb.mq.Shift()
        var buf = ObjQueue(sb.mq.Queue())
        return buf, true
    }
}

func (o Observable) Map(f func(Object)(Object)) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        pub.run(o, ctx, &observer {
            value:    func(v Object) { ob.value(f(v)) },
            error:    ob.error,
            complete: ob.complete,
        })
    })
}
func (o Observable) MapTo(v Object) Observable {
    return o.Map(func(_ Object)(Object) { return v })
}
func (o Observable) Filter(f func(Object)(bool)) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        pub.run(o, ctx, &observer {
            value:    func(v Object) { if f(v) { ob.value(v) } },
            error:    ob.error,
            complete: ob.complete,
        })
    })
}
func (o Observable) Scan(seed Object, f func(Object,Object)(Object)) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        var current = seed
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                current = f(current, v)
                ob.value(current)
            },
            error:    ob.error,
            complete: ob.complete,
        })
    })
}
func (o Observable) Reduce(initial Object, f func(Object,Object)(Object)) Observable {
    return o.Scan(initial, f).StartWith(initial)
}

func CombineLatest(all ...Observable) Observable {
    return Observable(func(pub DataPublisher) {
        if len(all) == 0 {
            pub.observer.complete()
            return
        }
        var ctx, dispose, ob = pub.useNewChildContext()
        var vv = createValueVector(len(all))
        var completed = 0
        for i_, o := range all {
            var index = i_
            pub.run(o, ctx, &observer {
                value: func(v Object) {
                    vv.SetItem(index, v)
                    if l, ok := vv.GetList(); ok {
                        ob.value(l)
                    }
                },
                error: func(err error) {
                    dispose()
                    ob.error(err)
                },
                complete: func() {
                    completed++
                    if completed == len(all) {
                        dispose()
                        ob.complete()
                    }
                },
            })
        }
    })
}
func (o Observable) CombineLatest(another Observable) Observable {
    return CombineLatest(o, another).Map(ListToPair)
}
type valueVector struct {
    values     [] Object
    available  [] bool
}
func createValueVector(size int) *valueVector {
    return &valueVector {
        values:    make([] Object, size),
        available: make([] bool, size),
    }
}
func (vv *valueVector) SetItem(index int, value Object) {
    vv.values[index] = value
    if vv.available != nil {
        vv.available[index] = true
    }
}
func (vv *valueVector) GetList() (Object, bool) {
    var ok = true
    if vv.available != nil {
        for _, item_ok := range vv.available {
        if !(item_ok) {
            ok = false
        }}
        if ok { vv.available = nil }
    }
    if ok {
        return ObjList(vv.values), true
    } else {
        return nil, false
    }
}

func (o Observable) Await(k func(Object)(Observable)) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        var current Object
        var ok = false
        pub.run(o, ctx, &observer {
            error: ob.error,
            value: func(v Object) {
                current = v
                ok = true
            },
            complete: func() {
                if ok {
                    var last = current
                    pub.run(k(last), ctx, ob)
                } else {
                    ob.complete()
                }
            },
        })
    })
}
func (o Observable) AwaitNoexcept(h RuntimeHandle, k func(Object)(Observable)) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        pub.run(o, ctx, awaitNoexceptObserver(ob, h, func(obj Object) {
            pub.run(k(obj), ctx, ob)
        }))
    })
}
func awaitNoexceptObserver(ob *observer, h RuntimeHandle, k func(Object)) *observer {
    var value, ok = Object(nil), false
    return &observer {
        value: func(v Object) {
            value = v
            ok = true
        },
        error: func(err error) {
            Crash(h, InvariantViolation, ("unexpected error: " + err.Error()))
        },
        complete: func() {
            if ok {
                k(value)
            } else {
                ob.complete()
            }
        },
    }
}

func (o Observable) Then(k Observable) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        pub.run(o, ctx, &observer {
            error:    ob.error,
            value:    func(_ Object) {},
            complete: func() { pub.run(k, ctx, ob) },
        })
    })
}

func (o Observable) With(bg Observable, log func(error)) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        pub.run(bg, ctx, &observer {
            error:    log,
            value:    func(_ Object) {},
            complete: func() {},
        })
        pub.run(o, ctx, ob)
    })
}
func (o Observable) And(bg Observable, log func(error)) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, ob = pub.useInheritedContext()
        pub.run(o, ctx, ob)
        pub.run(bg, ctx, &observer {
            error:    log,
            value:    func(_ Object) {},
            complete: func() {},
        })
    })
}

func Merge(forEach func(func(Observable))) Observable {
    return ObservableSequence(forEach).MergeMap(GetObservable)
}
func (o Observable) MergeMap(f func(Object)(Observable)) Observable {
    return o.MergeMap2(f, false)
}
func (o Observable) MergeMap2(f func(Object)(Observable), auto_stop bool) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, dispose, ob = pub.useNewChildContext()
        var proxy = createMergeProxy(ob, dispose)
        if auto_stop {
            proxy.outerClosed = true
        }
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                var inner = f(v)
                proxy.innerStart()
                pub.run(inner, ctx, &observer {
                    value:    proxy.pass,
                    error:    proxy.abort,
                    complete: proxy.innerExit,
                })
            },
            error:    proxy.abort,
            complete: proxy.outerClose,
        })
    })
}
type mergeProxy struct {
    observer      *observer
    ctxDispose    func()
    innerRunning  uint
    outerClosed   bool
}
func createMergeProxy(ob *observer, dispose func()) *mergeProxy {
    return &mergeProxy {
        observer:     ob,
        ctxDispose:   dispose,
        innerRunning: 0,
        outerClosed:  false,
    }
}
func (p *mergeProxy) pass(x Object) {
    p.observer.value(x)
}
func (p *mergeProxy) abort(e error) {
    p.observer.error(e)
    p.ctxDispose()
}
func (p *mergeProxy) innerStart() {
    p.innerRunning++
}
func (p *mergeProxy) innerExit() {
    if p.innerRunning == 0 { panic("something went wrong") }
    p.innerRunning--
    if p.innerRunning == 0 && p.outerClosed {
        p.observer.complete()
        p.ctxDispose()
    }
}
func (p *mergeProxy) outerClose() {
    p.outerClosed = true
    if p.innerRunning == 0 {
        p.observer.complete()
        p.ctxDispose()
    }
}

func Concat(forEach func(func(Observable))) Observable {
    return ObservableSequence(forEach).ConcatMap(GetObservable)
}
func Concurrent(limit int, forEach func(func(Observable))) Observable {
    return ObservableSequence(forEach).ConcurrentMap(limit, GetObservable)
}
func (o Observable) ConcatMap(f func(Object)(Observable)) Observable {
    return o.ConcurrentMap(1, f)
}
func (o Observable) ConcurrentMap(limit int, f func(Object)(Observable)) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, dispose, ob = pub.useNewChildContext()
        var proxy = createMergeProxy(ob, dispose)
        var queue = createRunQueue(limit, pub, ctx, &observer {
            value:    proxy.pass,
            error:    proxy.abort,
            complete: proxy.innerExit,
        })
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                var inner = f(v)
                proxy.innerStart()
                queue.push(inner)
            },
            error:    proxy.abort,
            complete: proxy.outerClose,
        })
    })
}
func ForkJoin(limit int, all ...Observable) Observable {
    return Observable(func(pub DataPublisher) {
        if len(all) == 0 {
            pub.observer.complete()
            return
        }
        var concurrent = Concurrent(limit, func(yield func(Observable)) {
            for i_, o := range all {
                var index = i_
                yield(o.Map(func(obj Object) Object {
                    return ObjPair(obj, ObjInt(index))
                }))
            }
        })
        var ctx, ob = pub.useInheritedContext()
        var vv = createValueVector(len(all))
        pub.run(concurrent, ctx, &observer {
            error: ob.error,
            value: func(obj Object) {
                var v, i_ = GetPair(obj)
                var index = GetInt(i_)
                vv.SetItem(index, v)
            },
            complete: func() {
                if l, ok := vv.GetList(); ok {
                    ob.value(l)
                }
                ob.complete()
            },
        })
    })
}
func (o Observable) ForkJoin(limit int, another Observable) Observable {
    return ForkJoin(limit, o, another).Map(ListToPair)
}
type runQueue struct {
    waiting    [] Observable
    running    int
    limit      int
    publisher  DataPublisher
    context    *context
    observer   *observer
}
func createRunQueue(limit int, pub DataPublisher, ctx *context, ob *observer) *runQueue {
    if limit < 1 {
        limit = 1
    }
    var q runQueue
    var proxy = &observer {
        value: ob.value,
        error: ob.error,
        complete: func() {
            ob.complete()
            q.running--
            if len(q.waiting) > 0 {
                var next = q.waiting[0]
                q.waiting[0] = nil
                q.waiting = q.waiting[1:]
                q.publisher.run(next, q.context, q.observer)
            }
        },
    }
    q = runQueue {
        waiting:   make([] Observable, 0),
        running:   0,
        limit:     limit,
        publisher: pub,
        context:   ctx,
        observer:  proxy,
    }
    return &q
}
func (q *runQueue) push(o Observable) {
    if q.running < q.limit {
        q.running++
        q.publisher.run(o, q.context, q.observer)
    } else {
        q.waiting = append(q.waiting, o)
    }
}

func (o Observable) SwitchMap(f func(Object)(Observable)) Observable {
    return Observable(func(pub DataPublisher) {
        var disposeInner func()
        var wrap = func(o Observable) Observable {
            return Observable(func(pub DataPublisher) {
                if disposeInner != nil {
                    disposeInner()
                }
                var ctx, dispose, ob = pub.useNewChildContext()
                disposeInner = func() {
                    dispose()
                    ob.complete()
                }
                pub.run(o, ctx, ob)
            })
        }
        var ctx, dispose, ob = pub.useNewChildContext()
        var proxy = createMergeProxy(ob, dispose)
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                var inner = f(v)
                proxy.innerStart()
                pub.run(wrap(inner), ctx, &observer {
                    value:    proxy.pass,
                    error:    proxy.abort,
                    complete: proxy.innerExit,
                })
            },
            error:    proxy.abort,
            complete: proxy.outerClose,
        })
    })
}
func (o Observable) ExhaustMap(f func(Object)(Observable)) Observable {
    return Observable(func(pub DataPublisher) {
        var ctx, dispose, ob = pub.useNewChildContext()
        var proxy = createMergeProxy(ob, dispose)
        pub.run(o, ctx, &observer {
            value: func(v Object) {
                if proxy.innerRunning == 0 {
                    var inner = f(v)
                    proxy.innerStart()
                    pub.run(inner, ctx, &observer {
                        value:    proxy.pass,
                        error:    proxy.abort,
                        complete: proxy.innerExit,
                    })
                }
            },
            error:    proxy.abort,
            complete: proxy.outerClose,
        })
    })
}


