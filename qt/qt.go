package qt

/*
#include <stdlib.h>
#include "qtbinding/qtbinding.h"
*/
// #cgo LDFLAGS: -L./build -L./build/release -lqtbinding -Wl,-rpath=\$ORIGIN/
import "C"

import (
    _ "embed"
    "unsafe"
    "reflect"
    "strings"
)


type Object struct { ptr unsafe.Pointer }

func Init() {
    C.QtInit()
}
func Main() {
    C.QtMain()
}
func Exit(code int) {
    C.QtExit(C.int(code))
}
func UUID() string {
    var s = consumeString(C.QtNewUUID())
    s = strings.TrimPrefix(s, "{")
    s = strings.TrimSuffix(s, "}")
    return s
}
func FontSize() int {
    return int(C.QtFontSize())
}
func Schedule(k func()) {
    var del_cb func()
    cb, del_cb := cgo_callback(func() {
        k()
        del_cb()
    })
    C.QtSchedule(cgo_callback_caller, cb)
}
func Connect(obj Object, signal string, ctx Pkg, k func()) {
    var pkg, dispose = CreatePkg(); defer dispose()
    var cb, del_cb = cgo_callback(k)
    var cb_obj_ptr = C.QtConnect(
        obj.ptr, str(signal, pkg),
        cgo_callback_caller, cb,
    )
    if cb_obj_ptr != nil {
        ctx.push(func() {
            C.QtDeleteObjectLater(cb_obj_ptr)
            del_cb()
        })
    } else {
        ctx.push(func() {})
        println("qt.Connect(): no such signal:", obj.ClassName(), "::", signal)
    }
}
func Listen(obj Object, kind EventKind, prevent bool, ctx Pkg, k func(Event)) {
    var l C.QtEventListener
    var cb, del_cb = cgo_callback(func() {
        var ev = C.QtGetCurrentEvent(l)
        k(Event{ev})
    })
    l = C.QtListen(
        obj.ptr, C.int(kind), makeBool(prevent),
        cgo_callback_caller, cb,
    )
    ctx.push(func() {
        C.QtUnlisten(obj.ptr, l)
        del_cb()
    })
}
const DP_InlineStyleSheetId = "qtbindingInlineStyleSheetId"
const P_StyleSheet = "styleSheet"
const P_Checkable = "checkable"
const P_Checked = "checked"
const P_Enabled = "enabled"
const P_WindowTitle = "windowTitle"
const P_Text = "text"
const P_PlainText = "plainText"
const P_Value = "value"
const P_CurrentIndex = "currentIndex"
const S_Triggered = "triggered()"
const S_Toggled = "toggled(bool)"
const S_Clicked = "clicked(bool)"
const S_Finished = "finished(int)"
const S_TextChanged0 = "textChanged()"
const S_TextChanged1 = "textChanged(const QString&)"
const S_ReturnPressed = "returnPressed()"
const S_ValueChanged = "valueChanged(int)"
const S_CurrentIndexChanged = "currentIndexChanged(int)"
type Event struct { QtEvent C.QtEvent }
type EventKind C.int
func EventMove() EventKind { return EventKind(C.QtEventMove) }
func EventResize() EventKind { return EventKind(C.QtEventResize) }
func EventShow() EventKind { return EventKind(C.QtEventShow) }
func EventClose() EventKind { return EventKind(C.QtEventClose) }
func EventFocusIn() EventKind { return EventKind(C.QtEventFocusIn) }
func EventFocusOut() EventKind { return EventKind(C.QtEventFocusOut) }
func EventWindowActivate() EventKind { return EventKind(C.QtEventWindowActivate) }
func EventWindowDeactivate() EventKind { return EventKind(C.QtEventWindowDeactivate) }
func EventDynamicPropertyChange() EventKind { return EventKind(C.QtEventDynamicPropertyChange) }
func (ev Event) ResizeEventGetWidth() int {
    return int(C.QtResizeEventGetWidth(ev.QtEvent))
}
func (ev Event) ResizeEventGetHeight() int {
    return int(C.QtResizeEventGetHeight(ev.QtEvent))
}
func (ev Event) DynamicPropertyChangeEventGetPropertyName() string {
    return consumeString(C.QtDynamicPropertyChangeEventGetPropertyName(ev.QtEvent))
}
func (obj Object) ClassName() string {
    return consumeString(C.QtObjectGetClassName(obj.ptr))
}
func (obj Object) GetPropString(prop string) string {
    var pkg, dispose = CreatePkg(); defer dispose()
    return consumeString(C.QtObjectGetPropString(obj.ptr, str(prop, pkg)))
}
func (obj Object) SetPropString(prop string, val string) {
    var pkg, dispose = CreatePkg(); defer dispose()
    C.QtObjectSetPropString(obj.ptr, str(prop, pkg), makeString(val, pkg))
}
func (obj Object) GetPropBool(prop string) bool {
    var pkg, dispose = CreatePkg(); defer dispose()
    return getBool(C.QtObjectGetPropBool(obj.ptr, str(prop, pkg)))
}
func (obj Object) SetPropBool(prop string, val bool) {
    var pkg, dispose = CreatePkg(); defer dispose()
    C.QtObjectSetPropBool(obj.ptr, str(prop, pkg), makeBool(val))
}
func (obj Object) GetPropInt(prop string) int {
    var pkg, dispose = CreatePkg(); defer dispose()
    return int(C.QtObjectGetPropInt(obj.ptr, str(prop, pkg)))
}
func (obj Object) SetPropInt(prop string, val int) {
    var pkg, dispose = CreatePkg(); defer dispose()
    C.QtObjectSetPropInt(obj.ptr, str(prop, pkg), C.int(val))
}
func (obj Object) GetPropFloat(prop string) float64 {
    var pkg, dispose = CreatePkg(); defer dispose()
    return float64(C.QtObjectGetPropDouble(obj.ptr, str(prop, pkg)))
}
func (obj Object) SetPropFloat(prop string, val float64) {
    var pkg, dispose = CreatePkg(); defer dispose()
    C.QtObjectSetPropDouble(obj.ptr, str(prop, pkg), C.double(val))
}

