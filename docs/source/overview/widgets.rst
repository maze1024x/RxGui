Basic Widgets
+++++++++++++

.. image :: widgets.png
    :scale: 62%

We'll explain the usage of basic widgets with the following example program.

.. literalinclude :: widgets.km
    :language: none

There are mainly two kinds of basic widgets.

A *source widget* receives user input.
It may take an initial value as an argument and produce an Observable as its output (such as what ``TextBox`` does),
or it may produce an Observable that emits discrete values (such as what ``PlainButton`` does).

A *sink widget* displays data to user.
It takes an Observable argument as its input (such as what ``Label`` does).

The example program demonstrates the most trivial usage of widgets: data echo,
i.e. data flowing from source widgets to sink widgets.
It is possible to use operators to process data from source widgets and
pass the processed data to sink widgets. The input of ``echo2`` is processed by
the ``debounce-time`` operator, the input of ``echo3`` is processed by the ``map``
operator, and the input of ``count`` is processed by the ``reduce`` operator and
the ``map`` operator successively.

.. Caution::
    The behavior of the ``reduce`` operator is different from RxJS.
    The ``reduce`` operator emits its initial value and
    emits its current value on each upstream emission.


