package core

import (
	"fmt"
	"errors"
	"rxgui/standalone/qt"
	"rxgui/standalone/ctn"
	"rxgui/standalone/util/pseudounion"
)


type guiObjectGuard struct {
	ptr *qt.Object
}
func makeGuiObjectGuard(Q qt.Object, dispose func()) (guiObjectGuard, func()) {
	var g = guiObjectGuard { &Q }
	// var class = Q.ClassName()
	// println("new", g.ptr, class)
	return g, func() {
		// println("delete", g.ptr, class)
		if *(g.ptr) != (qt.Object{}) {
			*(g.ptr) = (qt.Object{})
			dispose()
		} else {
			panic("double-free of GuiObject")
		}
	}
}
func (g guiObjectGuard) deref(h RuntimeHandle) qt.Object {
	if *(g.ptr) != (qt.Object{}) {
		return *(g.ptr)
	} else {
		return Crash1[qt.Object](h, UseAfterFree, "GuiObject already deleted")
	}
}

type GuiObject[T any] struct {
	g  guiObjectGuard
	f  func(qt.Object)(T)
}
func makeGuiObject[T any] (Q qt.Object, dispose func(), f func(qt.Object)(T)) (GuiObject[T], func()) {
	var g, delete_ = makeGuiObjectGuard(Q, dispose)
	return GuiObject[T] { g, f }, delete_
}
func (q GuiObject[T]) Deref(h RuntimeHandle) T {
	return q.f(q.g.deref(h))
}
func q_Widget(Q qt.Object) qt.Widget { return qt.Widget{Object:Q} }
func q_Action(Q qt.Object) qt.Action { return qt.Action{Object:Q} }
func q_ActionGroup(Q qt.Object) qt.ActionGroup { return qt.ActionGroup{Object:Q} }

type Action struct {
	GuiObject[qt.Action]
}
func WrapAction(k func(qt.Pkg)(qt.Action)) (Action, func()) {
	var pkg, dispose = qt.CreatePkg()
	var Q = k(pkg).Object
	var q, delete_ = makeGuiObject(Q, dispose, q_Action)
	return Action{q}, delete_
}

type Widget struct {
	GuiObject[qt.Widget]
}
func WrapWidget(k func(qt.Pkg)(qt.Widget)) (Widget, func()) {
	var pkg, dispose = qt.CreatePkg()
	var Q = k(pkg).Object
	var q, delete_ = makeGuiObject(Q, dispose, q_Widget)
	return Widget{q}, delete_
}
func WrapWidget2[X any] (k func(qt.Pkg)(X,qt.Widget)) (X, Widget, func()) {
	var pkg, dispose = qt.CreatePkg()
	var x, W = k(pkg)
	var Q = W.Object
	var q, delete_ = makeGuiObject(Q, dispose, q_Widget)
	var w = Widget{q}
	return x, w, delete_
}

type Signal struct {
	signature     string
	getValue      func(qt.Widget) Object
	getOnConnect  bool
}
//go:noinline
func MakeSignal(sig string) Signal {
	return Signal { sig, nil, false }
}
func MakeSignalWithValueGetter(sig string, v0 bool, f func(qt.Widget)(Object)) Signal {
	return Signal { sig, f, v0 }
}
func (s Signal) withPropGetter(f func(qt.Widget)(Object)) Signal {
	return Signal { s.signature, f, true }
}
func (s Signal) Connect(w Widget, h RuntimeHandle) Observable {
	return Observable(func(pub DataPublisher) {
		var pkg, dispose = qt.CreatePkg()
		if ((s.getValue != nil) && s.getOnConnect) {
			var current = s.getValue(w.Deref(h))
			pub.observer.value(current)
		}
		qt.Connect(w.Deref(h).Object, s.signature, pkg, func() {
			if s.getValue != nil {
				pub.observer.value(s.getValue(w.Deref(h)))
			} else {
				pub.observer.value(nil)
			}
		})
		pub.context.registerCleaner(dispose)
	})
}