type Action struct {
    Object
}
func CreateAction(icon string, text string, shortcut string, repeat bool, ctx Pkg) Action {
    var pkg, dispose = CreatePkg(); defer dispose()
    var shortcuts = strings.Split(shortcut, " ")
    var shortcuts_ptr, shortcuts_len = makeStrings(shortcuts, pkg)
    var ptr = C.QtCreateAction (
        makeIcon(icon, pkg),
        makeString(text, pkg),
        shortcuts_ptr,
        shortcuts_len,
        makeBool(repeat),
    )
    pushObjectDeletion(ctx, ptr)
    return Action{Object{ptr}}
}
func (a Action) Checkable() bool {
    return a.GetPropBool(P_Checkable)
}
func (a Action) SetCheckable(p bool) {
    a.SetPropBool(P_Checkable, p)
}
func (a Action) Checked() bool {
    return a.GetPropBool(P_Checked)
}
func (a Action) SetChecked(p bool) {
    a.SetPropBool(P_Checked, p)
}
func (a Action) SetEnabled(p bool) {
    a.SetPropBool(P_Enabled, p)
}
func (a Action) OnTrigger(ctx Pkg, cb func()) {
    Connect(a.Object, S_Triggered, ctx, cb)
}
func (a Action) OnCheck(ctx Pkg, cb func(bool)) {
    var k = func() {
        cb(a.Checked())
    }
    k()
    a.OnTrigger(ctx, k)
}

type ActionGroup struct {
    Object
}
func CreateActionGroup(ctx Pkg, items ([] ActionGroupItem), ok *bool) ActionGroup {
    var ptr = C.QtCreateActionGroup()
    pushObjectDeletion(ctx, ptr)
    setupActionGroup(ptr, items, ok)
    return ActionGroup{Object{ptr}}
}
type ActionGroupItem struct { Action Action; Selected bool }
func setupActionGroup(g unsafe.Pointer, items ([] ActionGroupItem), ok *bool) {
    var none_selected = true
    *ok = true
    var ok0 = true
    for i, item := range items {
        if getBool(C.QtActionInGroup(item.Action.ptr)) {
            *ok = false
            if i == 0 { ok0 = false }
            continue
        }
        if item.Selected {
            none_selected = false
            item.Action.SetChecked(true)
        } else {
            item.Action.SetChecked(false)
        }
        C.QtActionGroupAddAction(g, item.Action.ptr, C.int(i))
    }
    if none_selected {
        if ok0 {
            items[0].Action.SetChecked(true)
        }
    }
}
func (g ActionGroup) OnTrigger(ctx Pkg, cb func()) {
    Connect(g.Object, "triggered(QAction*)", ctx, cb)
}
func (g ActionGroup) OnCheck(ctx Pkg, cb func(int)) {
    var k = func() {
        var index = int(C.QtActionGroupGetCheckedActionIndex(g.ptr))
        if index >= 0 {
            cb(index)
        }
    }
    k()
    g.OnTrigger(ctx, k)
}

type Menu struct {
    // actually a Widget, but better to assume it's not
    Object
}
func CreateMenu(icon string, text string) Menu {
    var pkg, dispose = CreatePkg(); defer dispose()
    var m = C.QtCreateMenu(makeIcon(icon, pkg), makeString(text, pkg))
    return Menu{Object{m}}
}
func (m Menu) AddMenu(another Menu) {
    C.QtMenuAddMenu(m.ptr, another.ptr)
}
func (m Menu) AddAction(a Action) {
    C.QtMenuAddAction(m.ptr, a.ptr)
}
func (m Menu) AddSeparator() {
    C.QtMenuAddSeparator(m.ptr)
}

type ContextMenu struct {
    Menu
}
func CreateContextMenu(m Menu, ctx Pkg) ContextMenu {
    var ptr = m.ptr
    pushObjectDeletion(ctx, ptr)
    return ContextMenu{Menu{Object{ptr}}}
}
func (m ContextMenu) Bind(w Widget, ctx Pkg) bool {
    var binding_obj_ptr = C.QtBindContextMenu(w.ptr, m.ptr)
    if binding_obj_ptr != nil {
        pushObjectDeletion(ctx, binding_obj_ptr)
        return true
    } else {
        return false
    }
}

type Widget struct {
    Object
}
func CreateWidget(layout Layout, margin_x int, margin_y int, policy_x SizePolicy, policy_y SizePolicy, ctx Pkg) Widget {
    var ptr = C.QtCreateWidget(layout.ptr, C.int(margin_x), C.int(margin_y), C.int(policy_x), C.int(policy_y))
    pushObjectDeletion(ctx, ptr)
    return Widget{Object{ptr}}
}
func pushObjectDeletion(ctx Pkg, ptr unsafe.Pointer) {
    ctx.push(func() {
        C.QtDeleteObjectLater(ptr)
    })
}
func (w Widget) PointerNumber() uintptr {
    return uintptr(w.ptr)
}
func (w Widget) Show() {
    C.QtWidgetShow(w.ptr)
}
func (w Widget) Raise() {
    C.QtWidgetRaise(w.ptr)
}
func (w Widget) ActivateWindow() {
    C.QtWidgetActivateWindow(w.ptr)
}
func (w Widget) MoveToScreenCenter() {
    C.QtWidgetMoveToScreenCenter(w.ptr)
}
func (w Widget) SetEnabled(value bool) {
    w.SetPropBool(P_Enabled, value)
}
func (w Widget) SetStyleSheet(value string) {
    w.SetPropString(P_StyleSheet, value)
}
func (w Widget) OnFocusOut(ctx Pkg, cb func()) {
    Listen(w.Object, EventFocusOut(), false, ctx, func(_ Event) {
        cb()
    })
}
func (w Widget) ClearTextLater() {
    C.QtWidgetClearTextLater(w.ptr)
}
type SizePolicy C.int
func SizePolicyRigid() SizePolicy { return SizePolicy(C.QtSizePolicyRigid) }
func SizePolicyControlled() SizePolicy { return SizePolicy(C.QtSizePolicyControlled) }
func SizePolicyIncompressible() SizePolicy { return SizePolicy(C.QtSizePolicyIncompressible) }
func SizePolicyIncompressibleExpanding() SizePolicy { return SizePolicy(C.QtSizePolicyIncompressibleExpanding) }
func SizePolicyFree() SizePolicy { return SizePolicy(C.QtSizePolicyFree) }
func SizePolicyFreeExpanding() SizePolicy { return SizePolicy(C.QtSizePolicyFreeExpanding) }
func SizePolicyBounded() SizePolicy { return SizePolicy(C.QtSizePolicyBounded) }

