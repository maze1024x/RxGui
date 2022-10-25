package builtin

import (
    "os"
    "io"
    "fmt"
    "math"
    "time"
    "regexp"
    "errors"
    "strconv"
    "strings"
    "unicode"
    "reflect"
    "runtime"
    "net/url"
    "math/big"
    "math/rand"
    "unicode/utf8"
    "path/filepath"
    "rxgui/qt"
    "rxgui/util/ctn"
    "rxgui/util"
    "rxgui/util/pseudounion"
    "rxgui/interpreter/core"
)


var functionList = [] interface{} {
    // Reflection
    ReflectType,
    DebugInspect, DebugExpose, DebugTrace, DebugWatch,
    Serialize, Deserialize,
    ObjectPairEqualities,
    // Null, Error, undefined (panic)
    Null, Error, ErrorMessage, ErrorIsCancel, ErrorWrap, Undefined,
    // Bool
    BoolNo, BoolYes, BoolNot, BoolAnd, BoolOr, BoolEqual, BoolAssert,
    // Int
    IntPlus, IntMinus, IntTimes, IntQuo, IntRem, IntPow,
    IntEqual, IntLessThan, IntCompare,
    // Float
    FloatPlus, FloatMinus, FloatTimes, FloatQuo, FloatRem, FloatPow,
    FloatEqual, FloatLessThan,
    FloatInt, IntFloat, FloatNormal, FloatNaN, FloatInfinite,
    math.Floor, math.Ceil, math.Round,
    math.Sqrt, math.Cbrt, math.Exp, math.Log,
    math.Sin, math.Cos, math.Tan, math.Asin, math.Acos, math.Atan, math.Atan2,
    // Char
    Char,
    CharInt, CharUtf8Size,
    CharEqual, CharLessThan, CharCompare,
    // String
    String, StringFromChars, Quote, Unquote,
    StringEmpty,
    StringChars, StringFirstChar, StringNumberOfChars,
    StringUtf8Size,
    StringEqual, StringLessThan, StringCompare,
    StringShift, StringReverse,
    ListStringJoin, StringSplit, StringCut,
    StringHasPrefix, StringHasSuffix,
    StringTrimPrefix, StringTrimSuffix,
    StringTrim, StringTrimLeft, StringTrimRight,
    // RegExp
    RegExpString,
    StringAdvance, StringSatisfy, StringReplace,
    // To String
    BoolString, OrderingString, IntString, FloatString, CharString,
    // From String
    ParseInt, ParseFloat,
    // List
    ListConcat,
    core.Cons, core.Count, ListEmpty, ListFirst, ListLength, ListSeq,
    ListShift, ListReverse,
    ListSort, ListTake,
    ListWithIndex,
    ListMap, ListDeflateMap, ListFlatMap, ListFilter, ListScan, ListFold,
    // Seq
    SeqEmpty, SeqLast, SeqLength, SeqList,
    SeqAppend, SeqDeflateAppend, SeqFlatAppend, SeqSort, SeqFilter,
    // Queue
    Queue,
    QueueEmpty, QueueSize, QueueFirst, QueueList,
    QueueShift, QueueAppend,
    // Heap
    Heap,
    HeapEmpty, HeapSize, HeapFirst, HeapList,
    HeapShift, HeapInsert,
    // Set
    Set,
    SetEmpty, SetSize, SetList,
    SetHas, SetDelete, SetInsert,
    // Map
    Map,
    MapEmpty, MapSize, MapKeys, MapValues, MapEntries,
    MapHas, MapLookup, MapDelete, MapInsert,
    // Observable
    Observable,
    Go,
    core.WithChildContext, core.WithCancelTrigger, core.WithCancelTimeout,
    core.SetTimeout, core.SetInterval,
    Throw, Crash,
    ObservableCatch, ObservableRetry, ObservableLogError,
    ObservableDistinctUntilChanged,
    ObservableWithLatestFrom, ObservableMapToLatestFrom,
    ObservableWithCycle, ObservableWithIndex, ObservableWithTime,
    ObservableDelaySubscription, ObservableDelayValues,
    ObservableStartWith, ObservableEndWith,
    ObservableThrottle, ObservableDebounce,
    ObservableThrottleTime, ObservableDebounceTime,
    ObservableCompleteOnEmit,
    ObservableSkip,
    ObservableTake,
    ObservableTakeLast, ObservableTakeLastAsMaybe,
    ObservableTakeWhile, ObservableTakeWhileMaybeOK,
    ObservableTakeUntil,
    ObservableCount, ObservableCollect,
    ObservableBufferTime,
    ObservablePairwise, ObservableBufferCount,
    ObservableMap, ObservableMapTo,
    ObservableFilter,
    ObservableScan, ObservableReduce,
    ObservableCombineLatest, ListObservableCombineLatest,
    ObservableAwait, ObservableAwaitNoexcept,
    ObservableThen,
    ObservableWith, ObservableAnd,
    ObservableAutoMap,
    ListObservableMerge, ObservableMerge, ObservableMergeMap,
    ListObservableConcat, ObservableConcat, ObservableConcatMap,
    ObservableSwitchMap, ObservableExhaustMap,
    NumCPU,
    ListObservableConcurrent, ObservableConcurrentMap,
    ListObservableForkJoin, ObservableForkJoin,
    UUID,
    Random,
    Shuffle,
    // Subject
    CreateSubject,
    SubjectValues, SubjectPlug,
    core.Multicast, core.Loopback, core.SkipSync,
    // Time
    TimeString, TimeSubtractMillisecond, core.Now,
    // Request
    Get, Post, Put, Delete, Subscribe,
    // File
    FileString, FileEqual,
    ReadTextFile, WriteTextFile,
    // Config
    ReadConfig,
    WriteConfig,
    // Process
    Arguments,
    Environment,
    // GUI
    // general
    FontSize,
    // standard dialogs
    ShowInfo, ShowWarning, ShowCritical,
    ShowYesNo, ShowAbortRetryIgnore, ShowSaveDiscardCancel,
    GetChoice, GetLine, GetText, GetInt, GetFloat,
    GetFileListToOpen, GetFileToOpen, GetFileToSave,
    // action, context menu
    Action, ActionTriggers,
    ActionCheckBox, ActionComboBox,
    BindContextMenu,
    // widget
    ShowAndActivate,
    BindInlineStyleSheet,
    ComboBoxSelectedItem,
    core.CreateDynamicWidget,
    CreateWidget,
    CreateScrollArea,
    CreateGroupBox,
    CreateSplitter,
    CreateMainWindow,
    CreateDialog,
    CreateLabel,
    CreateIconLabel,
    CreateElidedLabel,
    CreateTextView,
    CreateCheckBox,
    CreateComboBox,
    CreatePushButton,
    CreateLineEdit,
    CreatePlainTextEdit,
    CreateSlider,
    CreateProgressBar,
    // signal and events
    Connect, Listen,
    SignalToggled,
    SignalClicked,
    SignalTextChanged0,
    SignalTextChanged1,
    SignalReturnPressed,
    SignalValueChanged,
    EventsShow,
    EventsClose,
    // prop
    Read,
    Bind,
    ClearTextLater,
    PropEnabled,
    PropWindowTitle,
    PropText,
    PropChecked,
    PropPlainText,
    PropValue,
    // list
    ListView,
    ListEditView,
}
var functionMap = (func() (map[string] core.NativeFunction) {
    var m = make(map[string] core.NativeFunction)
    for _, item := range functionList {
        var rv = reflect.ValueOf(item)
        var name = strings.Split(runtime.FuncForPC(rv.Pointer()).Name(), ".")[1]
        var obj = core.MakeNativeFunction(item)
        if _, duplicate := m[name]; duplicate {
            panic("duplicate native function: " + name)
        }
        m[name] = obj
    }
    return m
})()
func LookupFunction(id string) (core.NativeFunction, bool) {
    var f, ok = functionMap[id]
    return f, ok
}


func ReflectType() core.Object {
    // DummyReflectType
    return nil
}
func DebugInspect(hint string, rv core.ReflectValue, h core.RuntimeHandle) core.Object {
    var lg = core.MakeLogger(h)
    lg.Inspect(rv.Value(), rv.Type().CertainType(), hint)
    return rv.Value()
}
func DebugExpose(name string, rv core.ReflectValue, h core.RuntimeHandle) core.Observable {
    var lg = core.MakeLogger(h)
    return lg.Expose(name, rv.Value(), rv.Type().CertainType())
}
func DebugTrace(hint string, rv core.ReflectValue, h core.RuntimeHandle) core.Observable {
    var lg = core.MakeLogger(h)
    var ro = rv.CastToReflectObservable()
    return lg.Trace(ro.ObservableValue(), ro.InnerType().CertainType(), hint)
}
func DebugWatch(hint string, rv core.ReflectValue, h core.RuntimeHandle) core.Observable {
    var lg = core.MakeLogger(h)
    var ro = rv.CastToReflectObservable()
    return lg.Watch(ro.ObservableValue(), ro.InnerType().CertainType(), hint)
}
func Serialize(rv core.ReflectValue, h core.RuntimeHandle) core.Observable {
    var ctx = h.SerializationContext()
    return core.Observable(func(pub core.DataPublisher) {
        go (func() {
            var b, err = core.Marshal(rv, ctx)
            if err != nil { pub.AsyncThrow(err); return }
            var s = util.WellBehavedDecodeUtf8(b)
            pub.AsyncReturn(core.ObjString(s))
        })()
    })
}
func Deserialize(s string, rt core.ReflectType, h core.RuntimeHandle) core.Observable {
    var ctx = h.SerializationContext()
    return core.Observable(func(pub core.DataPublisher) {
        go (func() {
            var b = ([] byte)(s)
            var obj, err = core.Unmarshal(b, rt, ctx)
            if err != nil { pub.AsyncThrow(err); return }
            pub.AsyncReturn(obj)
        })()
    })
}
func marshal(rv core.ReflectValue, h core.RuntimeHandle) core.Observable {
    var ctx = h.SerializationContext()
    return core.Observable(func(pub core.DataPublisher) {
        go (func() {
            var b, err = core.Marshal(rv, ctx)
            if err != nil { pub.AsyncThrow(err); return }
            pub.AsyncReturn(core.ObjBytes(b))
        })()
    })
}
func unmarshal(b ([] byte), rt core.ReflectType, h core.RuntimeHandle) core.Observable {
    var ctx = h.SerializationContext()
    return core.Observable(func(pub core.DataPublisher) {
        go (func() {
            var obj, err = core.Unmarshal(b, rt, ctx)
            if err != nil { pub.AsyncThrow(err); return }
            pub.AsyncReturn(obj)
        })()
    })
}
func ObjectPairEqualities(o core.Observable) core.Observable {
    return o.ObjectPairEqualities()
}

