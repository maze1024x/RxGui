User Interface
++++++++++++++

API
===

Measurement
-----------

.. code-block:: none

    const FontSize Int

``FontSize`` is the default font size used in a GUI program,
which is equal to the absolute value of 16 relative pixels.

Lengths in all GUI APIs are measured by relative pixels.
The definition of relative pixel is as follows.

    (absolute pixels) = (relative pixels) * (length of shortest edge of screen) / 768

Hook
----

.. code-block:: none

    type Hook[T] record { [undocumented] }

.. code-block:: none

    function Hook[T] { value T } Hook[T]

.. code-block:: none

    operator use[A,B] { h Hook[A], f Lambda[A,Hook[B]] } Hook[B]
    function Hooks[T] { l List[Hook[T]] } Hook[List[T]]

.. code-block:: none

    operator run[T,X] { h Hook[T], f Lambda[T,$[X]] } $[X]
    operator run[T,X] { o $[Hook[T]], f Lambda[T,$[X]] } $[X]

.. code-block:: none

    function Effect { effect $[Null] } Hook[Null]
    function State[T] { initial T } Hook[State[T]]
    function Memo[T] { o $[T] } Hook[$[T]]
    function variadic Subject[T] { replay Int(0), items List[T] } Hook[Subject[T]]
    function Multicasting[T] { o $[T] } Hook[$[T]]
    function Style { w Widget, o $[String] } Hook[Null]

State
-----

.. code-block:: none

    type State[T] record { [undocumented] }  // a specialized Subject

.. code-block:: none

    function CreateState[T] { initial T } $[State[T]]

.. code-block:: none

    method State.Value $[T]
    method State.$ $[T]  // equivalent to "Value"

.. code-block:: none

    operator bind-update[T] { state State[T], updates $[Lambda[T,T]] } $[Null]
    operator bind-override[T] { state State[T], values $[T] } $[Null]
    operator update[T] { state State[T], update Lambda[T,T] } $[Null]
    operator override[T] { state State[T], value T } $[Null]

.. code-block:: none

    function MakeMemo[T] { o $[T] } $[$[T]]

Widget
------

.. code-block:: none

    type Widget native

.. code-block:: none

    type Widgets union { Widget, List[Widget] }

Window
------

.. code-block:: none

    type Window record { Widget Widget, [undocumented] }

.. code-block:: none

    function Window {
        title     $[String],
        layout    Layout,
        menu-bar  MenuBar (MenuBar()),
        tool-bar  ToolBar (ToolBar()),
        exit      Lambda[$[Null],$[Null]] ({ closes => closes }),
        margin-x  Int (6),
        margin-y  Int (6),
        width     Int (-1),
        height    Int (-1),
        icon      Icon (Icon('window'))
    }
    Hook[Window]

.. code-block:: none

    function ShowWindow { h Hook[Window] } $[Null]

Custom Dialog
-------------

.. code-block:: none

    type Dialog[T] record { Widget Widget, [undocumented] }

.. code-block:: none

    function Dialog[T] {
        title     $[String],
        layout    Layout,
        exit      Lambda[$[Null],$[T]],
        margin-x  Int (6),
        margin-y  Int (6),
        width     Int (-1),
        height    Int (-1),
        icon      Icon (Icon('window'))
    }
    Hook[Dialog[T]]

.. code-block:: none

    function ShowDialog[T] { h Hook[Dialog[T]] } $[T]

Standard Dialogs
----------------

.. Caution::
    The returned Observables of standard dialog functions
    complete without emitting any value
    when dialogs are closed without pressing a button (e.g. with the Esc key)
    or with pressing buttons like cancel/abort.

.. code-block:: none

    function ShowInfo     { content String, title String('Info') } $[Null]
    function ShowWarning  { content String, title String('Warning') } $[Null]
    function ShowCritical { content String, title String('Error') } $[Null]

.. code-block:: none

    type Retry/Ignore enum { Retry, Ignore }
    type Save/Discard enum { Save, Discard }
    function ShowAbortRetryIgnore  { content String, title String('Error') } $[Retry/Ignore]
    function ShowSaveDiscardCancel { content String, title String('Save Changes') } $[Save/Discard]
    function ShowYesNo { content String, title String('Question') } $[Bool]

.. code-block:: none

    function GetChoice[T] { prompt String, items List[ComboBoxItem[T]], title String('Select') } $[T]
    function GetLine { prompt String, initial String(''), title String('Input') } $[String]
    function GetText { prompt String, initial String(''), title String('Input') } $[String]
    function GetInt { prompt String, initial Int(0), title String('Input') } $[Int]
    function GetFloat { prompt String, initial Float(0), title String('Input') } $[Float]

.. code-block:: none

    function GetFileListToOpen[T] { filter String } $[List[File]]
    function GetFileToOpen[T] { filter String } $[File]
    function GetFileToSave[T] { filter String } $[File]

