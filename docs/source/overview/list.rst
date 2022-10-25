List Widgets
++++++++++++

.. image :: list.png
    :scale: 62%

We'll explain the usage of *list widgets*
with the help of the following example program,
which implements an ad-hoc todo list with an extra list echoing the current list.

.. literalinclude :: list.rxsc
    :language: none

There are two kinds of list widgets.

* ``ListEditView`` initializes with an initial list
  and receives list edit operations (insertion, modification, deletion, move, etc.)
  while producing an output Observable emitting the current list.
* ``ListView`` subscribes to an Observable emitting lists
  and updates its UI according to the newest list.
  Individual items are identified with a string key.

In the example, the output Observable of the ListEditView
is passed to the ListView, which causes the ListView to echo the current list
managed by the ListEditView.

Classic Editable List
=====================

.. image :: extension.png
    :scale: 62%

In a typical classic GUI app,
an editable list often doesn't have input widgets (such as TextBox) inside it.
Instead, there is an external editor widget beside the list,
and the list itself is used to pick a current editing item.
It may also have up/down buttons beside it to adjust the order of items.

Such a classic editable list can be easily implemented with ListEditView,
as shown by the following example program.

.. literalinclude :: extension.rxsc
    :language: none

Multi-column Data List
======================

.. image :: headers.png
    :scale: 62%

A multi-column data list with custom headers providing sort/filter features
can be easily implemented with ListView,
as shown by the following example program.

.. literalinclude :: headers.rxsc
    :language: none