func Null() core.Object {
    return nil
}
func Error(msg string) error {
    return errors.New(msg)
}
func ErrorMessage(err error) string {
    return err.Error()
}
func ErrorIsCancel(err error) bool {
    var _, is_cancel = err.(core.CancelError)
    return is_cancel
}
func ErrorWrap(err error, msg string) error {
    return fmt.Errorf("%s: %w", msg, err)
}
func Undefined(msg string, h core.RuntimeHandle) core.Object {
    return core.Crash1[core.Object] (
        h, core.ValueUndefined,
        strconv.Quote(msg),
    )
}

func BoolNo() bool {
    return false
}
func BoolYes() bool {
    return true
}
func BoolNot(p bool) bool {
    return !(p)
}
func BoolAnd(p bool, q bool) bool {
    return (p && q)
}
func BoolOr(p bool, q bool) bool {
    return (p || q)
}
func BoolEqual(p bool, q bool) bool {
    return (p == q)
}
func BoolAssert(ok bool, k func()(core.Object), h core.RuntimeHandle) core.Object {
    if ok {
        return k()
    } else {
        return core.Crash1[core.Object] (
            h, core.AssertionFailed,
            "given boolean value is false",
        )
    }
}

type CompareWithOperator func(core.Object,core.Object)(ctn.Ordering)
type EqualToOperator func(core.Object,core.Object)(bool)
type LessThanOperator func(core.Object,core.Object)(bool)

func IntPlus(a *big.Int, b *big.Int) *big.Int {
    var c big.Int
    return c.Add(a, b)
}
func IntMinus(a *big.Int, b *big.Int) *big.Int {
    var c big.Int
    return c.Sub(a, b)
}
func IntTimes(a *big.Int, b *big.Int) *big.Int {
    var c big.Int
    return c.Mul(a, b)
}
func IntQuo(a *big.Int, b *big.Int, h core.RuntimeHandle) *big.Int {
    if b.Sign() == 0 {
        core.Crash(h, core.InvalidArgument, "division by zero")
    }
    var c big.Int
    return c.Quo(a, b)
}
func IntRem(a *big.Int, b *big.Int, h core.RuntimeHandle) *big.Int {
    if b.Sign() == 0 {
        core.Crash(h, core.InvalidArgument, "division by zero")
    }
    var c big.Int
    return c.Rem(a, b)
}
func IntPow(a *big.Int, b *big.Int, h core.RuntimeHandle) *big.Int {
    var c big.Int
    if b.Sign() < 0 {
        core.Crash(h, core.InvalidArgument, "negative integer power")
    }
    return c.Exp(a, b, nil)
}
func IntEqual(a *big.Int, b *big.Int) bool {
    return (a.Cmp(b) == 0)
}
func IntLessThan(a *big.Int, b *big.Int) bool {
    return (a.Cmp(b) < 0)
}
func IntCompare(a *big.Int, b *big.Int) ctn.Ordering {
    var result = a.Cmp(b)
    if result < 0 {
        return ctn.Smaller
    } else if result > 0 {
        return ctn.Bigger
    } else {
        return ctn.Equal
    }
}

func FloatPlus(x float64, y float64) float64 {
    return (x + y)
}
func FloatMinus(x float64, y float64) float64 {
    return (x - y)
}
func FloatTimes(x float64, y float64) float64 {
    return (x * y)
}
func FloatQuo(x float64, y float64) float64 {
    return (x / y)
}
func FloatRem(x float64, y float64) float64 {
    return math.Mod(x, y)
}
func FloatPow(x float64, y float64) float64 {
    return math.Pow(x, y)
}
func FloatEqual(x float64, y float64) bool {
    return (x == y)
}
func FloatLessThan(x float64, y float64) bool {
    return (x < y)
}
func FloatInt(x float64, h core.RuntimeHandle) *big.Int {
    var a, ok = util.DoubleToInteger(x)
    if !(ok) {
        core.Crash(h, core.InvalidArgument,
            "cannot convert abnormal float to integer")
    }
    return a
}
func IntFloat(a *big.Int, h core.RuntimeHandle) float64 {
    var x, ok = util.IntegerToDouble(a)
    if !(ok) {
        core.Crash(h, core.InvalidArgument,
            "cannot convert large integer to float")
    }
    return x
}
func FloatNormal(x float64) bool {
    return util.IsNormalFloat(x)
}
func FloatNaN(x float64) bool {
    return (x != x)
}
func FloatInfinite(x float64) bool {
    return math.IsInf(x, 0)
}

func Char(n *big.Int) rune {
    if n.IsInt64() {
        n := n.Int64()
        if (0 <= n && n <= 0x10FFFF) && !(0xD800 <= n && n <= 0xDFFF) {
            return rune(n)
        }
    }
    return unicode.ReplacementChar
}
func CharInt(c rune) int {
    return int(c)
}
func CharUtf8Size(c rune) int {
    return utf8.RuneLen(c)
}
func CharEqual(c rune, d rune) bool {
    return (c == d)
}
func CharLessThan(c rune, d rune) bool {
    return (c < d)
}
func CharCompare(c rune, d rune) ctn.Ordering {
    return ctn.DefaultCompare[rune](c, d)
}

func String(fragments ([] ToString)) string {
    var buf strings.Builder
    for _, f := range fragments {
        buf.WriteString(f())
    }
    return buf.String()
}
func StringFromChars(chars ([] rune)) string {
    return string(chars)
}
func Quote(s string) string {
    return strconv.Quote(s)
}
func Unquote(s string) ctn.Maybe[string] {
    var unquoted, err = strconv.Unquote(s)
    if err == nil {
        return ctn.Just(unquoted)
    } else {
        return nil
    }
}
func StringEmpty(s string) bool {
    return (s == "")
}
func StringChars(s string) ([] rune) {
    return ([] rune)(s)
}
func StringFirstChar(s string) ctn.Maybe[rune] {
    for _, char := range s {
        return ctn.Just(char)
    }
    return nil
}
func StringNumberOfChars(s string) int {
    var n = 0
    for range s {
        n++
    }
    return n
}
func StringUtf8Size(s string) int {
    return len(s)
}
func StringEqual(a string, b string) bool {
    return (a == b)
}
func StringLessThan(a string, b string) bool {
    return (a < b)
}
func StringCompare(a string, b string) ctn.Ordering {
    return ctn.StringCompare(a, b)
}
func StringShift(s string) ctn.Maybe[ctn.Pair[rune,string]] {
    for _, char := range s {
        var rest = s[utf8.RuneLen(char):]
        return ctn.Just(ctn.MakePair(char, rest))
    }
    return nil
}
func StringReverse(s string) string {
    var chars = ([] rune)(s)
    var L = len(chars)
    var buf strings.Builder
    for i := (L-1); i >= 0; i -= 1 {
        buf.WriteRune(chars[i])
    }
    return buf.String()
}
func ListStringJoin(l core.List, sep string) string {
    var values = make([] string, 0)
    l.ForEach(func(item core.Object) {
        var value = core.GetString(item)
        values = append(values, value)
    })
    return strings.Join(values, sep)
}
func StringSplit(s string, sep string) core.List {
    return core.ToObjectList(strings.Split(s, sep))
}
func StringCut(s string, sep string) ctn.Maybe[ctn.Pair[string,string]] {
    if a, b, ok := strings.Cut(s, sep); ok {
        return ctn.Just(ctn.MakePair(a, b))
    } else {
        return nil
    }
}
func StringHasPrefix(s string, prefix string) bool {
    return strings.HasPrefix(s, prefix)
}
func StringHasSuffix(s string, suffix string) bool {
    return strings.HasSuffix(s, suffix)
}
func StringTrimPrefix(s string, prefix string) string {
    return strings.TrimPrefix(s, prefix)
}
func StringTrimSuffix(s string, suffix string) string {
    return strings.TrimSuffix(s, suffix)
}
func StringTrim(s string, chars ([] rune)) string {
    return strings.Trim(s, string(chars))
}
func StringTrimLeft(s string, chars ([] rune)) string {
    return strings.TrimLeft(s, string(chars))
}
func StringTrimRight(s string, chars ([] rune)) string {
    return strings.TrimRight(s, string(chars))
}