type Layout struct {
    Object
}
func CreateLayoutRow(spacing int) Layout {
    var l = C.QtCreateLayoutRow(C.int(spacing))
    return Layout{Object{l}}
}
func CreateLayoutColumn(spacing int) Layout {
    var l = C.QtCreateLayoutColumn(C.int(spacing))
    return Layout{Object{l}}
}
func CreateLayoutGrid(rowSpacing, columnSpacing int) Layout {
    var l = C.QtCreateLayoutGrid(C.int(rowSpacing), C.int(columnSpacing))
    return Layout{Object{l}}
}
func (l Layout) AddLayout(another Layout, span GridSpan, align Alignment) {
    var span_q = C.QtMakeGridSpan(C.int(span.Row), C.int(span.Column), C.int(span.RowSpan), C.int(span.ColumnSpan))
    var align_q = C.int(align)
    C.QtLayoutAddLayout(l.ptr, another.ptr, span_q, align_q)
}
func (l Layout) AddWidget(w Widget, span GridSpan, align Alignment) {
    var span_q = C.QtMakeGridSpan(C.int(span.Row), C.int(span.Column), C.int(span.RowSpan), C.int(span.ColumnSpan))
    var align_q = C.int(align)
    C.QtLayoutAddWidget(l.ptr, w.ptr, span_q, align_q)
}
func (l Layout) AddSpacer(width int, height int, expand bool, span GridSpan, align Alignment) {
    var span_q = C.QtMakeGridSpan(C.int(span.Row), C.int(span.Column), C.int(span.RowSpan), C.int(span.ColumnSpan))
    var align_q = C.int(align)
    C.QtLayoutAddSpacer(l.ptr, C.int(width), C.int(height), makeBool(expand), span_q, align_q)
}
func (l Layout) AddLabel(text string, span GridSpan, align Alignment) {
    var pkg, dispose = CreatePkg(); defer dispose()
    var span_q = C.QtMakeGridSpan(C.int(span.Row), C.int(span.Column), C.int(span.RowSpan), C.int(span.ColumnSpan))
    var align_q = C.int(align)
    C.QtLayoutAddLabel(l.ptr, makeString(text, pkg), span_q, align_q)
}
type GridSpan struct {
    Row int; Column int; RowSpan int; ColumnSpan int
}
type Alignment C.int
func (align Alignment) And(another Alignment) Alignment {
    return Alignment(C.int(align) | C.int(another))
}
func AlignDefault() Alignment { return Alignment(C.QtAlignDefault) }
func AlignLeft() Alignment { return Alignment(C.QtAlignLeft) }
func AlignRight() Alignment { return Alignment(C.QtAlignRight) }
func AlignHCenter() Alignment { return Alignment(C.QtAlignHCenter) }
func AlignTop() Alignment { return Alignment(C.QtAlignTop) }
func AlignBottom() Alignment { return Alignment(C.QtAlignBottom) }
func AlignVCenter() Alignment { return Alignment(C.QtAlignVCenter) }

type MainWindow struct {
    Widget
}
func CreateMainWindow(menu_bar MenuBar, tool_bar ToolBar, layout Layout, margin_x int, margin_y int, width int, height int, icon string, ctx Pkg) MainWindow {
    var pkg, dispose = CreatePkg(); defer dispose()
    var ptr = C.QtCreateMainWindow(menu_bar.ptr, tool_bar.ptr, layout.ptr, C.int(margin_x), C.int(margin_y), C.int(width), C.int(height), makeIcon(icon, pkg))
    pushObjectDeletion(ctx, ptr)
    return MainWindow{Widget{Object{ptr}}}
}
func (w MainWindow) OnClose(ctx Pkg, cb func()) {
    Listen(w.Object, EventClose(), true, ctx, func(_ Event) {
        cb()
    })
}

type MenuBar struct {
    Widget
}
func CreateMenuBar() MenuBar {
    var b = C.QtCreateMenuBar()
    return MenuBar{Widget{Object{b}}}
}
func (b MenuBar) AddMenu(m Menu) {
    C.QtMenuBarAddMenu(b.ptr, m.ptr)
}

type ToolBar struct {
    Widget
}
func CreateToolBar(tool_btn_style ToolButtonStyle) ToolBar {
    var b = C.QtCreateToolBar(C.int(tool_btn_style))
    return ToolBar{Widget{Object{b}}}
}
func (b ToolBar) AddMenu(m Menu) {
    C.QtToolBarAddMenu(b.ptr, m.ptr)
}
func (b ToolBar) AddAction(a Action) {
    C.QtToolBarAddAction(b.ptr, a.ptr)
}
func (b ToolBar) AddSeparator() {
    C.QtToolBarAddSeparator(b.ptr)
}
func (b ToolBar) AddWidget(w Widget) {
    C.QtToolBarAddWidget(b.ptr, w.ptr)
}
func (b ToolBar) AddSpacer(width int, height int, expand bool) {
    C.QtToolBarAddSpacer(b.ptr, C.int(width), C.int(height), makeBool(expand))
}
type ToolButtonStyle int
func ToolButtonIconOnly() ToolButtonStyle { return ToolButtonStyle(C.QtToolButtonIconOnly) }
func ToolButtonTextOnly() ToolButtonStyle { return ToolButtonStyle(C.QtToolButtonTextOnly) }
func ToolButtonTextBesideIcon() ToolButtonStyle { return ToolButtonStyle(C.QtToolButtonTextBesideIcon) }
func ToolButtonTextUnderIcon() ToolButtonStyle { return ToolButtonStyle(C.QtToolButtonTextUnderIcon) }

