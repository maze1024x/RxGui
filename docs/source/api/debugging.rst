Debugging
+++++++++

API
===

.. code-block:: none

    function DebugInspect[T] { hint String, v IntoReflectValue[T] } T

.. code-block:: none

    function DebugExpose[T] { name String, v IntoReflectValue[T] } $[Null]

.. code-block:: none

    function DebugTrace[T] { hint String, o IntoReflectValue[$[T]] } $[T]

.. code-block:: none

    function DebugWatch[T] { hint String, o IntoReflectValue[$[T]] } $[Null]