func RegExpString(r *regexp.Regexp) string {
    return r.String()
}
func StringAdvance(s string, pattern *regexp.Regexp) ctn.Maybe[ctn.Pair[string,string]] {
    var header_pattern = ("^(?:" + pattern.String() + ")")
    var r = assumeValidRegexp(header_pattern)
    if result := r.FindStringIndex(s); (result != nil) {
        if result[0] != 0 { panic("something went wrong") }
        var pos = result[1]
        var match = s[:pos]
        var rest = s[pos:]
        return ctn.Just(ctn.MakePair(match, rest))
    } else {
        return nil
    }
}
func StringSatisfy(s string, pattern *regexp.Regexp) bool {
    var full_pattern = ("^(?:" + pattern.String() + ")$")
    var r = assumeValidRegexp(full_pattern)
    if result := r.FindStringIndex(s); (result != nil) {
        if result[0] != 0 { panic("something went wrong") }
        if result[1] != len(s) { panic("something went wrong") }
        return true
    } else {
        return false
    }
}
func StringReplace(s string, pattern *regexp.Regexp, f func(string)(string)) string {
    return pattern.ReplaceAllStringFunc(s, f)
}

type ToString func()(string)
func BoolString(p bool) string {
    if p {
        return "Yes"
    } else {
        return "No"
    }
}
func OrderingString(o ctn.Ordering) string {
    return o.String()
}
func IntString(a *big.Int) string {
    return a.String()
}
func FloatString(x float64) string {
    return strconv.FormatFloat(x, 'g', -1, 64)
}
func CharString(c rune) string {
    return string([] rune { c })
}

func ParseInt(s string) ctn.Maybe[*big.Int] {
    var n = new(big.Int)
    if _, ok := n.SetString(s, 10); ok {
        return ctn.Just(n)
    } else {
        return nil
    }
}
func ParseFloat(s string) ctn.Maybe[float64] {
    var x, err = strconv.ParseFloat(s, 64)
    if err == nil {
        return ctn.Just(x)
    } else {
        return nil
    }
}

func ListConcat(l core.List) core.List {
    var buf core.ListBuilder
    l.ForEach(func(item core.Object) {
        core.GetList(item).ForEach(func(item core.Object) {
            buf.Append(item)
        })
    })
    return buf.Collect()
}
func ListEmpty(l core.List) bool {
    return l.Empty()
}
func ListFirst(l core.List) ctn.Maybe[core.Object] {
    return ctn.MakeMaybe(l.First())
}
func ListLength(l core.List) int {
    return l.Length()
}
func ListSeq(l core.List) core.Seq {
    return l.ToSeq()
}
func ListShift(l core.List) ctn.Maybe[ctn.Pair[core.Object,core.List]] {
    var head, tail, ok = l.Shifted()
    return ctn.MakeMaybe(ctn.MakePair(head, tail), ok)
}
func ListReverse(l core.List) core.List {
    return l.Reversed()
}
func ListSort(l core.List, lt LessThanOperator) core.List {
    return l.Sorted(ctn.Less[core.Object](lt))
}
func ListTake(l core.List, n int) core.List {
    return l.Take(n)
}
func ListWithIndex(l core.List) core.List {
    return l.WithIndex()
}
func ListMap(l core.List, f func(core.Object)(core.Object)) core.List {
    return l.Map(f)
}
func ListDeflateMap(l core.List, f func(core.Object)(ctn.Maybe[core.Object])) core.List {
    return l.DeflateMap(f)
}
func ListFlatMap(l core.List, f func(core.Object)(core.List)) core.List {
    return l.FlatMap(f)
}
func ListFilter(l core.List, f func(core.Object)(bool)) core.List {
    return l.Filter(f)
}
func ListScan(l core.List, b core.Object, f func(core.Object,core.Object)(core.Object)) core.List {
    return l.Scan(b, f)
}
func ListFold(l core.List, b core.Object, f func(core.Object,core.Object)(core.Object)) core.Object {
    return l.Fold(b, f)
}

func SeqEmpty(s core.Seq) bool {
    return s.Empty()
}
func SeqLast(s core.Seq) ctn.Maybe[core.Object] {
    return ctn.MakeMaybe(s.Last())
}
func SeqLength(s core.Seq) int {
    return s.Length()
}
func SeqList(s core.Seq) core.List {
    return s.ToList()
}
func SeqAppend(s core.Seq, item core.Object) core.Seq {
    return s.Appended(item)
}
func SeqDeflateAppend(s core.Seq, item ctn.Maybe[core.Object]) core.Seq {
    if item, ok := item.Value(); ok {
        return s.Appended(item)
    } else {
        return s
    }
}
func SeqFlatAppend(s core.Seq, items core.List) core.Seq {
    var draft = s
    items.ForEach(func(item core.Object) {
        draft = draft.Appended(item)
    })
    return draft
}
func SeqSort(s core.Seq, lt LessThanOperator) core.Seq {
    return s.Sorted(ctn.Less[core.Object](lt))
}
func SeqFilter(s core.Seq, f func(core.Object)(bool)) core.Seq {
    return s.Filter(f)
}

func Queue(items core.List) core.Queue {
    var queue = ctn.MakeQueue[core.Object]()
    items.ForEach(func(item core.Object) {
        queue = queue.Appended(item)
    })
    return core.Queue(queue)
}
func QueueEmpty(q core.Queue) bool {
    var queue = ctn.Queue[core.Object](q)
    return queue.IsEmpty()
}
func QueueSize(q core.Queue) int {
    var queue = ctn.Queue[core.Object](q)
    return queue.Size()
}
func QueueFirst(q core.Queue) ctn.Maybe[core.Object] {
    var queue = ctn.Queue[core.Object](q)
    return ctn.MakeMaybe(queue.First())
}
func QueueList(q core.Queue) core.List {
    var queue = ctn.Queue[core.Object](q)
    var buf core.ListBuilder
    queue.ForEach(func(item core.Object) {
        buf.Append(item)
    })
    return buf.Collect()
}
func QueueShift(q core.Queue) ctn.Maybe[ctn.Pair[core.Object,core.Queue]] {
    var queue = ctn.Queue[core.Object](q)
    var v, rest, ok = queue.Shifted()
    return ctn.MakeMaybe(ctn.MakePair(v, core.Queue(rest)), ok)
}
func QueueAppend(q core.Queue, item core.Object) core.Queue {
    var queue = ctn.Queue[core.Object](q)
    return core.Queue(queue.Appended(item))
}

func Heap(items core.List, lt LessThanOperator) core.Heap {
    var heap = ctn.MakeHeap(ctn.Less[core.Object](lt))
    items.ForEach(func(item core.Object) {
        heap = heap.Inserted(item)
    })
    return core.Heap(heap)
}
func HeapEmpty(h core.Heap) bool {
    var heap = ctn.Heap[core.Object](h)
    return heap.IsEmpty()
}
func HeapSize(h core.Heap) int {
    var heap = ctn.Heap[core.Object](h)
    return heap.Size()
}
func HeapFirst(h core.Heap) ctn.Maybe[core.Object] {
    var heap = ctn.Heap[core.Object](h)
    return ctn.MakeMaybe(heap.First())
}
func HeapList(h core.Heap) core.List {
    var heap = ctn.Heap[core.Object](h)
    var buf core.ListBuilder
    heap.ForEach(func(item core.Object) {
        buf.Append(item)
    })
    return buf.Collect()
}
func HeapShift(h core.Heap) ctn.Maybe[ctn.Pair[core.Object,core.Heap]] {
    var heap = ctn.Heap[core.Object](h)
    var v, rest, ok = heap.Shifted()
    return ctn.MakeMaybe(ctn.MakePair(v, core.Heap(rest)), ok)
}
func HeapInsert(h core.Heap, item core.Object) core.Heap {
    var heap = ctn.Heap[core.Object](h)
    return core.Heap(heap.Inserted(item))
}

func Set(items core.List, cmp CompareWithOperator) core.Set {
    var set = ctn.MakeSet(ctn.Compare[core.Object](cmp))
    items.ForEach(func(item core.Object) {
        set = set.Inserted(item)
    })
    return core.Set(set)
}
func SetEmpty(s core.Set) bool {
    var set = ctn.Set[core.Object](s)
    return set.IsEmpty()
}
func SetSize(s core.Set) int {
    var set = ctn.Set[core.Object](s)
    return set.Size()
}
func SetList(s core.Set) core.List {
    var set = ctn.Set[core.Object](s)
    var buf core.ListBuilder
    set.ForEach(func(item core.Object) {
        buf.Append(item)
    })
    return buf.Collect()
}
func SetHas(s core.Set, item core.Object) bool {
    var set = ctn.Set[core.Object](s)
    return set.Has(item)
}
func SetDelete(s core.Set, item core.Object) core.Set {
    var set = ctn.Set[core.Object](s)
    var new_set, _ = set.Deleted(item)
    return core.Set(new_set)
}
func SetInsert(s core.Set, item core.Object) core.Set {
    var set = ctn.Set[core.Object](s)
    return core.Set(set.Inserted(item))
}