type ScrollArea struct {
    Widget
}
func CreateScrollArea(direction ScrollDirection, layout Layout, margin_x int, margin_y int, ctx Pkg) ScrollArea {
    var ptr = C.QtCreateScrollArea(C.int(direction), layout.ptr, C.int(margin_x), C.int(margin_y))
    pushObjectDeletion(ctx, ptr)
    return ScrollArea{Widget{Object{ptr}}}
}
type ScrollDirection C.int
func ScrollBothDirection() ScrollDirection { return ScrollDirection(C.QtScrollBothDirection) }
func ScrollVerticalOnly() ScrollDirection { return ScrollDirection(C.QtScrollVerticalOnly) }
func ScrollHorizontalOnly() ScrollDirection { return ScrollDirection(C.QtScrollHorizontalOnly) }

type GroupBox struct {
    Widget
}
func CreateGroupBox(title string, layout Layout, margin_x int, margin_y int, ctx Pkg) GroupBox {
    var pkg, dispose = CreatePkg(); defer dispose()
    var ptr = C.QtCreateGroupBox(makeString(title, pkg), layout.ptr, C.int(margin_x), C.int(margin_y))
    pushObjectDeletion(ctx, ptr)
    return GroupBox{Widget{Object{ptr}}}
}

type Splitter struct {
    Widget
}
func CreateSplitter(widgets ([] Widget), ctx Pkg) Splitter {
    var widgets_ptr, widgets_len = ptrlen(widgets)
    var ptr = C.QtCreateSplitter(widgets_ptr, widgets_len)
    pushObjectDeletion(ctx, ptr)
    return Splitter{Widget{Object{ptr}}}
}

type Dialog struct {
    Widget
}
func CreateDialog(layout Layout, margin_x int, margin_y int, width int, height int, icon string, ctx Pkg) Dialog {
    var pkg, dispose = CreatePkg(); defer dispose()
    var ptr = C.QtCreateDialog(layout.ptr, C.int(margin_x), C.int(margin_y), C.int(width), C.int(height), makeIcon(icon, pkg))
    pushObjectDeletion(ctx, ptr)
    return Dialog{Widget{Object{ptr}}}
}
func (w Dialog) Accept() {
    C.QtDialogAccept(w.ptr)
}
func (w Dialog) Reject() {
    C.QtDialogReject(w.ptr)
}
func (w Dialog) OnFinish(ctx Pkg, cb func(bool)) {
    Connect(w.Object, S_Finished, ctx, func() {
        var result = int(C.QtDialogGetResult(w.ptr))
        var ok = (result != 0)
        cb(ok)
    })
}

type DialogButtonBox struct {
    Widget
}
func CreateDialogButtonBox(ctx Pkg) DialogButtonBox {
    var ptr = C.QtCreateDialogButtonBox()
    pushObjectDeletion(ctx, ptr)
    return DialogButtonBox{Widget{Object{ptr}}}
}
func (w DialogButtonBox) AddButton(kind BoxBtn) PushButton {
    var b = C.QtDialogButtonBoxAddButton(w.ptr, C.int(kind))
    return PushButton{Widget{Object{b}}}
}
type BoxBtn C.int
func BtnBoxOK() BoxBtn { return BoxBtn(C.QtBtnBoxOK) }
func BtnBoxCancel() BoxBtn { return BoxBtn(C.QtBtnBoxCancel) }

type DynamicWidget struct {
    Widget
}
func CreateDynamicWidget(ctx Pkg) DynamicWidget {
    var ptr = C.QtCreateDynamicWidget()
    pushObjectDeletion(ctx, ptr)
    return DynamicWidget{Widget{Object{ptr}}}
}
func (w DynamicWidget) SetWidget(new_widget Widget) {
    C.QtDynamicWidgetSetWidget(w.ptr, new_widget.ptr)
}

type Label struct {
    Widget
}
func CreateLabel(text string, align Alignment, selectable bool, ctx Pkg) Label {
    var pkg, dispose = CreatePkg(); defer dispose()
    var ptr = C.QtCreateLabel(makeString(text, pkg), C.int(align), makeBool(selectable))
    pushObjectDeletion(ctx, ptr)
    return Label{Widget{Object{ptr}}}
}
func CreateLabelLite(text string, ctx Pkg) Label {
    return CreateLabel(text, AlignLeft(), false, ctx)
}
func (w Label) SetText(content string) {
    w.SetPropString(P_Text, content)
}

// type IconLabel struct {
//     Label
// }
func CreateIconLabel(icon string, size int, ctx Pkg) Label {
    var pkg, dispose = CreatePkg(); defer dispose()
    var ptr = C.QtCreateIconLabel(makeIcon(icon, pkg), C.int(size))
    pushObjectDeletion(ctx, ptr)
    return Label{Widget{Object{ptr}}}
}

type ElidedLabel struct {
    Widget
}
func CreateElidedLabel(text string, ctx Pkg) ElidedLabel {
    var pkg, dispose = CreatePkg(); defer dispose()
    var ptr = C.QtCreateElidedLabel(makeString(text, pkg))
    pushObjectDeletion(ctx, ptr)
    return ElidedLabel{Widget{Object{ptr}}}
}
func (w ElidedLabel) SetText(content string) {
    w.SetPropString(P_Text, content)
}

