Arithmetic
++++++++++

API
===

Ordering
--------

.. code-block:: none

    type Ordering(~String) enum { L<R, L=R, L>R }
    operator == { o1 Ordering, o2 Ordering } Bool

.. code-block:: none

    type == [T] interface { Operator Lambda[Pair[T,T],Bool] }
    type <  [T] interface { Operator Lambda[Pair[T,T],Bool] }
    type <> [T] interface { Operator Lambda[Pair[T,T],Ordering] }

.. code-block:: none

    operator != [T] { a T, b T } { T/== ==[T] } Bool
    operator >  [T] { a T, b T } { T/< <[T] }   Bool
    operator <= [T] { a T, b T } { T/< <[T] }   Bool
    operator >= [T] { a T, b T } { T/< <[T] }   Bool
    function Min[T] { a T, b T } { T/< <[T] }   T
    function Max[T] { a T, b T } { T/< <[T] }   T

Int
---

.. code-block:: none

    type Int(~String) native

.. code-block:: none

    operator + { a Int, b Int } Int
    operator - { a Int, b Int } Int
    operator * { a Int, b Int } Int
    operator / { a Int, b Int } Int
    operator % { a Int, b Int } Int
    operator ^ { a Int, b Int } Int

.. code-block:: none

    operator == { a Int, b Int } Bool
    operator <  { a Int, b Int } Bool
    operator <> { a Int, b Int } Ordering

Float
-----

.. code-block:: none

    type Float(~String) native

.. code-block:: none

    operator + { x Float, y Float } Float
    operator - { x Float, y Float } Float
    operator * { x Float, y Float } Float
    operator / { x Float, y Float } Float
    operator % { x Float, y Float } Float
    operator ^ { x Float, y Float } Float

.. code-block:: none

    operator == { x Float, y Float } Bool
    operator <  { x Float, y Float } Bool

.. code-block:: none

    method Float.Int Int
    method Int.Float Float

.. code-block:: none

    const NaN Float
    const +Inf Float
    const -Inf Float
    method Float.Normal   Bool
    method Float.NaN      Bool
    method Float.Infinite Bool

.. code-block:: none

    const E  Float
    const PI Float
    function Floor { x Float } Float
    function Ceil  { x Float } Float
    function Round { x Float } Float
    function Sqrt  { x Float } Float
    function Cbrt  { x Float } Float
    function Exp   { x Float } Float
    function Log   { x Float } Float
    function Sin   { x Float } Float
    function Cos   { x Float } Float
    function Tan   { x Float } Float
    function Asin  { x Float } Float
    function Acos  { x Float } Float
    function Atan  { x Float } Float
    function Atan2 { y Float, x Float } Float