func Map(entries core.List, cmp CompareWithOperator) core.Map {
    var map_ = ctn.MakeMap[core.Object,core.Object](ctn.Compare[core.Object](cmp))
    entries.ForEach(func(entry core.Object) {
        var pair = core.FromObject[ctn.Pair[core.Object,core.Object]](entry)
        map_ = map_.Inserted(pair.Key(), pair.Value())
    })
    return core.Map(map_)
}
func MapEmpty(m core.Map) bool {
    var map_ = ctn.Map[core.Object,core.Object](m)
    return map_.IsEmpty()
}
func MapSize(m core.Map) int {
    var map_ = ctn.Map[core.Object,core.Object](m)
    return map_.Size()
}
func MapKeys(m core.Map) core.List {
    var map_ = ctn.Map[core.Object,core.Object](m)
    var buf core.ListBuilder
    map_.ForEach(func(k core.Object, _ core.Object) {
        buf.Append(k)
    })
    return buf.Collect()
}
func MapValues(m core.Map) core.List {
    var map_ = ctn.Map[core.Object,core.Object](m)
    var buf core.ListBuilder
    map_.ForEach(func(_ core.Object, v core.Object) {
        buf.Append(v)
    })
    return buf.Collect()
}
func MapEntries(m core.Map) core.List {
    var map_ = ctn.Map[core.Object,core.Object](m)
    var buf core.ListBuilder
    map_.ForEach(func(k core.Object, v core.Object) {
        buf.Append(core.ToObject(ctn.MakePair(k, v)))
    })
    return buf.Collect()
}
func MapHas(m core.Map, key core.Object) bool {
    var map_ = ctn.Map[core.Object,core.Object](m)
    return map_.Has(key)
}
func MapLookup(m core.Map, key core.Object) ctn.Maybe[core.Object] {
    var map_ = ctn.Map[core.Object,core.Object](m)
    return ctn.MakeMaybe(map_.Lookup(key))
}
func MapDelete(m core.Map, key core.Object) core.Map {
    var map_ = ctn.Map[core.Object,core.Object](m)
    var _, new_map, _ = map_.Deleted(key)
    return core.Map(new_map)
}
func MapInsert(m core.Map, entry ctn.Pair[core.Object,core.Object]) core.Map {
    var map_ = ctn.Map[core.Object,core.Object](m)
    return core.Map(map_.Inserted(entry.Key(), entry.Value()))
}

func Observable(l core.List) core.Observable {
    return core.Observable(func(pub core.DataPublisher) {
        pub.SyncGenerate(func(yield func(core.Object)) error {
            l.ForEach(yield)
            return nil
        })
    })
}
func Go(k func()(core.Object)) core.Observable {
    return core.Observable(func(pub core.DataPublisher) {
        go (func() {
            var v = k()
            pub.AsyncReturn(v)
        })()
    })
}
func Throw(e error) core.Observable {
    return core.Observable(func(pub core.DataPublisher) {
        pub.SyncReturn(func() (core.Object, error) {
            return nil, e
        })
    })
}
func Crash(err error, h core.RuntimeHandle) core.Observable {
    return core.Observable(func(pub core.DataPublisher) {
        pub.SyncReturn(func() (core.Object, error) {
            return core.Crash2[core.Object,error](
                h, core.ProgrammedCrash, err.Error(),
            )
        })
    })
}
func ObservableCatch(o core.Observable, f func(error,core.Observable)(core.Observable)) core.Observable {
    return o.Catch(f)
}
func ObservableRetry(o core.Observable, n int) core.Observable {
    return o.Retry(n)
}
func ObservableLogError(o core.Observable, h core.RuntimeHandle) core.Observable {
    return o.LogError(h)
}
func ObservableDistinctUntilChanged(o core.Observable, eq EqualToOperator) core.Observable {
    return o.DistinctUntilChanged(eq)
}
func ObservableWithLatestFrom(o core.Observable, another core.Observable) core.Observable {
    return o.WithLatestFrom(another)
}
func ObservableMapToLatestFrom(o core.Observable, another core.Observable) core.Observable {
    return o.MapToLatestFrom(another)
}
func ObservableWithCycle(o core.Observable, l core.List) core.Observable {
    return o.WithCycle(l)
}
func ObservableWithIndex(o core.Observable) core.Observable {
    return o.WithIndex()
}
func ObservableWithTime(o core.Observable) core.Observable {
    return o.WithTime()
}
func ObservableDelaySubscription(o core.Observable, ms int) core.Observable {
    return o.DelayRun(ms)
}
func ObservableDelayValues(o core.Observable, ms int) core.Observable {
    return o.DelayValues(ms)
}
func ObservableStartWith(o core.Observable, first core.Object) core.Observable {
    return o.StartWith(first)
}
func ObservableEndWith(o core.Observable, last core.Object) core.Observable {
    return o.EndWith(last)
}
func ObservableThrottle(o core.Observable, f func(core.Object)(core.Observable)) core.Observable {
    return o.Throttle(f)
}
func ObservableDebounce(o core.Observable, f func(core.Object)(core.Observable)) core.Observable {
    return o.Debounce(f)
}
func ObservableThrottleTime(o core.Observable, ms int) core.Observable {
    return o.ThrottleTime(ms)
}
func ObservableDebounceTime(o core.Observable, ms int) core.Observable {
    return o.DebounceTime(ms)
}
func ObservableCompleteOnEmit(o core.Observable) core.Observable {
    return o.CompleteOnEmit()
}
func ObservableSkip(o core.Observable, n int) core.Observable {
    return o.Skip(n)
}
func ObservableTake(o core.Observable, n int) core.Observable {
    return o.Take(n)
}
func ObservableTakeLast(o core.Observable) core.Observable {
    return o.TakeLast()
}
func ObservableTakeLastAsMaybe(o core.Observable) core.Observable {
    return o.TakeLast().Map(core.Just).EndWith(core.Nothing()).Take(1)
}
func ObservableTakeWhile(o core.Observable, f func(core.Object)(bool)) core.Observable {
    return o.TakeWhile(f)
}
func ObservableTakeWhileMaybeOK(o core.Observable) core.Observable {
    return o.TakeWhile(func(obj core.Object) bool {
        var _, ok = core.UnwrapMaybe(obj)
        return ok
    }).Map(func(obj core.Object) core.Object {
        var v, ok = core.UnwrapMaybe(obj)
        if !(ok) { panic("something went wrong") }
        return v
    })
}
func ObservableTakeUntil(o core.Observable, stop core.Observable) core.Observable {
    return o.TakeUntil(stop)
}
func ObservableCount(o core.Observable) core.Observable {
    return o.Count()
}
func ObservableCollect(o core.Observable) core.Observable {
    return o.Collect()
}
func ObservableBufferTime(o core.Observable, ms int) core.Observable {
    return o.BufferTime(ms)
}
func ObservablePairwise(o core.Observable) core.Observable {
    return o.Pairwise()
}
func ObservableBufferCount(o core.Observable, n int) core.Observable {
    return o.BufferCount(n)
}
func ObservableMap(o core.Observable, f func(core.Object)(core.Object)) core.Observable {
    return o.Map(f)
}
func ObservableMapTo(o core.Observable, v core.Object) core.Observable {
    return o.MapTo(v)
}
func ObservableFilter(o core.Observable, f func(core.Object)(bool)) core.Observable {
    return o.Filter(f)
}
func ObservableScan(o core.Observable, init core.Object, f func(core.Object,core.Object)(core.Object)) core.Observable {
    return o.Scan(init, f)
}
func ObservableReduce(o core.Observable, init core.Object, f func(core.Object,core.Object)(core.Object)) core.Observable {
    return o.Reduce(init, f)
}
func ObservableCombineLatest(o core.Observable, another core.Observable) core.Observable {
    return o.CombineLatest(another)
}
func ListObservableCombineLatest(l ([] core.Observable)) core.Observable {
    return core.CombineLatest(l...)
}
func ObservableAwait(o core.Observable, k func(core.Object)(core.Observable)) core.Observable {
    return o.Await(k)
}
func ObservableAwaitNoexcept(o core.Observable, k func(core.Object)(core.Observable), h core.RuntimeHandle) core.Observable {
    return o.AwaitNoexcept(h, k)
}
func ObservableThen(o core.Observable, k core.Observable) core.Observable {
    return o.Then(k)
}
func ObservableWith(o core.Observable, bg core.Observable, h core.RuntimeHandle) core.Observable {
    return o.With(bg, core.ErrorLogger(h))
}
func ObservableAnd(o core.Observable, bg core.Observable, h core.RuntimeHandle) core.Observable {
    return o.And(bg, core.ErrorLogger(h))
}
func ObservableAutoMap(o core.Observable, f func(core.Object)(core.Observable)) core.Observable {
    return o.MergeMap2(f, true)
}
func ListObservableMerge(l core.List) core.Observable {
    return core.Merge(forEachObservable(l))
}
func ObservableMerge(o1 core.Observable, o2 core.Observable) core.Observable {
    return core.Merge(core.YieldObservables(o1, o2))
}
func ObservableMergeMap(o core.Observable, f func(core.Object)(core.Observable)) core.Observable {
    return o.MergeMap(f)
}
func ListObservableConcat(l core.List) core.Observable {
    return core.Concat(forEachObservable(l))
}
func ObservableConcat(o1 core.Observable, o2 core.Observable) core.Observable {
    return core.Concat(core.YieldObservables(o1, o2))
}
func ObservableConcatMap(o core.Observable, f func(core.Object)(core.Observable)) core.Observable {
    return o.ConcatMap(f)
}
func ObservableSwitchMap(o core.Observable, f func(core.Object)(core.Observable)) core.Observable {
    return o.SwitchMap(f)
}
func ObservableExhaustMap(o core.Observable, f func(core.Object)(core.Observable)) core.Observable {
    return o.ExhaustMap(f)
}
func NumCPU() int {
    return runtime.NumCPU()
}
func ListObservableConcurrent(l core.List, n int) core.Observable {
    return core.Concurrent(n, forEachObservable(l))
}
func ObservableConcurrentMap(o core.Observable, f func(core.Object)(core.Observable), n int) core.Observable {
    return o.ConcurrentMap(n, f)
}
func ListObservableForkJoin(l ([] core.Observable), n int) core.Observable {
    return core.ForkJoin(n, l...)
}
func ObservableForkJoin(o core.Observable, another core.Observable, n int) core.Observable {
    return o.ForkJoin(n, another)
}
func UUID() core.Observable {
    return core.Observable(func(pub core.DataPublisher) {
        pub.SyncReturn(func() (core.Object, error) {
            var uuid = qt.UUID()
            return core.ObjString(uuid), nil
        })
    })
}
func Random(supremum *big.Int, h core.RuntimeHandle) core.Observable {
    return core.RandBigInt(supremum, h)
}
func Shuffle(l core.List, h core.RuntimeHandle) core.Observable {
    return core.Rand(h, func(r *rand.Rand) core.Object {
        var nodes = make([] core.ListNode, 0)
        l.ForEach(func(value core.Object) {
            nodes = append(nodes, core.ListNode {
                Value: value,
            })
        })
        r.Shuffle(len(nodes), func(i, j int) {
            nodes[i], nodes[j] = nodes[j], nodes[i]
        })
        var shuffled = core.NodesToList(nodes)
        return core.Obj(shuffled)
    })
}