Action
------

.. code-block:: none

    type Action native

.. code-block:: none

    function Action {
        icon      Icon,
        text      String,
        shortcut  String (''),
        repeat    Bool (No),
        enable    $[Bool] ($())
    }
    Hook[Action]

.. code-block:: none

    method Action.Triggers $[Null]

.. code-block:: none

    type ActionCheckBox record { Checked $[Bool] }
    function ActionCheckBox { action Action, checked Bool } Hook[ActionCheckBox]

.. code-block:: none

    type ActionComboBox[T] record { SelectedItem $[T] }
    function ActionComboBox[T] { items List[ActionComboBoxItem[T]] } Hook[ActionComboBox[T]]
    type ActionComboBoxItem[T] record { Action Action, Value T, Selected Bool }
    function ActionComboBoxItem[T] { action Action, value T, selected Bool } ActionComboBoxItem[T]

Menu
----

.. code-block:: none

    type Menu record { Icon Icon, Name String, Items List[MenuItem] }

.. code-block:: none

    function variadic Menu { icon Icon, name String, items List[MenuItem] } Menu

.. code-block:: none

    type MenuItem union { Menu, Action, Separator }
    type Separator record {}
    function Separator {} Separator

Context Menu
------------

.. code-block:: none

    function ContextMenu { w Widget, m Menu } Hook[Null]

Menu Bar
--------

.. code-block:: none

    type MenuBar record { Menus List[Menu] }

.. code-block:: none

    function variadic MenuBar { menus List[Menu] } MenuBar

Tool Bar
--------

.. code-block:: none

    type ToolBar record { Mode ToolBarMode, Items List[ToolBarItem] }

.. code-block:: none

    function variadic ToolBar { mode ToolBarMode(IconOnly), items List[ToolBarItem] } ToolBar

.. code-block:: none

    type ToolBarMode enum { IconOnly, TextOnly, TextBesideIcon, TextUnderIcon }

.. code-block:: none

    type ToolBarItem union { Menu, Action, Separator, Widget, Spacer }

Icon
----

.. code-block:: none

    type Icon record { Name String }

.. code-block:: none

    function Icon { name String('') } Icon

Icon name formats are as follows.

* Tango icons: ``name``
* External file icons: ``file:relative/path/to/icon``

Layout
------

.. code-block:: none

    type Layout union { Row, Column, Grid }

.. code-block:: none

    type LayoutItem union { Layout, Widget, Spacer, String }
    type Spacer record { Width Int, Height Int, Expand Bool }
    function Spacer { width Int(0), height Int(0), expand Bool(Yes) } LayoutItem

.. code-block:: none

    type Row record { Items List[LayoutItem], Spacing Int(4) }
    function variadic Row { items List[LayoutItem] } Layout

.. code-block:: none

    type Column record { Items List[LayoutItem], Spacing Int(4) }
    function variadic Column { items List[LayoutItem] } Layout

.. code-block:: none

    type Grid record { Items List[Span], RowSpacing Int(4), ColumnSpacing Int(4) }
    type Span record { Item LayoutItem, Row Int, Column Int, RowSpan Int, ColumnSpan Int, Align Align }
    function variadic Grid { spans List[Span] } Layout
    function Span { item LayoutItem, row Int, column Int, align Align(Default), row-span Int(1), column-span Int(1) } Span

.. code-block:: none

    type Align enum { Default, Center, Left, Right, Top, Bottom, LeftTop, LeftBottom, RightTop, RightBottom }
    function Aligned { align Align, item LayoutItem } Layout

.. code-block:: none

    function variadic Form { pairs List[Pair[LayoutItem,LayoutItem]] } Layout

Wrapper
-------

.. code-block:: none

    type Wrapper record { Widget Widget }

.. code-block:: none

    function Wrapper { layout Layout, policy-x SizePolicy(Free), policy-y SizePolicy(Free) } Hook[Wrapper]

.. code-block:: none

    function WrapperWithMargins {
        layout    Layout,
        margin-x  Int (6),
        margin-y  Int (4),
        policy-x  SizePolicy (Free),
        policy-y  SizePolicy (Free)
    }
    Hook[Wrapper]

.. code-block:: none

    type SizePolicy enum {
        Rigid,
        Controlled,
        Incompressible,
        IncompressibleExpanding,
        Free,
        FreeExpanding,
        Bounded
    }

ScrollArea
----------

.. code-block:: none

    type ScrollArea record { Widget Widget }

.. code-block:: none

    function ScrollArea { scroll Scroll, layout Layout } Hook[ScrollArea]

.. code-block:: none

    function ScrollAreaWithMargins {
        scroll    Scroll,
        layout    Layout,
        margin-x  Int (6),
        margin-y  Int(4)
    }
    Hook[ScrollArea]