type Events struct {
	kind     qt.EventKind
	prevent  bool
}
//go:noinline
func MakeEvents(kind qt.EventKind, prevent bool) Events {
	return Events { kind, prevent }
}
func (e Events) Listen(w Widget, h RuntimeHandle) Observable {
	return Observable(func(pub DataPublisher) {
		var pkg, dispose = qt.CreatePkg()
		var callback = eventTransformer(e.kind, pub.observer.value)
		qt.Listen(w.Deref(h).Object, e.kind, e.prevent, pkg, callback)
		pub.context.registerCleaner(dispose)
	})
}
func eventTransformer(kind qt.EventKind, yield func(Object)) func(qt.Event) {
	return func(ev qt.Event) {
		for _, k := range trivialEventKinds {
			if kind == k {
				yield(nil)
				return
			}
		}
		switch kind {
			// ...
		default:
			panic("unsupported event kind")
		}
	}
}
var trivialEventKinds = [] qt.EventKind {
	qt.EventShow(), qt.EventClose(),
}

type Prop struct {
	name   string
	type_  PropType
}
type PropType int
const ( PropString PropType = iota; PropBool; PropInt )
//go:noinline
func MakeProp(type_ PropType, name string) Prop {
	return Prop { name, type_ }
}
func (p Prop) get(w qt.Widget) Object {
	switch p.type_ {
	case PropString:
		return ObjString(w.GetPropString(p.name))
	case PropBool:
		return ObjBool(w.GetPropBool(p.name))
	case PropInt:
		return ObjInt(w.GetPropInt(p.name))
	default:
		panic("unsupported prop type")
	}
}
func (p Prop) set(v Object, w qt.Widget) {
	switch p.type_ {
	case PropString:
		w.SetPropString(p.name, GetString(v))
	case PropBool:
		w.SetPropBool(p.name, GetBool(v))
	case PropInt:
		w.SetPropInt(p.name, GetInt(v))
	default:
		panic("unsupported prop type")
	}
}
func (p Prop) Read(s Signal) Signal {
	return s.withPropGetter(p.get)
}
func (p Prop) Bind(o Observable, w Widget, h RuntimeHandle) Observable {
	return o.ConcatMap(func(v Object) Observable {
		return Observable(func(pub DataPublisher) {
			p.set(v, w.Deref(h))
			pub.observer.complete()
		})
	})
}

func DialogGenerate(pub DataPublisher) (func(Object), func()) {
	return pub.observer.value, pub.observer.complete
}
func ConnectActionTrigger(a Action, h RuntimeHandle) Observable {
	return onSync(func() func(qt.Pkg, func()) {
		return a.Deref(h).OnTrigger
	})
}
func SetActionEnabled(a Action, p bool, h RuntimeHandle) Observable {
	return doSync(func() {
		a.Deref(h).SetEnabled(p)
	})
}
func ActionCheckBox(a Action, initial bool, h RuntimeHandle, k func(Observable)(Object)) Hook {
	return MakeHook(func() (Object, func()) {
		if a.Deref(h).Checkable() {
			Crash(h, InvariantViolation, "action already checkable")
		}
		a.Deref(h).SetCheckable(true)
		var set_uncheckable = func() {
			a.Deref(h).SetCheckable(false)
		}
		a.Deref(h).SetChecked(initial)
		var checked = onSync1(func() func(qt.Pkg, func(bool)) {
			return a.Deref(h).OnCheck
		})
		return k(checked), set_uncheckable
	})
}
type ActionComboBoxItem struct { Action Action; Selected bool }
func ActionComboBox(items ([] ActionComboBoxItem), h RuntimeHandle, k func(Observable)(Object)) Hook {
	if len(items) == 0 {
		panic("invalid argument")
	}
	return MakeHook(func() (Object, func()) {
		for _, item := range items {
			if item.Action.Deref(h).Checkable() {
				Crash(h, InvariantViolation, "some action already checkable")
			}
			item.Action.Deref(h).SetCheckable(true)
		}
		var set_uncheckable = func() {
			for _, item := range items {
				item.Action.Deref(h).SetCheckable(false)
			}
		}
		var items = ctn.MapEach(items, func(item ActionComboBoxItem) qt.ActionGroupItem {
			return qt.ActionGroupItem {
				Action:   item.Action.Deref(h),
				Selected: item.Selected,
			}
		})
		var pkg, dispose = qt.CreatePkg()
		var ok bool
		var G = qt.CreateActionGroup(pkg, items, &ok)
		if !(ok) {
			Crash(h, InvariantViolation, "some action already grouped")
		}
		var Q = G.Object
		var g, delete_ = makeGuiObject(Q, dispose, q_ActionGroup)
		var index_on_trigger = onSync1(func() func(qt.Pkg, func(int)) {
			return g.Deref(h).OnCheck
		})
		var index = index_on_trigger.DistinctUntilChanged(func(a Object, b Object) bool {
			return (GetInt(a) == GetInt(b))
		})
		return k(index), func() {
			delete_()
			set_uncheckable()
		}
	})
}