type TextView struct {
    Widget
}
func CreateTextView(text string, format TextFormat, ctx Pkg) TextView {
    var pkg, dispose = CreatePkg(); defer dispose()
    var ptr = C.QtCreateTextView(makeString(text, pkg), C.int(format))
    pushObjectDeletion(ctx, ptr)
    return TextView{Widget{Object{ptr}}}
}
type TextFormat C.int
func TextFormatPlain() TextFormat { return TextFormat(C.QtTextPlain) }
func TextFormatHtml() TextFormat { return TextFormat(C.QtTextHtml) }
func TextFormatMarkdown() TextFormat { return TextFormat(C.QtTextMarkdown) }

type CheckBox struct {
    Widget
}
func CreateCheckBox(text string, checked bool, ctx Pkg) CheckBox {
    var pkg, dispose = CreatePkg(); defer dispose()
    var ptr = C.QtCreateCheckBox(makeString(text, pkg), makeBool(checked))
    pushObjectDeletion(ctx, ptr)
    return CheckBox{Widget{Object{ptr}}}
}

type ComboBox struct {
    Widget
}
func CreateComboBox(items ([] ComboBoxItem), ctx Pkg) ComboBox {
    var pkg, dispose = CreatePkg(); defer dispose()
    var ptr = C.QtCreateComboBox()
    pushObjectDeletion(ctx, ptr)
    setupComboBox(ptr, items, pkg)
    return ComboBox{Widget{Object{ptr}}}
}
type ComboBoxItem struct { Icon string; Name string; Selected bool }
func setupComboBox(b unsafe.Pointer, items ([] ComboBoxItem), ctx Pkg) {
    if len(items) == 0 {
        panic("invalid argument")
    }
    var none_selected = true
    for _, item := range items {
        C.QtComboBoxAddItem(b, makeIcon(item.Icon, ctx), makeString(item.Name, ctx), makeBool(item.Selected))
        if item.Selected {
            none_selected = false
        }
    }
    if none_selected {
        Widget{Object{b}}.SetPropInt(ComboBox_CurrentIndex, 0)
    }
}
const ComboBox_CurrentIndex = P_CurrentIndex
const ComboBox_CurrentIndexChanged = S_CurrentIndexChanged

type ComboBoxDialog struct {
    Dialog
}
func CreateComboBoxDialog(items ([] ComboBoxItem), title string, prompt string) ComboBoxDialog {
    var pkg, dispose = CreatePkg(); defer dispose()
    var d = C.QtCreateComboBoxDialog(makeString(title, pkg), makeString(prompt, pkg))
    var b = C.QtComboBoxDialogGetComboBox(d)
    setupComboBox(b, items, pkg)
    return ComboBoxDialog{Dialog{Widget{Object{d}}}}
}
func (dialog ComboBoxDialog) ComboBox() ComboBox {
    var b = C.QtComboBoxDialogGetComboBox(dialog.ptr)
    return ComboBox{Widget{Object{b}}}
}
func (dialog ComboBoxDialog) Consume(k func(int,bool)) {
    var del_cb func()
    cb, del_cb := cgo_callback(func() {
        if getBool(C.QtDialogGetResultBoolean(dialog.ptr)) {
            var combo_box = dialog.ComboBox()
            k(int(combo_box.GetPropInt(ComboBox_CurrentIndex)), true)
        } else {
            k(-1, false)
        }
        del_cb()
    })
    C.QtConsumeDialog(dialog.ptr, cgo_callback_caller, cb)
}

type PushButton struct {
    Widget
}
func CreatePushButton(icon string, text string, tooltip string, ctx Pkg) PushButton {
    var pkg, dispose = CreatePkg(); defer dispose()
    var ptr = C.QtCreatePushButton(makeIcon(icon, pkg), makeString(text, pkg), makeString(tooltip, pkg))
    pushObjectDeletion(ctx, ptr)
    return PushButton{Widget{Object{ptr}}}
}
func (w PushButton) OnClick(ctx Pkg, cb func()) {
    Connect(w.Object, PushButton_Clicked, ctx, cb)
}
const PushButton_Clicked = S_Clicked

type LineEdit struct {
    Widget
}
func CreateLineEdit(text string, ctx Pkg) LineEdit {
    var pkg, dispose = CreatePkg(); defer dispose()
    var ptr = C.QtCreateLineEdit(makeString(text, pkg))
    pushObjectDeletion(ctx, ptr)
    return LineEdit{Widget{Object{ptr}}}
}
func (w LineEdit) SetText(content string) {
    w.SetPropString(P_Text, content)
}
func (w LineEdit) Text() string {
    return w.GetPropString(P_Text)
}
func (w LineEdit) OnTextChange(ctx Pkg, cb func()) {
    Connect(w.Object, LineEdit_TextChanged, ctx, cb)
}
const LineEdit_TextChanged = S_TextChanged1

type TextEdit struct {
    Widget
}
func (w TextEdit) SetPlainText(content string) {
    w.SetPropString(P_PlainText, content)
}
func (w TextEdit) PlainText() string {
    return w.GetPropString(P_PlainText)
}
func (w TextEdit) OnTextChange(ctx Pkg, cb func()) {
    Connect(w.Object, TextEdit_TextChanged, ctx, cb)
}
const TextEdit_TextChanged = S_TextChanged0

type PlainTextEdit struct {
    Widget
}
func CreatePlainTextEdit(text string, ctx Pkg) PlainTextEdit {
    var pkg, dispose = CreatePkg(); defer dispose()
    var ptr = C.QtCreatePlainTextEdit(makeString(text, pkg))
    pushObjectDeletion(ctx, ptr)
    return PlainTextEdit{Widget{Object{ptr}}}
}
func (w PlainTextEdit) SetText(content string) {
    w.SetPropString(P_PlainText, content)
}
func (w PlainTextEdit) Text() string {
    return w.GetPropString(P_PlainText)
}
func (w PlainTextEdit) OnTextChange(ctx Pkg, cb func()) {
    Connect(w.Object, PlainTextEdit_TextChanged, ctx, cb)
}
const PlainTextEdit_TextChanged = S_TextChanged0