.. code-block:: none

    type Scroll enum { BothDirection, VerticalOnly, HorizontalOnly }

GroupBox
--------

.. code-block:: none

    type GroupBox record { Widget Widget }

.. code-block:: none

    function GroupBox { title String, layout Layout, margin-x Int(0), margin-y Int(0) }

Splitter
--------

.. code-block:: none

    type Splitter record { Widget Widget }

.. code-block:: none

    function variadic Splitter { content List[Widget] } Hook[Splitter]

Switchable
----------

.. code-block:: none

    type Switchable record { Widget Widget }

.. code-block:: none

    function Switchable { widgets $[Widget] } Hook[Switchable]

Reloadable
----------

.. code-block:: none

    type Reloadable record { Widget Widget }

.. code-block:: none

    function Reloadable { hooks $[Hook[Widget]] } Hook[Reloadable]

.. code-block:: none

    function LazyReloadable { hooks $[Hook[Widget]] } Hook[Reloadable]

Label
-----

.. code-block:: none

    type Label record { Widget Widget }

.. code-block:: none

    function Label { text $[String], align Align(Left), selectable Bool(No) } Hook[Label]

ElidedLabel
-----------

.. code-block:: none

    type ElidedLabel record { Widget Widget }

.. code-block:: none

    function ElidedLabel { text $[String] } Hook[ElidedLabel]

IconLabel
---------

.. code-block:: none

    type IconLabel record { Widget Widget }

.. code-block:: none

    function IconLabel { icon Icon, size IconSize(Auto) } Hook[IconLabel]

.. code-block:: none

    type IconSize enum { Auto, Small, Medium, Large }

TextView
--------

.. code-block:: none

    type TextView record { Widget Widget }

.. code-block:: none

    function TextView { text $[String], format TextFormat(Plain) } Hook[TextView]

.. code-block:: none

    type TextFormat enum { Plain, Html, Markdown }

CheckBox
--------

.. code-block:: none

    type CheckBox record { Widget Widget }

.. code-block:: none

    method CheckBox.Checked $[Bool]

.. code-block:: none

    function CheckBox { text String, checked Bool } Hook[CheckBox]


ComboBox
--------

.. code-block:: none

    type ComboBox[T] record { Widget Widget, SelectedItem $[T] }

.. code-block:: none

    function ComboBox[T] { items List[ComboBoxItem[T]] } Hook[ComboBox[T]]

.. code-block:: none

    type ComboBoxItem[T] record { Icon Icon, Name String, Value T, Selected Bool }
    function ComboBoxItem[T] { icon Icon, name String, value T, selected Bool } ComboBoxItem[T]

Button
------

.. code-block:: none

    type Button record { Widget Widget }

.. code-block:: none

    method Button.Clicks $[Null]

.. code-block:: none

    function Button { icon Icon, text String, tooltip String(''), enable $[Bool]($()) } Hook[Button]
    function PlainButton { text String, enable $[Bool]($()) } Hook[Button]
    function IconButton { icon Icon, tooltip String, enable $[Bool]($()) } Hook[Button]

TextBox
-------

.. code-block:: none

    type TextBox record { Widget Widget }

.. code-block:: none

    method TextBox.Text $[String]
    method TextBox.Enters $[Null]
    method TextBox.TextOn Lambda[$[Null],$[String]]
    method TextBox.TextOnEnters $[String]

.. code-block:: none

    function TextBox { text String('') } Hook[TextBox]

.. code-block:: none

    operator bind-override { edit TextBox, text $[String] } $[Null]

TextArea
--------

.. code-block:: none

    type TextArea record { Widget Widget }

.. code-block:: none

    method TextArea.Text $[String]

.. code-block:: none

    function TextArea { text String('') } Hook[TextArea]

Slider
------

.. code-block:: none

    type Slider record { Widget Widget }

.. code-block:: none

    method Slider.Value $[Int]

.. code-block:: none

    function Slider { value Int, min Int, max Int } Hook[Slider]

ProgressBar
-----------

.. code-block:: none

    type ProgressBar record { Widget Widget }

.. code-block:: none

    function ProgressBar { value $[Int], max Int, format String('') } Hook[ProgressBar]

ListView
--------

.. code-block:: none

    type ListView[T] record {
        Widget     Widget,
        Extension  $[Maybe[Widget]],
        Current    $[Maybe[String]],
        Selection  $[List[String]]
    }

.. code-block:: none

    function ListView[T] {
        data     $[List[T]],
        key      Lambda[T, String],
        content  Lambda[Pair[$[T],ItemInfo], Hook[ItemView]],
        headers  List[HeaderView] (List()),
        stretch  Int (-1),
        select   ItemSelect (Single)
    }
    Hook[ListView[T]]

