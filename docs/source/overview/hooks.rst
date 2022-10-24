Hooks
+++++

.. image :: counter.png
    :scale: 62%

**Hook** is a kind of specialized Observable,
which enables a resource management mechanism
similar to RAII but deallocates resources asynchronously.

The idea of Hook and the terminology "Hook" is borrowed from ReactJS.
As a result, usage of Hook-related APIs may look very similar to ReactJS.
However, the concept of Hook is generalized,
and the underlying mechanics is completely different.

The following program implements a simple counter.
We'll explain the concept of Hook with this example.

.. literalinclude :: counter.km
    :language: none

``State``, ``PlainButton``, ``Effect`` and ``Label`` are functions
that return a Hook. The ``use`` operator uses the returned Hooks to
acquire resources (initialize widgets, create effects, etc.),
and more importantly, know how to release resources when the window is closed.

When the window is closed, the following happens successively.

* The Window is deleted,
  with its child widgets hidden and detached.
* The Label ``num`` is deleted,
  with its input Observable unsubscribed.
* The Effect involving ``inc.Clicks`` is disposed,
  which means the underlying click event handler is removed from the PlainButton ``inc``.
* The Effect involving ``dec.Clicks`` is disposed,
  which means the underlying click event handler is removed from the PlainButton ``dec``.
* The PlainButton ``inc`` is deleted.
* The PlainButton ``dec`` is deleted.
* The State ``count`` is disposed,
  which means a no-op in current implementation.

.. Note::
    If a widget is already deleted, accessing it will trigger a crash.
    But in normal circumstances, such crash is almost impossible.

The order of resource release is the reverse order of resource acquisition,
which is similar to RAII.

Also, unlike the ``useEffect`` API in ReactJS,
in the ``Effect`` calls, there is no code about how to dispose the effect,
but when the window is closed,
the click event handlers on ``inc`` and ``dec`` are removed anyways.
This is because ``inc.Clicks`` and ``dec.Clicks`` are Observables,
when they are unsubscribed,
a pre-defined behavior of removing event handlers
is triggered automatically.
This is like when using RAII to acquire a resource,
there is no need to specify a destructor explicitly,
because a pre-defined destructor gets called automatically
when program execution exited the scope.

In the example program, resources are released when the only window is closed,
which seems not that useful.
But the Hook mechanism is also used in
custom modal dialogs, conditional widgets and list widgets.
For example, when a list item is deleted,
all corresponding effects created by the item are disposed automatically.