type Slider struct {
    Widget
}
func CreateSlider(min int, max int, value int, ctx Pkg) Slider {
    var ptr = C.QtCreateSlider(C.int(min), C.int(max), C.int(value))
    pushObjectDeletion(ctx, ptr)
    return Slider{Widget{Object{ptr}}}
}

type ProgressBar struct {
    Widget
}
func CreateProgressBar(format string, max int, ctx Pkg) ProgressBar {
    var pkg, dispose = CreatePkg(); defer dispose()
    var ptr = C.QtCreateProgressBar(makeString(format, pkg), C.int(max))
    pushObjectDeletion(ctx, ptr)
    return ProgressBar{Widget{Object{ptr}}}
}

type DummyFocusableWidget struct {
    Widget
}
func CreateDummyFocusableWidget(ctx Pkg) DummyFocusableWidget {
    var ptr = C.QtCreateDummyFocusableWidget()
    pushObjectDeletion(ctx, ptr)
    return DummyFocusableWidget{Widget{Object{ptr}}}
}

type InputDialog struct {
    Dialog
}
func CreateInputDialog(mode InputDialogMode, value interface{}, title string, content string) InputDialog {
    var pkg, dispose = CreatePkg(); defer dispose()
    var d = C.QtCreateInputDialog (
        C.int(mode),
        makeVariant(value, pkg),
        makeString(title, pkg),
        makeString(content, pkg),
    )
    return InputDialog{Dialog{Widget{Object{d}}}}
}
func (dialog InputDialog) UseMultilineText() {
    C.QtInputDialogUseMultilineText(dialog.ptr)
}
func (dialog InputDialog) UseChoiceItems(items ([] string)) {
    var pkg, dispose = CreatePkg(); defer dispose()
    var items_array, items_size = makeStrings(items, pkg)
    C.QtInputDialogUseChoiceItems(dialog.ptr, items_array, items_size)
}
func (dialog InputDialog) Consume(k func(string,int,float64,bool)) {
    var del_cb func()
    cb, del_cb := cgo_callback(func() {
        if getBool(C.QtDialogGetResultBoolean(dialog.ptr)) {
            var v_text = consumeString(C.QtInputDialogGetTextValue(dialog.ptr))
            var v_int = int(C.QtInputDialogGetIntValue(dialog.ptr))
            var v_double = float64(C.QtInputDialogGetDoubleValue(dialog.ptr))
            k(v_text, v_int, v_double, true)
        } else {
            k("", 0, 0.0, false)
        }
        del_cb()
    })
    C.QtConsumeDialog(dialog.ptr, cgo_callback_caller, cb)
}
type InputDialogMode C.int
func InputText() InputDialogMode { return InputDialogMode(C.QtInputText) }
func InputInt() InputDialogMode { return InputDialogMode(C.QtInputInt) }
func InputDouble() InputDialogMode { return InputDialogMode(C.QtInputDouble) }

type MessageBox struct {
    Dialog
}
func CreateMessageBox(icon MsgBoxIcon, buttons MsgBoxBtn, title string, content string) MessageBox {
    var pkg, dispose = CreatePkg(); defer dispose()
    var d = C.QtCreateMessageBox(
        C.int(icon), C.int(buttons),
        makeString(title, pkg), makeString(content, pkg),
    )
    return MessageBox{Dialog{Widget{Object{d}}}}
}
func (msgbox MessageBox) SetDefaultButton(btn MsgBoxBtn) {
    C.QtMessageBoxSetDefaultButton(msgbox.ptr, C.int(btn))
}
func (msgbox MessageBox) Consume(k func(MsgBoxBtn)) {
    var del_cb func()
    cb, del_cb := cgo_callback(func() {
        var btn = MsgBoxBtn(int(C.QtMessageBoxGetResultButton(msgbox.ptr)))
        k(btn)
        del_cb()
    })
    C.QtConsumeDialog(msgbox.ptr, cgo_callback_caller, cb)
}
type MsgBoxIcon C.int
func MsgBoxInfo() MsgBoxIcon { return MsgBoxIcon(C.QtMsgBoxInfo) }
func MsgBoxWarning() MsgBoxIcon { return MsgBoxIcon(C.QtMsgBoxWarning) }
func MsgBoxCritical() MsgBoxIcon { return MsgBoxIcon(C.QtMsgBoxCritical) }
func MsgBoxQuestion() MsgBoxIcon { return MsgBoxIcon(C.QtMsgBoxQuestion) }
type MsgBoxBtn C.int
func MsgBoxOK() MsgBoxBtn { return MsgBoxBtn(C.QtMsgBoxOK) }
func MsgBoxCancel() MsgBoxBtn { return MsgBoxBtn(C.QtMsgBoxCancel) }
func MsgBoxYes() MsgBoxBtn { return MsgBoxBtn(C.QtMsgBoxYes) }
func MsgBoxNo() MsgBoxBtn { return MsgBoxBtn(C.QtMsgBoxNo) }
func MsgBoxAbort() MsgBoxBtn { return MsgBoxBtn(C.QtMsgBoxAbort) }
func MsgBoxRetry() MsgBoxBtn { return MsgBoxBtn(C.QtMsgBoxRetry) }
func MsgBoxIgnore() MsgBoxBtn { return MsgBoxBtn(C.QtMsgBoxIgnore) }
func MsgBoxSave() MsgBoxBtn { return MsgBoxBtn(C.QtMsgBoxSave) }
func MsgBoxDiscard() MsgBoxBtn { return MsgBoxBtn(C.QtMsgBoxDiscard) }
func (btn MsgBoxBtn) And(another MsgBoxBtn) MsgBoxBtn {
    return MsgBoxBtn(C.int(btn) | C.int(another))
}