func CreateSubject(replay int, items ([] core.Object), h core.RuntimeHandle) core.Observable {
    return core.Observable(func(pub core.DataPublisher) {
        pub.SyncReturn(func() (core.Object, error) {
            var bus = core.CreateSubject(h, replay, items...)
            return core.Obj(bus), nil
        })
    })
}
func SubjectValues(b core.Subject) core.Observable {
    return b.Observe()
}
func SubjectPlug(b core.Subject, input core.Observable) core.Observable {
    return b.Plug(input)
}

func TimeString(t time.Time) string {
    return t.String()
}
func TimeSubtractMillisecond(t time.Time, u time.Time) int {
    return int(t.Sub(u).Milliseconds())
}

func Get(endpoint string, rt core.ReflectType, token string, h core.RuntimeHandle) core.Observable {
    return request(core.GET, endpoint, rt, token, nil, h)
}
func Post(data core.ReflectValue, endpoint string, rt core.ReflectType, token string, h core.RuntimeHandle) core.Observable {
    return request(core.POST, endpoint, rt, token, ctn.Just(data), h)
}
func Put(data core.ReflectValue, endpoint string, rt core.ReflectType, token string, h core.RuntimeHandle) core.Observable {
    return request(core.PUT, endpoint, rt, token, ctn.Just(data), h)
}
func Delete(endpoint string, rt core.ReflectType, token string, h core.RuntimeHandle) core.Observable {
    return request(core.DELETE, endpoint, rt, token, nil, h)
}
func Subscribe(endpoint string, rt core.ReflectType, token string, h core.RuntimeHandle) core.Observable {
    return request(core.SUBSCRIBE, endpoint, rt, token, nil, h)
}
func request(method core.RequestMethod, endpoint_ string, rt core.ReflectType, token string, body_ ctn.Maybe[core.ReflectValue], h core.RuntimeHandle) core.Observable {
    var endpoint, err = url.Parse(endpoint_)
    if err != nil {
        var msg = ("invalid endpoint URL: " + endpoint_)
        core.Crash(h, core.InvalidArgument, msg)
    }
    var marshal_req = (func() core.Observable {
        if rv, ok := body_.Value(); ok {
            return marshal(rv, h)
        } else {
            return core.ObservableSyncValue(core.ObjBytes(nil))
        }
    })()
    var unmarshal_resp = func(obj core.Object) core.Observable {
        var binary = core.GetBytes(obj)
        return unmarshal(binary, rt, h)
    }
    return marshal_req.Await(func(obj core.Object) core.Observable {
        var body = core.GetBytes(obj)
        var req = core.Request {
            Method:      method,
            Endpoint:    endpoint,
            AuthToken:   token,
            BodyContent: body,
        }
        var resp = req.Observe(core.MakeLogger(h))
        return resp.ConcatMap(unmarshal_resp)
    })
}

func FileString(f core.File) string {
    return f.Path
}
func FileEqual(f core.File, g core.File) bool {
    return (f.Path == g.Path)
}
func ReadTextFile(f core.File, h core.RuntimeHandle) core.Observable {
    return core.Observable(func(pub core.DataPublisher) {
        go (func() {
            core.MakeLogger(h).LogFileRead(f)
            var content, err = os.ReadFile(f.Path)
            if err != nil { pub.AsyncThrow(err); return }
            var text, ok = util.WellBehavedTryDecodeUtf8(content)
            if !(ok) {
                var err = errors.New("invalid UTF-8")
                { pub.AsyncThrow(err); return }
            }
            pub.AsyncReturn(core.ObjString(text))
        })()
    })
}
func WriteTextFile(f core.File, text string, h core.RuntimeHandle) core.Observable {
    return core.Observable(func(pub core.DataPublisher) {
        go (func() {
            var content = ([] byte)(text)
            core.MakeLogger(h).LogFileWrite(f, content)
            var err = os.WriteFile(f.Path, content, 0666)
            if err != nil { pub.AsyncThrow(err); return }
            pub.AsyncReturn(nil)
        })()
    })
}

func ReadConfig(dir string, name string, default_ core.ReflectValue, h core.RuntimeHandle) core.Observable {
    return core.Observable(func(pub core.DataPublisher) {
        go (func() {
            var ctx = h.SerializationContext()
            var f, err = getConfigFile(dir, name)
            if err != nil { pub.AsyncThrow(err); return }
            if _, err := os.Stat(f); err != nil {
                var binary, err = core.Marshal(default_, ctx)
                if err != nil { pub.AsyncThrow(err); return }
                { var fd, err = os.Create(f)
                if err != nil { pub.AsyncThrow(err); return }
                { var _, err = fd.Write(binary)
                if err != nil { pub.AsyncThrow(err); return }
                pub.AsyncReturn(default_.Value()) }}
            } else {
                var fd, err = os.Open(f)
                if err != nil { pub.AsyncThrow(err); return }
                { var binary, err = io.ReadAll(fd)
                if err != nil { pub.AsyncThrow(err); return }
                { var value, err = core.Unmarshal(binary, default_.Type(), ctx)
                if err != nil { pub.AsyncThrow(err); return }
                pub.AsyncReturn(value) }}
            }
        })()
    })
}
func WriteConfig(dir string, name string, rv core.ReflectValue, h core.RuntimeHandle) core.Observable {
    return core.Observable(func(pub core.DataPublisher) {
        go (func() {
            var ctx = h.SerializationContext()
            var f, err = getConfigFile(dir, name)
            if err != nil { pub.AsyncThrow(err); return }
            { var binary, err = core.Marshal(rv, ctx)
            if err != nil { pub.AsyncThrow(err); return }
            { var fd, err = os.Create(f)
            if err != nil { pub.AsyncThrow(err); return }
            { var _, err = fd.Write(binary)
            if err != nil { pub.AsyncThrow(err); return }
            pub.AsyncReturn(nil) }}}
        })()
    })
}
func getConfigFile(dir string, name string) (string, error) {
    for _, elem := range [] string { dir, name } {
        if strings.ContainsAny(elem, "./\\:") {
            return "", errors.New("invalid dir/file name")
        }
    }
    var config, err = os.UserConfigDir()
    if err != nil { return "", err }
    { var dir = filepath.Join(config, dir)
    _ = os.Mkdir(dir, 0777)
    var name = filepath.Join(dir, name)
    return name, nil }
}

func Arguments(h core.RuntimeHandle) ([] string) {
    return h.ProgramArgs()
}
func Environment() ([] string) {
    return os.Environ()
}

func FontSize() int {
    return qt.FontSize()
}