func BindInlineStyleSheet(w Widget, o Observable, h RuntimeHandle) Observable {
	return ObservableFlattenLast(doSync2(func() (Observable, func()) {
		var id, initial = addInlineStyleSheet(w, h)
		var selector = fmt.Sprintf(`[%s="%s"]`, qt.DP_InlineStyleSheetId, id)
		var value = o.Map(func(obj Object) Object {
			var decls = GetString(obj)
			var ruleset = fmt.Sprintf("%s { %s }\n", selector, decls)
			return ObjString(ruleset + initial)
		})
		var prop = MakeProp(PropString, qt.P_StyleSheet)
		return prop.Bind(value, w, h), func() {
			removeInlineStyleSheet(initial, w, h)
		}
	}))
}
func addInlineStyleSheet(w Widget, h RuntimeHandle) (string, string) {
	var W = w.Deref(h)
	if W.GetPropString(qt.DP_InlineStyleSheetId) != "" {
		Crash(h, InvariantViolation, "duplicate styling")
	}
	var id = fmt.Sprintf("%x", w.Deref(h).PointerNumber())
	W.SetPropString(qt.DP_InlineStyleSheetId, id)
	return id, W.GetPropString(qt.P_StyleSheet)
}
func removeInlineStyleSheet(initial string, w Widget, h RuntimeHandle) {
	var W = w.Deref(h)
	W.SetPropString(qt.P_StyleSheet, initial)
	W.SetPropString(qt.DP_InlineStyleSheetId, "")
}

func BindContextMenu(w Widget, m_ qt.Menu, h RuntimeHandle) Observable {
	return Observable(func(pub DataPublisher) {
		var pkg, dispose = qt.CreatePkg()
		var m = qt.CreateContextMenu(m_, pkg)
		var ok = m.Bind(w.Deref(h), pkg)
		if !(ok) {
			Crash(h, InvariantViolation, "duplicate context menu")
		}
		pub.context.registerCleaner(dispose)
	})
}

func CreateDynamicWidget(widgets Observable, h RuntimeHandle) Hook {
	return Hook { Observable(func(pub DataPublisher) {
		var ctx, ob = pub.useInheritedContext()
		var W, w, w_dispose = WrapWidget2(func(ctx qt.Pkg) (qt.DynamicWidget, qt.Widget) {
			var W = qt.CreateDynamicWidget(ctx)
			return W, W.Widget
		})
		ctx.registerCleaner(w_dispose)
		pub.run(widgets, ctx, &observer {
			value: func(obj Object) {
				W.SetWidget(GetWidget(obj).Deref(h))
			},
			error:    ErrorLogger(h),
			complete: func() {},
		})
		ob.value(Obj(w))
		ob.complete()
	})}
}