type FileDialog struct {
    Dialog
}
func CreateFileDialog(mode FileDialogMode, filters string) FileDialog {
    var pkg, dispose = CreatePkg(); defer dispose()
    var d = C.QtCreateFileDialog(C.int(mode), makeString(filters, pkg))
    return FileDialog{Dialog{Widget{Object{d}}}}
}
func (dialog FileDialog) Consume(k func([] string, bool)) {
    var del_cb func()
    cb, del_cb := cgo_callback(func() {
        if getBool(C.QtDialogGetResultBoolean(dialog.ptr)) {
            var count = int(C.QtFileDialogGetResultFileCount(dialog.ptr))
            var files = make([] string, count)
            for i := 0; i < count; i += 1 {
                var f = consumeString(
                    C.QtFileDialogGetResultFileItem(dialog.ptr, C.int(i)),
                )
                files[i] = f
            }
            k(files, true)
        } else {
            k(nil, false)
        }
        del_cb()
    })
    C.QtConsumeDialog(dialog.ptr, cgo_callback_caller, cb)
}
type FileDialogMode C.int
func FileDialogModeSave() FileDialogMode { return FileDialogMode(C.QtFileDialogModeSave) }
func FileDialogModeOpenSingle() FileDialogMode { return FileDialogMode(C.QtFileDialogModeOpenSingle) }
func FileDialogModeOpenMultiple() FileDialogMode { return FileDialogMode(C.QtFileDialogModeOpenMultiple) }

func ClipboardReadText16() ([] uint16) {
    return consumeString16(C.QtClipboardReadText())
}
func ClipboardWriteText16(text ([] uint16)) {
    var pkg, dispose = CreatePkg(); defer dispose()
    C.QtClipboardWriteText(makeString16(text, pkg))
}

// ListWidgetInterface
type Lwi struct { ptr unsafe.Pointer }
func (I Lwi) objectPointer() unsafe.Pointer {
    return I.CastToWidget().Object.ptr
}
const DefaultListWidget_CurrentChanged = "currentItemChanged(QTreeWidgetItem*,QTreeWidgetItem*)"
const DefaultListWidget_SelectionChanged = "itemSelectionChanged()"
func LwiCreateFromDefaultListWidget(columns int, select_ ItemSelectionMode, headers ([] Widget), stretch int, ctx Pkg) Lwi {
    var headers_ptr, headers_len = ptrlen(headers)
    var ptr = C.QtLwiCreateFromDefaultListWidget (
        C.size_t(columns),
        C.int(select_),
        headers_ptr, headers_len,
        C.int(stretch),
    )
    var I = Lwi{ptr}
    pushObjectDeletion(ctx, I.objectPointer())
    return I
}
func (I Lwi) CastToWidget() Widget {
    var ptr = C.QtLwiCastToWidget(I.ptr)
    return Widget{Object{ptr}}
}
func (I Lwi) Current() (string, bool) {
    var exists C.QtBool
    var key = consumeString(C.QtLwiCurrent(I.ptr, &exists))
    return key, getBool(exists)
}
func (I Lwi) All() ([] string) {
    return consumeStringList(C.QtLwiAll(I.ptr))
}
func (I Lwi) Selection() ([] string) {
    return consumeStringList(C.QtLwiSelection(I.ptr))
}
func (I Lwi) ContiguousSelection() ([] string) {
    return consumeStringList(C.QtLwiContiguousSelection(I.ptr))
}
func (I Lwi) Prepend(key string, widgets ([] Widget)) {
    var pkg, dispose = CreatePkg(); defer dispose()
    var widgets_ptr, widgets_len = ptrlen(widgets)
    C.QtLwiInsert(
        I.ptr, C.QtLwiPrepend,
        makeString("", pkg),
        makeString(key, pkg),
        widgets_ptr, widgets_len,
    )
}
func (I Lwi) Append(key string, widgets ([] Widget)) {
    var pkg, dispose = CreatePkg(); defer dispose()
    var widgets_ptr, widgets_len = ptrlen(widgets)
    C.QtLwiInsert(
        I.ptr, C.QtLwiAppend,
        makeString("", pkg),
        makeString(key, pkg),
        widgets_ptr, widgets_len,
    )
}
func (I Lwi) InsertAbove(pivot string, key string, widgets ([] Widget)) {
    var pkg, dispose = CreatePkg(); defer dispose()
    var widgets_ptr, widgets_len = ptrlen(widgets)
    C.QtLwiInsert(
        I.ptr, C.QtLwiInsertAbove,
        makeString(pivot, pkg),
        makeString(key, pkg),
        widgets_ptr, widgets_len,
    )
}
func (I Lwi) InsertBelow(pivot string, key string, widgets ([] Widget)) {
    var pkg, dispose = CreatePkg(); defer dispose()
    var widgets_ptr, widgets_len = ptrlen(widgets)
    C.QtLwiInsert(
        I.ptr, C.QtLwiInsertBelow,
        makeString(pivot, pkg),
        makeString(key, pkg),
        widgets_ptr, widgets_len,
    )
}
func (I Lwi) Update(key string) {
    var pkg, dispose = CreatePkg(); defer dispose()
    C.QtLwiUpdate(I.ptr, makeString(key, pkg))
}
func (I Lwi) MoveUp(key string) bool {
    var pkg, dispose = CreatePkg(); defer dispose()
    return getBool(C.QtLwiMove(I.ptr, C.QtLwiUp, makeString(key, pkg)))
}
func (I Lwi) MoveDown(key string) bool {
    var pkg, dispose = CreatePkg(); defer dispose()
    return getBool(C.QtLwiMove(I.ptr, C.QtLwiDown, makeString(key, pkg)))
}
func (I Lwi) Delete(key string) {
    var pkg, dispose = CreatePkg(); defer dispose()
    C.QtLwiDelete(I.ptr, makeString(key, pkg))
}
func (I Lwi) Reorder(order ([] string)) {
    var pkg, dispose = CreatePkg(); defer dispose()
    var order_ptr, order_len = makeStrings(order, pkg)
    C.QtLwiReorder(I.ptr, order_ptr, order_len)
}
type ItemSelectionMode C.int
func ItemNoSelection() ItemSelectionMode { return ItemSelectionMode(C.QtItemNoSelection) }
func ItemSingleSelection() ItemSelectionMode { return ItemSelectionMode(C.QtItemSingleSelection) }
func ItemMultiSelection() ItemSelectionMode { return ItemSelectionMode(C.QtItemMultiSelection) }
func ItemExtendedSelection() ItemSelectionMode { return ItemSelectionMode(C.QtItemExtendedSelection) }