func ShowInfo(content string, title string) core.Observable {
    return showMessageBox(
        qt.MsgBoxInfo(),
        qt.MsgBoxOK(), qt.MsgBoxOK(),
        title, content,
        func(btn qt.MsgBoxBtn) (core.Object, bool) {
            return nil, (btn == qt.MsgBoxOK())
        },
    )
}
func ShowWarning(content string, title string) core.Observable {
    return showMessageBox(
        qt.MsgBoxWarning(),
        qt.MsgBoxOK(), qt.MsgBoxOK(),
        title, content,
        func(btn qt.MsgBoxBtn) (core.Object, bool) {
            return nil, (btn == qt.MsgBoxOK())
        },
    )
}
func ShowCritical(content string, title string) core.Observable {
    return showMessageBox(
        qt.MsgBoxCritical(),
        qt.MsgBoxOK(), qt.MsgBoxOK(),
        title, content,
        func(btn qt.MsgBoxBtn) (core.Object, bool) {
            return nil, (btn == qt.MsgBoxOK())
        },
    )
}
func ShowYesNo(content string, title string) core.Observable {
    return showMessageBox(
        qt.MsgBoxQuestion(),
        qt.MsgBoxYes().And(qt.MsgBoxNo()),
        qt.MsgBoxYes(),
        title,
        content,
        func(btn qt.MsgBoxBtn) (core.Object, bool) {
            switch btn {
            case qt.MsgBoxYes(): return core.ObjBool(true), true
            case qt.MsgBoxNo():  return core.ObjBool(false), true
            default:             return nil, false
            }
        },
    )
}
func ShowAbortRetryIgnore(content string, title string) core.Observable {
    return showMessageBox(
        qt.MsgBoxCritical(),
        qt.MsgBoxAbort().And(qt.MsgBoxRetry()).And(qt.MsgBoxIgnore()),
        qt.MsgBoxRetry(),
        title,
        content,
        func(btn qt.MsgBoxBtn) (core.Object, bool) {
            switch btn {
            case qt.MsgBoxRetry():  return core.ToObject(RI_Retry), true
            case qt.MsgBoxIgnore(): return core.ToObject(RI_Ignore), true
            default:                return nil, false
            }
        },
    )
}
func ShowSaveDiscardCancel(content string, title string) core.Observable {
    return showMessageBox(
        qt.MsgBoxWarning(),
        qt.MsgBoxSave().And(qt.MsgBoxDiscard()).And(qt.MsgBoxCancel()),
        qt.MsgBoxCancel(),
        title,
        content,
        func(btn qt.MsgBoxBtn) (core.Object, bool) {
            switch btn {
            case qt.MsgBoxSave():    return core.ToObject(SD_Save), true
            case qt.MsgBoxDiscard(): return core.ToObject(SD_Discard), true
            default:                 return nil, false
            }
        },
    )
}
type RetryIgnore int
const (
    RI_Retry RetryIgnore = iota
    RI_Ignore
)
type SaveDiscard int
const (
    SD_Save SaveDiscard = iota
    SD_Discard
)

func GetChoice(prompt string, items ([] ComboBoxItem), title string, h core.RuntimeHandle) core.Observable {
    if len(items) == 0 {
        core.Crash(h, core.InvalidArgument, "item list cannot be empty")
    }
    return showComboBoxInputBox(items, title, prompt, h)
}
func GetLine(prompt string, initial string, title string) core.Observable {
    return showInputBox(qt.InputText(), false, nil, initial, title, prompt, func(s string, _ int, _ float64) core.Object {
        return core.ObjString(s)
    })
}
func GetText(prompt string, initial string, title string) core.Observable {
    return showInputBox(qt.InputText(), true, nil, initial, title, prompt, func(s string, _ int, _ float64) core.Object {
        return core.ObjString(s)
    })
}
func GetInt(prompt string, initial int, title string) core.Observable {
    return showInputBox(qt.InputInt(), false, nil, initial, title, prompt, func(_ string, i int, _ float64) core.Object {
        return core.ObjInt(i)
    })
}
func GetFloat(prompt string, initial float64, title string) core.Observable {
    return showInputBox(qt.InputDouble(), false, nil, initial, title, prompt, func(_ string, _ int, x float64) core.Object {
        return core.ObjFloat(x)
    })
}

func GetFileListToOpen(filter string) core.Observable {
    return showFileDialog(qt.FileDialogModeOpenMultiple(), filter)
}
func GetFileToOpen(filter string) core.Observable {
    return showFileDialog(qt.FileDialogModeOpenSingle(), filter)
}
func GetFileToSave(filter string) core.Observable {
    return showFileDialog(qt.FileDialogModeSave(), filter)
}
func showFileDialog(mode qt.FileDialogMode, filter string) core.Observable {
    return core.Observable(func(pub core.DataPublisher) {
        var dialog = qt.CreateFileDialog(mode, filter)
        dialog.Consume(func(files ([] string), ok bool) {
            var yield, complete = core.DialogGenerate(pub); defer complete()
            if ok {
                switch mode {
                default:
                    yield(core.ObjFile(files[0]))
                case qt.FileDialogModeOpenMultiple():
                    yield(core.ObjList(ctn.MapEach(files, core.ObjFile)))
                }
            }
        })
    })
}

func Action(icon Icon, text string, shortcut string, repeat bool, enable core.Observable, h core.RuntimeHandle) core.Hook {
    return core.MakeHookWithEffect(h, func() (core.Action, core.Observable, func()) {
        var action, delete_ = core.WrapAction(func(ctx qt.Pkg) qt.Action {
            return qt.CreateAction(convertIcon(icon, h), text, shortcut, repeat, ctx)
        })
        var bind_enable = enable.ConcatMap(func(obj core.Object) core.Observable {
            return core.SetActionEnabled(action, core.GetBool(obj), h)
        })
        return action, bind_enable, delete_
    })
}
func ActionTriggers(a core.Action, h core.RuntimeHandle) core.Observable {
    return core.ConnectActionTrigger(a, h)
}

type ActionCheckBox_T struct { Checked core.Observable }
func ActionCheckBox(a core.Action, initial bool, h core.RuntimeHandle) core.Hook {
    return core.ActionCheckBox(a, initial, h, func(checked core.Observable) core.Object {
        return core.ToObject(ActionCheckBox_T { checked })
    })
}

type ActionComboBox_T struct { SelectedItem core.Observable }
func ActionComboBox(items ([] ActionComboBoxItem), h core.RuntimeHandle) core.Hook {
    if len(items) == 0 {
        core.Crash(h, core.InvalidArgument, "item list cannot be empty")
    }
    var adapted = ctn.MapEach(items, adaptActionComboBoxItem)
    return core.ActionComboBox(adapted, h, func(selected_ core.Observable) core.Object {
        var selected = selected_.Map(func(obj core.Object) core.Object {
            var index = core.GetInt(obj)
            return items[index].Value
        })
        return core.ToObject(ActionComboBox_T { selected })
    })
}

type MenuBar struct {
    Menus  [] Menu
}
func convertMenuBar(b MenuBar, h core.RuntimeHandle) qt.MenuBar {
    if len(b.Menus) == 0 {
        return qt.MenuBar {}
    }
    var menu_bar = qt.CreateMenuBar()
    for _, m := range b.Menus {
        menu_bar.AddMenu(convertMenu(m, h))
    }
    return menu_bar
}

type ToolBar struct {
    Mode   ToolBarMode
    Items  [] ToolBarItem
}
type ToolBarMode int
const (
    TBM_IconOnly ToolBarMode = iota
    TBM_TextOnly
    TBM_TextBesideIcon
    TBM_TextUnderIcon
)
type ToolBarItem struct {
    pseudounion.Tag
    Menu; Action(core.Action); Separator; Widget(core.Widget); Spacer
}
func convertToolBar(b ToolBar, h core.RuntimeHandle) qt.ToolBar {
    if len(b.Items) == 0 {
        return qt.ToolBar {}
    }
    var tool_bar = qt.CreateToolBar(convertToolBarMode(b.Mode))
    for _, item := range b.Items {
        addToolBarItem(item, tool_bar, h)
    }
    return tool_bar
}
func convertToolBarMode(m ToolBarMode) qt.ToolButtonStyle {
    switch m {
    case TBM_IconOnly:       return qt.ToolButtonIconOnly()
    case TBM_TextOnly:       return qt.ToolButtonTextOnly()
    case TBM_TextBesideIcon: return qt.ToolButtonTextBesideIcon()
    case TBM_TextUnderIcon:  return qt.ToolButtonTextUnderIcon()
    default:
        panic("impossible branch")
    }
}
func addToolBarItem(i ToolBarItem, b qt.ToolBar, h core.RuntimeHandle) {
    switch I := pseudounion.Load(i).(type) {
    case Menu:
        b.AddMenu(convertMenu(I, h))
    case core.Action:
        b.AddAction(I.Deref(h))
    case Separator:
        b.AddSeparator()
    case core.Widget:
        b.AddWidget(I.Deref(h))
    case Spacer:
        b.AddSpacer(I.Width, I.Height, I.Expand)
    default:
        panic("impossible branch")
    }
}

func BindContextMenu(w core.Widget, m Menu, h core.RuntimeHandle) core.Observable {
    return core.BindContextMenu(w, convertMenu(m, h), h)
}
type Menu struct {
    Icon   Icon
    Name   string
    Items  [] MenuItem
}
type MenuItem struct {
    pseudounion.Tag
    Menu; Action(core.Action); Separator
}
type Separator struct {}
func convertMenu(m Menu, h core.RuntimeHandle) qt.Menu {
    var menu = qt.CreateMenu(convertIcon(m.Icon, h), m.Name)
    for _, item := range m.Items {
        addMenuItem(item, menu, h)
    }
    return menu
}
func addMenuItem(i MenuItem, m qt.Menu, h core.RuntimeHandle) {
    switch I := pseudounion.Load(i).(type) {
    case Menu:
        m.AddMenu(convertMenu(I, h))
    case core.Action:
        m.AddAction(I.Deref(h))
    case Separator:
        m.AddSeparator()
    default:
        panic("impossible branch")
    }
}