type ListViewConfig struct {
	CreateInterface  func(ctx qt.Pkg) qt.Lwi
	ReturnObject     func(w Widget, e Observable, c Observable, s Observable) Object
}
func ListView(config ListViewConfig, data Observable, getKey func(Object)(string), p ItemViewProvider, h RuntimeHandle) Hook {
	return Hook { Observable(func(pub DataPublisher) {
		var ctx, ob = pub.useInheritedContext()
		var lg = MakeLogger(h)
		var I, w, w_dispose = WrapWidget2(func(ctx qt.Pkg) (qt.Lwi, qt.Widget) {
			var I = config.CreateInterface(ctx)
			return I, I.CastToWidget()
		})
		ctx.registerCleaner(w_dispose)
		var m = make(map[string] *listViewItem)
		var update = withLwiLatestKeys(data, I, h, func(obj Object, old_keys_ ([] string)) Observable {
			var new_keys_ = make([] string, 0)
			var new_keys = make(map[string] struct{})
			var mapping = make(map[string] listViewItemValueWithPos)
			var list = GetList(obj)
			var total = list.Length()
			list.ForEachWithIndex(func(index int, value Object) {
				var key = getKey(value)
				var _, dup = new_keys[key]
				if !(dup) {
					new_keys_ = append(new_keys_, key)
					new_keys[key] = struct{}{}
					mapping[key] = listViewItemValueWithPos {
						Value: value,
						Pos:   ItemPos { index, total },
					}
				} else {
					lg.LogError(errors.New("duplicate key ignored: " + key))
				}
			})
			var tasks = make([] Observable, 0)
			var old_keys = make(map[string] struct{})
			for _, K := range old_keys_ {
				var old_key = K
				old_keys[old_key] = struct{}{}
				var _, exists = new_keys[old_key]
				if removed := !(exists); removed {
					var remove = doSync(func() {
						m[old_key].Dispose()
					})
					tasks = append(tasks, remove)
				}
			}
			for _, K := range new_keys_ {
				var new_key = K
				var new_value = mapping[new_key]
				if _, exists := old_keys[new_key]; exists {
					var item_update = doSync(func() {
						m[new_key].Buffer.value(ToObject(new_value))
						I.Update(new_key)
					})
					tasks = append(tasks, item_update)
				} else {
					var add = retrieveBufAndHook(new_key, p, h, func(buf Subject, hook Hook) Observable {
						return retrieveObjectInChildContext(ctx, hook.Job, h, func(view ItemView, ctx *context, dispose func()) Observable {
						return doSync(func() {
							var widgets = view.Widgets.Deref(h)
							var extension = view.Extension
							I.Append(new_key, widgets)
							ctx.registerCleaner(func() {
								I.Delete(new_key)
							})
							var item = &listViewItem {
								Context:   ctx,
								Dispose:   dispose,
								Buffer:    buf,
								Extension: extension,
							}
							m[new_key] = item
							ctx.registerCleaner(func() {
								delete(m, new_key)
							})
							buf.value(ToObject(new_value))
							I.Update(new_key)
						})
					})})
					tasks = append(tasks, add)
				}
			}
			var reorder = doSync(func() {
				I.Reorder(new_keys_)
			})
			tasks = append(tasks, reorder)
			return Concat(func(yield func(Observable)) {
				for _, task := range tasks {
					yield(task)
				}
			})
		})
		var watch = func(name string, k func(Object)(Observable)) Observable {
			return MakeSignal(name).Connect(w, h).StartWith(nil).ConcatMap(k)
		}
		var e = watch(qt.DefaultListWidget_CurrentChanged, func(_ Object) Observable {
			return doSync1(func() ctn.Maybe[Widget] {
				if cur, ok := I.Current(); ok {
				if item, ok := m[cur]; ok {
					return item.Extension
				}}
				return nil
			})
		})
		var c = watch(qt.DefaultListWidget_CurrentChanged, func(_ Object) Observable {
			return doSync1(func() ctn.Maybe[string] {
				return ctn.MakeMaybe(I.Current())
			})
		})
		var s = watch(qt.DefaultListWidget_SelectionChanged, func(_ Object) Observable {
			return doSync1(func() ([] string) {
				return I.Selection()
			})
		})
		pub.run(update, ctx, &observer {
			value:    func(Object) {},
			error:    lg.LogError,
			complete: func() {},
		})
		ob.value(config.ReturnObject(w, e, c, s))
		ob.complete()
	})}
}
func withLwiLatestKeys(data Observable, I qt.Lwi, h RuntimeHandle, k func(Object,[]string)(Observable)) Observable {
	return data.ConcatMap(func(obj Object) Observable {
		var get_keys = doSync1(func() ([] string) { return I.All() })
		return retrieveObject(get_keys, h, func(keys ([] string)) Observable {
			return k(obj, keys)
		})
	})
}
func retrieveBufAndHook(key string, p ItemViewProvider, h RuntimeHandle, k func(Subject,Hook)(Observable)) Observable {
	var create_buf = doSync1(func() Subject {
		return CreateSubject(h, 1)
	})
	return retrieveObject(create_buf, h, func(buf Subject) Observable {
		var o = listViewItemValue(buf.Observe())
		var pos = listViewItemPos(buf.Observe())
		var info = ItemInfo {
			Key: key,
			Pos: pos,
		}
		var hook = p(o, info)
		return k(buf, hook)
	})
}
type listViewItem struct {
	Context    *context
	Dispose    func()
	Buffer     Subject
	Extension  ctn.Maybe[Widget]
}
type listViewItemValueWithPos struct {
	Value  Object
	Pos    ItemPos
}
func listViewItemValue(o Observable) Observable {
	var raw = o.Map(func(obj Object) Object {
		return FromObject[listViewItemValueWithPos](obj).Value
	})
	return raw.DistinctUntilObjectChanged()
}
func listViewItemPos(o Observable) Observable {
	var raw = o.Map(func(obj Object) Object {
		return ToObject(FromObject[listViewItemValueWithPos](obj).Pos)
	})
	return DistinctUntilItemPosChanged(raw)
}