func getBool(number C.QtBool) bool {
    return (number != 0)
}
func makeBool(p bool) C.QtBool {
    if p {
        return C.int(int(1))
    } else {
        return C.int(int(0))
    }
}

func makeVariant(v interface{}, ctx Pkg) C.QtVariant {
    var variant = (func() C.QtVariant {
        switch V := v.(type) {
        case int:
            return C.QtCreateVariantInt(C.int(V))
        case float64:
            return C.QtCreateVariantDouble(C.double(V))
        case string:
            return C.QtCreateVariantString(makeString(V, ctx))
        default:
            return C.QtCreateVariantInvalid()
        }
    })()
    ctx.push(func() {
        C.QtDeleteVariant(variant)
    })
    return variant
}

func makeString(s string, ctx Pkg) C.QtString {
    var q_str C.QtString
    if s != "" {
        var hdr = (*reflect.StringHeader)(unsafe.Pointer(&s))
        var ptr = (*C.uint8_t)(unsafe.Pointer(hdr.Data))
        var size = (C.size_t)(len(s))
        q_str = C.QtNewStringUTF8(ptr, size)
    } else {
        q_str = C.QtNewStringUTF8(nil, 0)
    }
    var del = func() { C.QtDeleteString(q_str) }
    ctx.push(del)
    return q_str
}
func makeString16(s ([] uint16), ctx Pkg) C.QtString {
    var q_str C.QtString
    if len(s) > 0 {
        var ptr = (*C.uint16_t)(unsafe.Pointer(&s[0]))
        var size = (C.size_t)(len(s))
        q_str = C.QtNewStringUTF16(ptr, size)
    } else {
        q_str = C.QtNewStringUTF16(nil, 0)
    }
    var del = func() { C.QtDeleteString(q_str) }
    ctx.push(del)
    return q_str
}
func makeStrings(all ([] string), ctx Pkg) (*C.QtString, C.size_t) {
    if len(all) == 0 {
        ctx.push(func() {})
        return nil, 0
    }
    var L = len(all)
    var buf = make([] C.QtString, len(all))
    for i := range all {
        buf[i] = makeString(all[i], ctx)
    }
    return &(buf[0]), C.size_t(L)
}
func copyString(q_str C.QtString) string {
    var size16 = uint(C.QtStringUTF16Length(q_str))
    if size16 == 0 {
        return ""
    }
    var buf = make([] rune, size16)
    var size32 = uint(C.QtStringWriteToUTF32Buffer(q_str,
        (*C.uint32_t)(unsafe.Pointer(&buf[0]))))
    buf = buf[:size32]
    return string(buf)
}
func copyString16(q_str C.QtString) ([] uint16) {
    var size = uint(C.QtStringUTF16Length(q_str))
    if size == 0 {
        return [] uint16 {}
    }
    var buf = make([] uint16, size)
    C.QtStringWriteToUTF16Buffer(q_str,
        (*C.uint16_t)(unsafe.Pointer(&buf[0])))
    return buf
}
func consumeString(q_str C.QtString) string {
    var go_str = copyString(q_str)
    C.QtDeleteString(q_str)
    return go_str
}
func consumeString16(q_str C.QtString) ([] uint16) {
    var buf = copyString16(q_str)
    C.QtDeleteString(q_str)
    return buf
}
func consumeStringList(q_list C.QtStringList) ([] string) {
    var size = int(C.QtStringListGetSize(q_list))
    var list = make([] string, size)
    for i := 0; i < size; i += 1 {
        list[i] = consumeString(C.QtStringListGetItem(q_list, C.size_t(i)))
    }
    C.QtDeleteStringList(q_list)
    return list
}

type ByteArray struct { QtByteArray C.QtByteArray }
func consumeByteArray(b ByteArray) ([] byte) {
    var buf = unsafe.Pointer(C.QtByteArrayGetBuffer(b.QtByteArray))
    var size = int(C.QtByteArrayGetSize(b.QtByteArray))
    var data = make([] byte, size, size)
    for i := 0; i < size; i += 1 {
        data[i] = *(*byte)(unsafe.Pointer(uintptr(buf) + uintptr(i)))
    }
    C.QtDeleteByteArray(b.QtByteArray)
    return data
}

func makeIcon(name string, ctx Pkg) C.QtIcon {
    var pkg, dispose = CreatePkg(); defer dispose()
    var icon C.QtIcon
    if strings.HasPrefix(name, FileIconNamePrefix) {
        var path = strings.TrimPrefix(name, FileIconNamePrefix)
        icon = C.QtCreateIconFromFile(makeString(path, pkg))
    } else {
        if name == "" {
            pkg.push(func() {})
            icon = C.QtCreateNullIcon()
        } else {
            icon = C.QtCreateIconFromStock(makeString(name, pkg))
        }
    }
    var del = func() { C.QtDeleteIcon(icon) }
    ctx.push(del)
    return icon
}
const FileIconNamePrefix = "file:"