func ShowAndActivate(w core.Widget, h core.RuntimeHandle) core.Observable {
    return core.Observable(func(pub core.DataPublisher) {
        pub.SyncReturn(func() (core.Object, error) {
            var W = w.Deref(h)
            W.Show()
            W.MoveToScreenCenter()
            W.ActivateWindow()
            return nil, nil
        })
    })
}
func BindInlineStyleSheet(w core.Widget, o core.Observable, h core.RuntimeHandle) core.Observable {
    return core.BindInlineStyleSheet(w, o, h)
}
func ComboBoxSelectedItem(w core.Widget, items ([] ComboBoxItem), h core.RuntimeHandle) core.Observable {
    var signal = core.MakeSignalWithValueGetter(qt.ComboBox_CurrentIndexChanged, true, func(W qt.Widget) core.Object {
        return core.ObjInt(int(W.GetPropInt(qt.ComboBox_CurrentIndex)))
    })
    return signal.Connect(w, h).Map(func(obj core.Object) core.Object {
        var index = core.GetInt(obj)
        return items[index].Value
    })
}
func CreateWidget(layout Layout, margin_x int, margin_y int, policy_x SizePolicy, policy_y SizePolicy, h core.RuntimeHandle) core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        return qt.CreateWidget(convertLayout(layout, h), margin_x, margin_y, convertSizePolicy(policy_x), convertSizePolicy(policy_y), ctx)
    })
}
func CreateScrollArea(scroll Scroll, layout Layout, margin_x int, margin_y int, h core.RuntimeHandle) core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        return qt.CreateScrollArea(convertScroll(scroll), convertLayout(layout, h), margin_x, margin_y, ctx).Widget
    })
}
func CreateGroupBox(title string, layout Layout, margin_x int, margin_y int, h core.RuntimeHandle) core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        return qt.CreateGroupBox(title, convertLayout(layout, h), margin_x, margin_y, ctx).Widget
    })
}
func CreateSplitter(list ([] core.Widget), h core.RuntimeHandle) core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        var converted = make([] qt.Widget, len(list))
        for i := range list {
            converted[i] = list[i].Deref(h)
        }
        return qt.CreateSplitter(converted, ctx).Widget
    })
}
func CreateMainWindow(menu_bar MenuBar, tool_bar ToolBar, layout Layout, margin_x int, margin_y int, width int, height int, icon Icon, h core.RuntimeHandle) core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        return qt.CreateMainWindow(convertMenuBar(menu_bar, h), convertToolBar(tool_bar, h), convertLayout(layout, h), margin_x, margin_y, width, height, convertIcon(icon, h), ctx).Widget
    })
}
func CreateDialog(layout Layout, margin_x int, margin_y int, width int, height int, icon Icon, h core.RuntimeHandle) core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        return qt.CreateDialog(convertLayout(layout, h), margin_x, margin_y, width, height, convertIcon(icon, h), ctx).Widget
    })
}
func CreateLabel(align Align, selectable bool) core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        return qt.CreateLabel("", convertAlign(align), selectable, ctx).Widget
    })
}
func CreateIconLabel(icon Icon, size int, h core.RuntimeHandle) core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        return qt.CreateIconLabel(convertIcon(icon, h), size, ctx).Widget
    })
}
func CreateElidedLabel() core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        return qt.CreateElidedLabel("", ctx).Widget
    })
}
func CreateTextView(format TextFormat) core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        return qt.CreateTextView("", convertTextFormat(format), ctx).Widget
    })
}
func CreateCheckBox(text string, checked bool) core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        return qt.CreateCheckBox(text, checked, ctx).Widget
    })
}
func CreateComboBox(items ([] ComboBoxItem), h core.RuntimeHandle) core.Hook {
    if len(items) == 0 {
        core.Crash(h, core.InvalidArgument, "item list cannot be empty")
    }
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        var items = ctn.MapEach(items, func(item ComboBoxItem) qt.ComboBoxItem {
            return convertComboBoxItem(item, h)
        })
        return qt.CreateComboBox(items, ctx).Widget
    })
}
func CreatePushButton(icon Icon, text string, tooltip string, h core.RuntimeHandle) core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        return qt.CreatePushButton(convertIcon(icon, h), text, tooltip, ctx).Widget
    })
}
func CreateLineEdit(text string) core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        return qt.CreateLineEdit(text, ctx).Widget
    })
}
func CreatePlainTextEdit(text string) core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        return qt.CreatePlainTextEdit(text, ctx).Widget
    })
}
func CreateSlider(value int, min int, max int) core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        return qt.CreateSlider(min, max, value, ctx).Widget
    })
}
func CreateProgressBar(max int, format string) core.Hook {
    return createWidget(func(ctx qt.Pkg) qt.Widget {
        return qt.CreateProgressBar(format, max, ctx).Widget
    })
}

func Connect(s core.Signal, w core.Widget, h core.RuntimeHandle) core.Observable {
    return s.Connect(w, h)
}
func Listen(s core.Events, w core.Widget, h core.RuntimeHandle) core.Observable {
    return s.Listen(w, h)
}
func SignalToggled() core.Signal {
    return core.MakeSignal(qt.S_Toggled)
}
func SignalClicked() core.Signal {
    return core.MakeSignal(qt.S_Clicked)
}
func SignalTextChanged0() core.Signal {
    return core.MakeSignal(qt.S_TextChanged0)
}
func SignalTextChanged1() core.Signal {
    return core.MakeSignal(qt.S_TextChanged1)
}
func SignalReturnPressed() core.Signal {
    return core.MakeSignal(qt.S_ReturnPressed)
}
func SignalValueChanged() core.Signal {
    return core.MakeSignal(qt.S_ValueChanged)
}
func EventsShow() core.Events {
    return core.MakeEvents(qt.EventShow(), false)
}
func EventsClose() core.Events {
    return core.MakeEvents(qt.EventClose(), true)
}

func Read(p core.Prop, s core.Signal) core.Signal {
    return p.Read(s)
}
func Bind(p core.Prop, o core.Observable, w core.Widget, h core.RuntimeHandle) core.Observable {
    return p.Bind(o, w, h)
}
func ClearTextLater(w core.Widget, o core.Observable, h core.RuntimeHandle) core.Observable {
    return o.ConcatMap(func(obj core.Object) core.Observable {
        return core.Observable(func(pub core.DataPublisher) {
            pub.SyncReturn(func() (core.Object, error) {
                w.Deref(h).ClearTextLater()
                return obj, nil
            })
        })
    })
}
func PropEnabled() core.Prop {
    return core.MakeProp(core.PropBool, qt.P_Enabled)
}
func PropWindowTitle() core.Prop {
    return core.MakeProp(core.PropString, qt.P_WindowTitle)
}
func PropText() core.Prop {
    return core.MakeProp(core.PropString, qt.P_Text)
}
func PropChecked() core.Prop {
    return core.MakeProp(core.PropBool, qt.P_Checked)
}
func PropPlainText() core.Prop {
    return core.MakeProp(core.PropString, qt.P_PlainText)
}
func PropValue() core.Prop {
    return core.MakeProp(core.PropInt, qt.P_Value)
}

type ListView_T struct {
    Widget     core.Widget
    Extension  core.Observable
    Current    core.Observable
    Selection  core.Observable
}
func ListView(data core.Observable, key func(core.Object)(string), p core.ItemViewProvider, headers ([] HeaderView), stretch int, select_ ItemSelect, h core.RuntimeHandle) core.Hook {
    var config = core.ListViewConfig {
        CreateInterface: func(ctx qt.Pkg) qt.Lwi {
            return qt.LwiCreateFromDefaultListWidget (
                len(headers),
                convertItemSelect(select_),
                getHeaderWidgets(headers, h, ctx),
                stretch,
                ctx,
            )
        },
        ReturnObject: func(w core.Widget, e core.Observable, c core.Observable, s core.Observable) core.Object {
            return core.ToObject(ListView_T {
                Widget:    w,
                Extension: e,
                Current:   c,
                Selection: s,
            })
        },
    }
    return core.ListView(config, data, key, p, h)
}

type ListEditView_T struct {
    Widget     core.Widget
    Output     core.Observable
    Extension  core.Observable
    EditOps    core.Subject
}
func ListEditView(initial core.List, p core.ItemEditViewProvider, headers ([] HeaderView), stretch int, select_ ItemSelect, h core.RuntimeHandle) core.Hook {
    var config = core.ListEditViewConfig {
        CreateInterface: func(ctx qt.Pkg) qt.Lwi {
            return qt.LwiCreateFromDefaultListWidget (
                len(headers),
                convertItemSelect(select_),
                getHeaderWidgets(headers, h, ctx),
                stretch,
                ctx,
            )
        },
        ReturnObject: func(w core.Widget, o core.Observable, e core.Observable, O core.Subject) core.Object {
            return core.ToObject(ListEditView_T {
                Widget:    w,
                Output:    o,
                Extension: e,
                EditOps:   O,
            })
        },
    }
    return core.ListEditView(config, initial, p, h)
}

type Icon struct {
    Name  string
}
func convertIcon(icon Icon, h core.RuntimeHandle) string {
    var name = icon.Name
    if strings.HasPrefix(name, qt.FileIconNamePrefix) {
        var rel_path = strings.TrimPrefix(name, qt.FileIconNamePrefix)
        var abs_path = filepath.Join(filepath.Dir(h.ProgramPath()), rel_path)
        return (qt.FileIconNamePrefix + abs_path)
    } else {
        return name
    }
}

type SizePolicy int
const (
    SP_Rigid SizePolicy = iota
    SP_Controlled
    SP_Incompressible; SP_IncompressibleExpanding
    SP_Free; SP_FreeExpanding
    SP_Bounded
)
func convertSizePolicy(p SizePolicy) qt.SizePolicy {
    switch p {
    case SP_Rigid:                   return qt.SizePolicyRigid()
    case SP_Controlled:              return qt.SizePolicyControlled()
    case SP_Incompressible:          return qt.SizePolicyIncompressible()
    case SP_IncompressibleExpanding: return qt.SizePolicyIncompressibleExpanding()
    case SP_Free:                    return qt.SizePolicyFree()
    case SP_FreeExpanding:           return qt.SizePolicyFreeExpanding()
    case SP_Bounded:                 return qt.SizePolicyBounded()
    default:
        panic("impossible branch")
    }
}