type ListEditViewConfig struct {
	CreateInterface  func(ctx qt.Pkg) qt.Lwi
	ReturnObject     func(w Widget, o Observable, e Observable, O Subject) Object
}
func ListEditView(config ListEditViewConfig, initial List, p ItemEditViewProvider, h RuntimeHandle) Hook {
	return Hook { Observable(func(pub DataPublisher) {
		var ctx, ob = pub.useInheritedContext()
		var lg = MakeLogger(h)
		var I, w, w_dispose = WrapWidget2(func(ctx qt.Pkg) (qt.Lwi, qt.Widget) {
			var I = config.CreateInterface(ctx)
			return I, I.CastToWidget()
		})
		ctx.registerCleaner(w_dispose)
		var O = CreateSubject(h, 0)
		var store = createItemEditDataStore(h)
		var update = O.Observe().ConcatMap(func(op__ Object) Observable {
			var op_ = FromObject[ListEditOperation](op__)
			switch op := pseudounion.Load(op_).(type) {
			case Prepend, Append, InsertAbove, InsertBelow:
				var v, i = listEditViewMatchInsertionOperation(op)
				return retrieveKeyAndHook(v, p, store, h, func(key string, hook Hook) Observable {
				return retrieveObjectInChildContext(ctx, hook.Job, h, func(view ItemEditView, ctx *context, dispose func()) Observable {
					var ext = view.Extension
					var ops = view.EditOps(key)
					var enable_ops = O.Plug(SkipSync(ops))
					var insert = doSync(func() {
						var widgets = view.Widgets.Deref(h)
						store.insert(key, I, i, widgets, ctx, dispose, v, ext)
					})
					return insert.And(enable_ops, lg.LogError)
				})})
			case Update:
				return doSync(func() {
					store.update(op.Key, I, op.Value)
				})
			case Delete:
				return doSync(func() {
					store.delete(op.Key, I)
				})
			case MoveUp:
				return doSync(func() {
					store.moveUp(op.Key, I)
				})
			case MoveDown:
				return doSync(func() {
					store.moveDown(op.Key, I)
				})
			case MoveTop:
				return doSync(func() {
					store.moveTop(op.Key, I)
				})
			case MoveBottom:
				return doSync(func() {
					store.moveBottom(op.Key, I)
				})
			case Reorder:
				return doSync(func() {
					store.reorder(I, op.Reorder, lg)
				})
			default:
				panic("impossible branch")
			}
		})
		var o = store.output()
		var watch = func(name string, k func(Object)(Observable)) Observable {
			return MakeSignal(name).Connect(w, h).StartWith(nil).ConcatMap(k)
		}
		var e = watch(qt.DefaultListWidget_CurrentChanged, func(_ Object) Observable {
			return doSync1(func() ctn.Maybe[Widget] {
				return store.getCurrentExtension(I)
			})
		})
		pub.run(update, ctx, &observer {
			value:    func(Object) {},
			error:    lg.LogError,
			complete: func() {},
		})
		initial.ForEach(func(v Object) {
			O.value(ToObject(pseudounion.Store[ListEditOperation] (
				Append { v },
			)))
		})
		ob.value(config.ReturnObject(w, o, e, O))
		ob.complete()
	})}
}
var listEditKeyCounter = uint64(0)
func generateListEditKey() string {
	var i = listEditKeyCounter
	listEditKeyCounter++
	var key = fmt.Sprint(i)
	return key
}
func retrieveKeyAndHook(value Object, p ItemEditViewProvider, store *listEditViewDataStore, h RuntimeHandle, k func(string,Hook)(Observable)) Observable {
	var gen_key = doSync1(generateListEditKey)
	return retrieveObject(gen_key, h, func(key string) Observable {
		return k(key, p(value, store.info(key)))
	})
}
type listEditViewDataStore struct {
	ItemsMap      map[string] *listEditItem
	PositionMap   map[string] ItemPos
	OutputBuffer  Subject
}
type listEditItem struct {
	Context    *context
	Dispose    func()
	Value      Object
	Extension  ctn.Maybe[Widget]
}
func createItemEditDataStore(h RuntimeHandle) *listEditViewDataStore {
	return &listEditViewDataStore {
		ItemsMap:     make(map[string] *listEditItem),
		PositionMap:  make(map[string] ItemPos),
		OutputBuffer: CreateSubject(h, 1),
	}
}
func (s *listEditViewDataStore) getCurrentExtension(I qt.Lwi) ctn.Maybe[Widget] {
	if cur, ok := I.Current(); ok {
	if item, ok := s.ItemsMap[cur]; ok {
		return item.Extension
	}}
	return nil
}
func (s *listEditViewDataStore) output() Observable {
	return s.OutputBuffer.Observe()
}
func (s *listEditViewDataStore) info(key string) ItemInfo {
	return ItemInfo {
		Key: key,
		Pos: s.pos(key),
	}
}
func (s *listEditViewDataStore) pos(key string) Observable {
	var raw = Observable(func(pub DataPublisher) {
		var ctx, ob = pub.useInheritedContext()
		pub.run(s.output(), ctx, &observer {
			value: func(_ Object) {
				if pos, ok := s.PositionMap[key]; ok {
					ob.value(ToObject(pos))
				}
			},
			error:    ob.error,
			complete: ob.complete,
		})
	})
	return DistinctUntilItemPosChanged(raw)
}
func (s *listEditViewDataStore) __updateOutput(I qt.Lwi) {
	var keys = I.All()
	var L = len(keys)
	var pm = make(map[string] ItemPos)
	var nodes = make([] ListNode, len(keys))
	for i, key := range keys {
		pm[key] = ItemPos { i, L }
		nodes[i].Value = s.ItemsMap[key].Value
	}
	var list = NodesToList(nodes)
	s.PositionMap = pm
	s.OutputBuffer.value(Obj(list))
}
type listEditViewInsertion func(I qt.Lwi, key string, widgets ([] qt.Widget))
func listEditViewPrepend() listEditViewInsertion {
	return func(I qt.Lwi, key string, widgets ([] qt.Widget)) {
		I.Prepend(key, widgets)
	}
}
func listEditViewAppend() listEditViewInsertion {
	return func(I qt.Lwi, key string, widgets ([] qt.Widget)) {
		I.Append(key, widgets)
	}
}
func listEditViewInsertAbove(pivot string) listEditViewInsertion {
	return func(I qt.Lwi, key string, widgets []qt.Widget) {
		I.InsertAbove(pivot, key, widgets)
	}
}
func listEditViewInsertBelow(pivot string) listEditViewInsertion {
	return func(I qt.Lwi, key string, widgets []qt.Widget) {
		I.InsertBelow(pivot, key, widgets)
	}
}
func listEditViewMatchInsertionOperation(op_ interface{}) (Object,listEditViewInsertion) {
	switch op := op_.(type) {
	case Prepend:
		return op.Value, listEditViewPrepend()
	case Append:
		return op.Value, listEditViewAppend()
	case InsertAbove:
		return op.Value, listEditViewInsertAbove(op.PivotKey)
	case InsertBelow:
		return op.Value, listEditViewInsertBelow(op.PivotKey)
	default:
		panic("invalid argument")
	}
}
func (s *listEditViewDataStore) insert(key string, I qt.Lwi, i listEditViewInsertion, widgets ([] qt.Widget), ctx *context, dispose func(), val Object, ext ctn.Maybe[Widget]) {
	var _, exists = s.ItemsMap[key]
	if exists { panic("something went wrong") }
	var add, remove = i, func() {
		I.Delete(key)
	}
	add(I, key, widgets)
	ctx.registerCleaner(remove)
	var item = &listEditItem {
		Context:   ctx,
		Dispose:   dispose,
		Value:     val,
		Extension: ext,
	}
	s.ItemsMap[key] = item
	ctx.registerCleaner(func() {
		delete(s.ItemsMap, key)
	})
	s.__updateOutput(I)
	I.Update(key)
}
func (s *listEditViewDataStore) update(key string, I qt.Lwi, val Object) {
	if val != s.ItemsMap[key].Value {
		s.ItemsMap[key].Value = val
		s.__updateOutput(I)
		I.Update(key)
	}
}
func (s *listEditViewDataStore) delete(key_ ctn.Maybe[string], I qt.Lwi) {
	for _, key := range getListEditOperandKeys(false, key_, I) {
		s.ItemsMap[key].Dispose()
		s.__updateOutput(I)
	}
}
func (s *listEditViewDataStore) moveUp(key_ ctn.Maybe[string], I qt.Lwi) bool {
	var keys = getListEditOperandKeys(true, key_, I)
	if len(keys) == 0 {
		return false
	}
	for i, key := range keys {
		var ok = I.MoveUp(key)
		if ((i == 0) && !(ok)) {
			return false
		}
		s.__updateOutput(I)
	}
	return true
}
func (s *listEditViewDataStore) moveDown(key_ ctn.Maybe[string], I qt.Lwi) bool {
	var keys = ctn.Reverse(getListEditOperandKeys(true, key_, I))
	if len(keys) == 0 {
		return false
	}
	for i, key := range keys {
		var ok = I.MoveDown(key)
		if ((i == 0) && !(ok)) {
			return false
		}
		s.__updateOutput(I)
	}
	return true
}
func (s *listEditViewDataStore) moveTop(key_ ctn.Maybe[string], I qt.Lwi) {
	for s.moveUp(key_, I) {}
}
func (s *listEditViewDataStore) moveBottom(key_ ctn.Maybe[string], I qt.Lwi) {
	for s.moveDown(key_, I) {}
}
func (s *listEditViewDataStore) reorder(I qt.Lwi, f func(List)(List), lg Logger) {
	var buf ListBuilder
	var rev = make(map[Object] string)
	for key, item := range s.ItemsMap {
		var value = Object(new(ObjectImpl))
		*value = *(item.Value)
		buf.Append(value)
		rev[value] = key
	}
	var L = len(rev)
	var order = make([] string, L)
	var reordered = make([] Object, 0)
	var reordered_ = f(buf.Collect())
	reordered_.ForEach(func(value Object) {
		reordered = append(reordered, value)
	})
	if len(reordered) == L {
		for i, value := range reordered {
			if key, ok := rev[value]; ok {
				delete(rev, value)
				order[i] = key
			} else {
				goto NG
			}
		}
		I.Reorder(order)
		return
	}
	NG:
	lg.LogError(errors.New("invalid reorder, ignored"))
}
func getListEditOperandKeys(contiguous bool, key_ ctn.Maybe[string], I qt.Lwi) ([] string) {
	if key, ok := key_.Value(); ok {
		return [] string { key }
	} else {
		if contiguous {
			return I.ContiguousSelection()
		} else {
			return I.Selection()
		}
	}
}


