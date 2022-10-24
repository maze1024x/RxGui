Conditional Widgets
+++++++++++++++++++

.. image :: conditional.png
    :scale: 62%

A *conditional widget* is a dynamic wrapper widget
that shows different inner widget according to the newest value of an Observable.

We'll explain the usage of conditional widgets with the following example,
which is an arbitrary program just for demonstration purpose.

.. literalinclude :: conditional.km
    :language: none

There are two kinds of built-in conditional widgets available.

* ``Switchable`` switches between existing widgets
  according to an Observable of widgets.
* ``Reloadable``  loads new widgets
  according to an Observable of Hooks.
  It unloads previously loaded widget and loads a new widget
  each time the Observable emits a new Hook of widget.

In the example, when ``combo.SelectedItem`` emits a new value,
the Switchable ``s`` changes its inner widget
to one of the existing TextBoxes ``a``, ``b``, or ``c``,
while the Reloadable ``r`` deletes its old TextBox
and creates a new TextBox as its inner widget.

If the TextBox in the Switchable is edited, the edited text is preserved.
On the contrary, if the TextBox in the Reloadable is edited,
the edited text is discarded when another option in the ComboBox is selected
(this can be verified by selecting back to the original option),
due to the fact that the inner widget of Reloadable is actually reloaded.

The inner widgets of Switchable/Reloadable are all TextBoxes in the example.
However, it is completely possible to switch
the content of a Switchable/Reloadable between any widget kinds.