type Align int
const (
    A_Default Align = iota
    A_Center
    A_Left; A_Right; A_Top; A_Bottom
    A_LeftTop; A_LeftBottom; A_RightTop; A_RightBottom
)
func convertAlign(a Align) qt.Alignment {
    switch a {
    case A_Default:     return qt.AlignDefault()
    case A_Center:      return qt.AlignHCenter().And(qt.AlignVCenter())
    case A_Left:        return qt.AlignLeft().And(qt.AlignVCenter())
    case A_Right:       return qt.AlignRight().And(qt.AlignVCenter())
    case A_Top:         return qt.AlignTop().And(qt.AlignHCenter())
    case A_Bottom:      return qt.AlignBottom().And(qt.AlignHCenter())
    case A_LeftTop:     return qt.AlignLeft().And(qt.AlignTop())
    case A_LeftBottom:  return qt.AlignLeft().And(qt.AlignBottom())
    case A_RightTop:    return qt.AlignRight().And(qt.AlignTop())
    case A_RightBottom: return qt.AlignRight().And(qt.AlignBottom())
    default:
        panic("impossible branch")
    }
}

type Scroll int
const (
    S_BothDirection Scroll = iota
    S_VerticalOnly
    S_HorizontalOnly
)
func convertScroll(s Scroll) qt.ScrollDirection {
    switch s {
    case S_BothDirection:  return qt.ScrollBothDirection()
    case S_VerticalOnly:   return qt.ScrollVerticalOnly()
    case S_HorizontalOnly: return qt.ScrollHorizontalOnly()
    default:
        panic("impossible branch")
    }
}

type TextFormat int
const (
    TF_Plain TextFormat = iota
    TF_Html
    TF_Markdown
)
func convertTextFormat(f TextFormat) qt.TextFormat {
    switch f {
    case TF_Plain:    return qt.TextFormatPlain()
    case TF_Html:     return qt.TextFormatHtml()
    case TF_Markdown: return qt.TextFormatMarkdown()
    default:
        panic("impossible branch")
    }
}

type Layout struct {
    pseudounion.Tag
    Row; Column; Grid
}
type Row struct { Items ([] LayoutItem); Spacing int }
type Column struct { Items ([] LayoutItem); Spacing int }
type Grid struct { Spans ([] Span); RowSpacing int; ColumnSpacing int }
type Span struct {
    Item LayoutItem
    Row int; Column int; RowSpan int; ColumnSpan int; Align Align
}
type LayoutItem struct {
    pseudounion.Tag
    Layout; Widget(core.Widget); Spacer; String(string)
}
type Spacer struct { Width int; Height int; Expand bool }
type Wrapper struct { Widget core.Widget }
func convertLayout(l Layout, h core.RuntimeHandle) qt.Layout {
    var q_span qt.GridSpan
    var q_align = qt.AlignDefault()
    switch L := pseudounion.Load(l).(type) {
    case Row:
        var layout = qt.CreateLayoutRow(L.Spacing)
        for _, item := range L.Items {
            addLayoutItem(item, layout, q_span, q_align, h)
        }
        return layout
    case Column:
        var layout = qt.CreateLayoutColumn(L.Spacing)
        for _, item := range L.Items {
            addLayoutItem(item, layout, q_span, q_align, h)
        }
        return layout
    case Grid:
        var layout = qt.CreateLayoutGrid(L.RowSpacing, L.ColumnSpacing)
        for _, span := range L.Spans {
            var item = span.Item
            q_span.Row = span.Row
            q_span.Column = span.Column
            q_span.RowSpan = span.RowSpan
            q_span.ColumnSpan = span.ColumnSpan
            q_align = convertAlign(span.Align)
            addLayoutItem(item, layout, q_span, q_align, h)
        }
        return layout
    default:
        panic("impossible branch")
    }
}
func addLayoutItem(i LayoutItem, l qt.Layout, span qt.GridSpan, align qt.Alignment, h core.RuntimeHandle) {
    switch I := pseudounion.Load(i).(type) {
    case Layout:
        l.AddLayout(convertLayout(I, h), span, align)
    case Spacer:
        l.AddSpacer(I.Width, I.Height, I.Expand, span, align)
    case string:
        l.AddLabel(I, span, align)
    case core.Widget:
        l.AddWidget(I.Deref(h), span, align)
    default:
        panic("impossible branch")
    }
}

type ComboBoxItem struct {
    Icon Icon; Name string; Value core.Object; Selected bool
}
func convertComboBoxItem(item ComboBoxItem, h core.RuntimeHandle) qt.ComboBoxItem {
    return qt.ComboBoxItem {
        Icon:     convertIcon(item.Icon, h),
        Name:     item.Name,
        Selected: item.Selected,
    }
}

type ActionComboBoxItem struct {
    Action core.Action; Value core.Object; Selected bool
}
func adaptActionComboBoxItem(item ActionComboBoxItem) core.ActionComboBoxItem {
    return core.ActionComboBoxItem {
        Action:   item.Action,
        Selected: item.Selected,
    }
}

type HeaderView struct {
    pseudounion.Tag
    String string; Widget core.Widget
}
func getHeaderWidgets(headers ([] HeaderView), h core.RuntimeHandle, ctx qt.Pkg) ([] qt.Widget) {
    return ctn.MapEach(headers, func(v HeaderView) qt.Widget {
        switch V := pseudounion.Load(v).(type) {
        case string:
            var label = qt.CreateLabelLite(V, ctx)
            var spacing = 0
            var margin_x, margin_y = 6, 4
            var policy_x, policy_y = qt.SizePolicyFree(), qt.SizePolicyFree()
            var row = qt.CreateLayoutRow(spacing)
            row.AddWidget(label.Widget, qt.GridSpan{}, qt.AlignDefault())
            return qt.CreateWidget(row, margin_x, margin_y, policy_x, policy_y, ctx)
        case core.Widget:
            return V.Deref(h)
        default:
            panic("impossible branch")
        }
    })
}

type ItemSelect int
const (
    IS_NA ItemSelect = iota
    IS_Single
    IS_Multiple
    IS_MaybeMultiple
)
func convertItemSelect(m ItemSelect) qt.ItemSelectionMode {
    switch m {
    case IS_NA:            return qt.ItemNoSelection()
    case IS_Single:        return qt.ItemSingleSelection()
    case IS_Multiple:      return qt.ItemMultiSelection()
    case IS_MaybeMultiple: return qt.ItemExtendedSelection()
    default:
        panic("impossible branch")
    }
}


func assumeValidRegexp(pattern string) *regexp.Regexp {
    var value, err = regexp.Compile(pattern)
    if err != nil { panic("something went wrong") }
    return value
}
func forEachObservable(l core.List) func(func(core.Observable)) {
    return func(yield func(core.Observable)) {
        l.ForEach(func(item core.Object) {
            yield(core.GetObservable(item))
        })
    }
}
func createWidget(k func(qt.Pkg)(qt.Widget)) core.Hook {
    return core.MakeHook(func() (core.Widget, func()) {
        return core.WrapWidget(k)
    })
}
func showComboBoxInputBox(items ([] ComboBoxItem), title string, prompt string, h core.RuntimeHandle) core.Observable {
    if len(items) == 0 {
        panic("invalid argument")
    }
    return core.Observable(func(pub core.DataPublisher) {
        var items_ = ctn.MapEach(items, func(item ComboBoxItem) qt.ComboBoxItem {
            return convertComboBoxItem(item, h)
        })
        var dialog = qt.CreateComboBoxDialog(items_, title, prompt)
        dialog.Consume(func(index int, ok bool) {
            var yield, complete = core.DialogGenerate(pub); defer complete()
            if ok {
                yield(items[index].Value)
            }
        })
    })
}
func showInputBox(mode qt.InputDialogMode, multilineText bool, choiceItems ([] string), value interface{}, title string, prompt string, f func(string,int,float64)(core.Object)) core.Observable {
    return core.Observable(func(pub core.DataPublisher) {
        var dialog = qt.CreateInputDialog(mode, value, title, prompt)
        if multilineText {
            dialog.UseMultilineText()
        }
        if len(choiceItems) > 0 {
            dialog.UseChoiceItems(choiceItems)
        }
        dialog.Consume(func(s string, i int, x float64, ok bool) {
            var yield, complete = core.DialogGenerate(pub); defer complete()
            if ok {
                yield(f(s,i,x))
            }
        })
    })
}
func showMessageBox(icon qt.MsgBoxIcon, buttons qt.MsgBoxBtn, default_ qt.MsgBoxBtn, title string, content string, f func(qt.MsgBoxBtn)(core.Object,bool)) core.Observable {
    return core.Observable(func(pub core.DataPublisher) {
        var msgbox = qt.CreateMessageBox(icon, buttons, title, content)
        msgbox.SetDefaultButton(default_)
        msgbox.Consume(func(btn qt.MsgBoxBtn) {
            var yield, complete = core.DialogGenerate(pub); defer complete()
            if obj, ok := f(btn); ok {
                yield(obj)
            }
        })
    })
}


