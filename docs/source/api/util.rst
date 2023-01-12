Utilities & Miscellaneous
+++++++++++++++++++++++++

API
===

Reflection
----------

.. code-block:: none

    type IntoReflectType[T] native
    type IntoReflectValue[T] native

.. code-block:: none

    type DummyReflectType native
    const ReflectType DummyReflectType

Lambda
------

.. code-block:: none

    type Lambda[A,B] native

.. code-block:: none

    operator *  [A,B,C] { f Lambda[B,C], g Lambda[A,B] } Lambda[A,C]
    operator -> [A,B]   { x A, f Lambda[A,B] } B
    operator call [A,B] { f Lambda[A,B], x A } B

Tuples
------

.. code-block:: none

    type Pair[A,B] record {
        First   A,
        Second  B
    }

.. code-block:: none

    type Triple[A,B,C] record {
        First   A,
        Second  B,
        Third   C
    }

.. code-block:: none

    function Pair[A,B] { first A, second B } Pair[A,B]
    function Triple[A,B,C] { first A, second B, third C } Triple[A,B,C]

Null
----

.. code-block:: none

    type Null native
    const Null Null

Bool
----

.. code-block:: none

    type Bool(~String) native
    const No  Bool
    const Yes Bool

.. code-block:: none

    function Not { p Bool } Bool
    operator and { p Bool, q Bool } Bool
    operator or  { p Bool, q Bool } Bool

.. code-block:: none

    operator ==  { p Bool, q Bool } Bool

Maybe
-----

.. code-block:: none

    type Maybe[OK] union {
        Null,
        OK
    }

.. code-block:: none

    function Nothing[T] {} Maybe[T]
    function Just[T] { value T } Maybe[T]

.. code-block:: none

    method Maybe.List List[OK]
    method Maybe.$ $[OK]

.. code-block:: none

    operator ??[T] { value? Maybe[T], fallback T } T
    operator map[A,B] { v? Maybe[A], f Lambda[A,B] } Maybe[B]
    operator filter[T] { v? Maybe[T], f Lambda[T,Bool] } Maybe[T]
    operator maybe[A,B] { v? Maybe[A], k Lambda[A,Maybe[B]] } Maybe[B]

Lens
----

.. code-block:: none

    type Lens1[Whole,Part] record {
        Value   Part,
        Assign  Lambda[Part,Whole]
    }
    type Lens2[Abstract,Concrete] record {
        Value   Maybe[Concrete],
        Assign  Lambda[Concrete,Abstract]
    }

.. code-block:: none

    method Lens1.Update Lambda[Lambda[Part,Part],Whole]
    method Lens2.Update Lambda[Lambda[Maybe[Concrete],Concrete],Abstract]
    method Lens1.Update? Lambda[Lambda[Part,Maybe[Part]],Maybe[Whole]]
    method Lens2.Update? Lambda[Lambda[Maybe[Concrete],Maybe[Concrete]],Maybe[Abstract]]

.. code-block:: none

    operator compose1[A,B,C] { ab Lens1[A,B], f Lambda[B,Lens1[B,C]] } Lens1[A,C]
    operator compose2[A,B,C] { ab Lens1[A,B], f Lambda[B,Lens2[B,C]] } Lens2[A,C]
    operator compose[A,B,C]  { ab Lens2[A,B], f Lambda[Maybe[B],Lens2[B,C]] } Lens2[A,C]

Assertion
---------

.. code-block:: none

    operator assert[T] { ok Bool, k Lambda[Null,T] } T

.. code-block:: none

    function Undefined[T] { msg String('') } T

Error
-----

.. code-block:: none

    type Error native

.. code-block:: none

    function variadic Error { msg List[~String] } Error

.. code-block:: none

    method Error.Message String
    method Error.IsCancel Bool 

.. code-block:: none

    operator wrap { err Error, msg String } Error

Result
------

.. code-block:: none

    type Result[OK] union {
        Error,
        OK
    }

.. code-block:: none

    function Success[T] { value T } Result[T]

.. code-block:: none

    method Result.Maybe Maybe[OK]

.. code-block:: none

    operator map[A,B] { r Result[A], f Lambda[A,B] } Result[B]