.. code-block:: none

    type ItemView record { [undocumented] }

.. code-block:: none

    function ItemView { widgets Widgets, extension Maybe[Widget] (Null) } Hook[ItemView]

ListEditView
------------

.. code-block:: none

    type ListEditView[T] record {
        Widget     Widget,
        Output     $[List[T]],
        Extension  $[Maybe[Widget]],
        [undocumented]
    }

.. code-block:: none

    function ListEditView[T] {
        initial  List[T],
        content  Lambda[Pair[T,ItemInfo], Hook[ItemEditView[T]]],
        headers  List[HeaderView] (List()),
        stretch  Int (-1),
        select   ItemSelect (Single)
    }
    Hook[ListEditView[T]]

.. code-block:: none

    type ItemEditView[T] record { [undocumented] }

.. code-block:: none

    function ItemEditView[T] {
        widgets       Widgets,
        extension     Maybe[Widget] (Null),
        update        $[T] ($()),
        delete        $[Null] ($()),
        move-up       $[Null] ($()),
        move-down     $[Null] ($()),
        move-top      $[Null] ($()),
        move-bottom   $[Null] ($()),
        insert-above  $[T] ($()),
        insert-below  $[T] ($())
    }
    Hook[ItemEditView[T]]

.. code-block:: none

    operator bind-update[T] {
        list         ListEditView[T],
        prepend      $[T] ($()),
        append       $[T] ($()),
        delete       $[Null] ($()),
        move-up      $[Null] ($()),
        move-down    $[Null] ($()),
        move-top     $[Null] ($()),
        move-bottom  $[Null] ($()),
        reorder      $[Lambda[List[T],List[T]]] ($())
    }
    $[Null]

List Commons
------------

.. code-block:: none

    type HeaderView union { String, Widget }

.. code-block:: none

    type ItemSelect enum { N/A, Single, Multiple, MaybeMultiple }

.. code-block:: none

    type ItemInfo record { Key String, Pos $[ItemPos] }
    type ItemPos record { Index Int, Total Int }
    method ItemInfo.IsFirst $[Bool]
    method ItemInfo.IsLast  $[Bool]

Editor
------

.. code-block:: none

    type Editor[T] record {
        Widget  Widget,
        [undocumented]
    }

.. code-block:: none

    function Editor[T] {
        initial  EditorDocument[T],
        content  Lambda[T,Hook[EditorView[T]]],
        open     EditorOpenBehavior[T] ({ => $() }),
        save     EditorSaveBehavior[T] ({ doc => $(doc.File?) })
    }
    Hook[Editor[T]]

.. code-block:: none

    type EditorDocument[T] record {
        File?  Maybe[File],
        Data   T
    }
    function EditorDocument[T] { file? Maybe[File], data T } EditorDocument[T]

.. code-block:: none

    type EditorView[T] record {
        Widget     Widget,
        NewValues  $[T]
    }
    function EditorView[T] { w Widget, new-values $[T] } Hook[EditorView[T]]

.. code-block:: none

    type EditorOpenBehavior[T] interface {
        Open Lambda[Bool,$[EditorDocument[T]]]
    }
    type EditorSaveBehavior[T] interface {
        Save Lambda[EditorDocument[T],$[Maybe[File]]]
    }

.. code-block:: none

    method Editor.Document $[EditorDocument[T]]
    method Editor.File $[Maybe[File]]
    method Editor.Output $[T]
    method Editor.LastSave $[T]
    method Editor.Modified $[Bool]

.. code-block:: none

    operator bind-override[T] { e Editor[T], values $[T] } $[Null]
    operator bind-reset[T] { e Editor[T], triggers $[Null] } $[Null]
    operator bind-open[T] { e Editor[T], triggers $[Null] } $[Null]
    operator bind-save[T] { e Editor[T], triggers $[Null] } $[Null]
    operator bind-save-as[T] { e Editor[T], triggers $[Null] } $[Null]
    operator ask-for-save { e Editor[T], message String, title String('Save Changes') } $[Null]

Todo List Example
=================

.. image :: todo.png
    :scale: 62%

.. literalinclude :: todo.rxsc
    :language: none

.. Note::
    The example program above mimics the todo list example
    commonly used by Web frameworks.
    However, when implementing something like a todo list in a classic widget app,
    instead of mimicking Web,
    it is better to use a classic editable list,
    which has debuted in the overview part of the documentation.

SDI Text Editor Example
=======================

.. image :: sdi.png
    :scale: 62%

.. literalinclude :: sdi.rxsc
    :language: none

.. Tip::
    The ``auto-map`` operator behaves like the ``merge-map`` operator,
    but unsubscribes upstream when a running child Observable completes
    while there is no other running child Observables.
    This behavior makes the example program exit
    when the last window is closed.


