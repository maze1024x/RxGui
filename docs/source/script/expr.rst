Expressions
+++++++++++

Type Casting
============

Cast ``EXPR`` to type ``TYPE``.

.. code-block:: none

    ([TYPE]) EXPR

.. code-block:: none

    EXPR.([TYPE])

Literals
========

Number
------

.. code-block:: none

    1
    -2

.. code-block:: none

    1.0
    -0.5
    2e3
    -2e-3

Bytes
-----

.. code-block:: none

    \x00\x01\x02\x03\xFF

Char
----

.. code-block:: none

    `c`
    \n
    \uFFFD

String
------

.. code-block:: none

    'raw string'
    "string"

.. code-block:: none

    'foo' .. "bar" .. \n
        .. "baz" .. \n

Lambda
======

.. code-block:: none

    { PATTERN => EXPR }

.. code-block:: none

    { PATTERN &SELF => EXPR }

The type of a lambda is inferred.
There is no way to specify the type of a lambda inside the lambda itself.
When using a standalone lambda,
an explicit type cast to ``Lambda[I,O]`` (or a SAM interface type) is required.

Block
=====

.. code-block:: none

    { EXPR }

.. code-block:: none

    {
        BINDING-1,
        BINDING-2,
        EXPR
    }

``BINDING`` can be a ``let`` binding, a ``const`` binding, or a CPS Binding.

Let Binding
-----------

.. code-block:: none

    let PATTERN = EXPR

Const Binding
-------------

.. code-block:: none

    const PATTERN = EXPR
    // non-const bindings are invisible in EXPR

CPS Binding
-----------

.. code-block:: none

    @OPERATOR PATTERN = EXPR
    // OPERATOR is called with EXPR as the 1st argument
    // and { PATTERN => remaining block content } as the 2nd argument.

Function Call
=============

Calling an ordinary function.

.. code-block:: none

    FooBar(value1, value2)

.. code-block:: none

    FooBar { arg1: value1, arg2: value2 }

Calling an operator.

.. code-block:: none

    value1 | foo-bar(value2, value3)

.. code-block:: none

    value1 | foo-bar { arg2: value2, arg3: value3 }

Calling a method.

.. code-block:: none

    value.FooBar

Using a constant.

.. code-block:: none

    FooBar

Syntactic Sugars for Calling an Operator
----------------------------------------

**Infix Call**, an additional syntax for calling a binary operator.

.. code-block:: none

    (value1 foo-bar value2)

**CPS Call**, an additional syntax for calling a binary operator
with a lambda as the 2nd argument,
which has already debuted in the section of Block (`CPS Binding`_).

.. code-block:: none

    { @foo-bar PATTERN = value1, EXPR }
    // equivalent to (value1 foo-bar { PATTERN => EXPR })

Implicit Argument
-----------------

Passing an implicit argument implicitly.

.. code-block:: none

    let ctx = value,
    FooBar(value1, value2)

.. code-block:: none

    let ctx = value,
    FooBar { arg1: value1, arg2: value2 }

Passing an implicit argument explicitly.

.. code-block:: none

    FooBar { arg1: value1, arg2: value2, ctx: value }

The sample code above calls an ordinary function,
but implicit arguments can be used in operator calls as well.

When an implicit argument is not explicitly specified in a function call,
its name is queried in the local scope first
and queried in global functions second.

For a generic function,
it is possible to declare an implicit argument with a name like ``T/op``,
in which the ``T`` part is identical to the name of a type parameter.
In this case, when ``T/op`` is not explicitly specified in a function call,
and there is no binding named ``T/op`` in the local scope,
the name ``op`` is queried in the namespace of the actual type of ``T``.
For example,
declaring an implicit argument ``T/==`` retrieves the ``==`` operator on ``T``,
declaring an implicit argument ``T/<`` retrieves the ``<`` operator on ``T``,
etc.

Usage of Bool & Maybe & Enum & Union
====================================

Bool
----

``Bool`` constants.

* ``Yes`` (aka true)
* ``No`` (aka false)

``Bool`` functions (logical operators).

* ``Not(BOOL)``
* ``(BOOL1 and BOOL2)``
* ``(BOOL1 or BOOL2)``

``if`` expressions which use ``Bool`` values as conditions.

.. code-block:: none

    if (BOOL) YES-BLOCK
    else NO-BLOCK

.. code-block:: none

    if (BOOL-1) 1-YES-BLOCK
    if (BOOL-2) 2-YES-BLOCK
    else BOTH-NO-BLOCK

.. Tip::
    There is no "else if" or "elif".
    Just write another ``if``.

.. code-block:: none

    if (BOOL-1, BOOL-2) BOTH-YES-BLOCK
    else EITHER-NO-BLOCK

.. Note::
    What ``if (BOOL-1, BOOL-2)`` does is different from ``if ((BOOL-1 and BOOL-2))``.
    The 2nd condition ``BOOL-2`` is lazy evaluated in ``if (BOOL-1, BOOL-2)``
    but eager evaluated in ``if ((BOOL-1 and BOOL-2))``.
    Operators such as ``and`` are just functions,
    whose arguments are evaluated anyways. 

Maybe
-----

``Maybe[OK]`` is a union type of ``Null`` and ``OK``.

There is a constant named ``Null``,
which is the only value of the type named ``Null``.

In most cases, ``Maybe`` values are constructed by implicit conversions.
When a value of type ``OK`` or a ``Null`` value
is passed to a context requiring a value of ``Maybe[OK]`` type,
the value is converted to ``Maybe[OK]`` type automatically.
In the cases that an implicit conversion is not applicable,
it is required to explicitly cast the value to ``Maybe[OK]`` type
or call the constructors ``Just(value)`` and ``Nothing[OK]()``.

``Maybe`` values can be used in ``if`` expressions just like ``Bool`` values.
Also, when using a ``Maybe`` value as a condition,
a ``let`` binding can be added on it.

.. code-block:: none

    if (let PATTERN = MAYBE) YES-BLOCK

.. code-block:: none

    if (let PATTERN-1 = MAYBE-1, let PATTERN-2 = MAYBE-2) BOTH-YES-BLOCK

In the sample code above,
bindings in ``PATTERN-1`` and ``PATTERN-2``
are all available in ``BOTH-YES-BLOCK``.
And especially, bindings in ``PATTERN-1``
are also available in ``MAYBE-2``.

Enum
----

Assume there is an enum type ``Foo`` with items ``A``, ``B`` and ``C``.
When the identifiers ``A``, ``B`` or ``C``
is passed to a context requiring a value of ``Foo`` type,
a ``Foo`` enum value is constructed.
When there is no such context
and a ``Foo`` enum value is desired to be constructed,
use a explicit type cast like ``([Foo]) A``.

An enum value can be passed to a context requiring an integer
or explicitly converted to a an integer.

To enumerate all items of a enum type,
use an ``each`` expression.

.. code-block:: none

    each(Foo) { A => 'A', B => 'B', C => 'C' }
    // evalutes to List('A', 'B', 'C')

All enum items must present in an ``each`` expression.

To determine the value of a enum value,
use a ``when`` expression.

.. code-block:: none

    when (foo) {
        A => EXPR-1,
        B => EXPR-2,
        C => EXPR-3
    }

By default, it is required to list all possible branches in a ``when`` expression.
A default branch can be provided to ignore this restriction.

.. code-block:: none

    _ => EXPR  // a default branch

In addition, enum values can be converted to a ``Lens2`` type,
which can be used in ``if`` expressions,
as shown below.

.. code-block:: none

    if (foo.(A)) A-BLOCK
    else OTHERWISE-BLOCK

.. code-block:: none

    if (foo.(A), bar.(B)) A-B-BLOCK
    if (foo.(B), bar.(A)) B-A-BLOCK
    else OTHERWISE-BLOCK

.. Note::
    Although an enum value is basically an integer,
    an enum type does *NOT* have comparison operators
    like ``<`` ``==`` ``<>`` pre-defined.
    Comparison operators on enums are defined on demand.

    For example, the ``==`` operator on ``Foo`` can be defined as follows.
    
    .. code-block:: none

        operator == { a Foo, b Foo } Bool { (([Int]) a == b) }

Union
-----

Assume there is a union type ``Foo[Other]``
with item types ``Type1``, ``Type2[String]`` and ``Other``.
When a value is passed to a context requiring a value of ``Foo[X]`` type,
the value is converted to ``Foo[X]`` type automatically
if its type is one of ``Type1``, ``Type2[String]`` and ``X``.

To enumerate all items types of a *concrete* union type,
use an ``each`` expression.

.. code-block:: none

    each(Foo[X]) { Type1 t1 => t1(a), Type2 t2 => t2(b), Other o => o(x) }
    // evaluates to List[Foo[X]](([Type1]) a, ([Type2[String]]) b, ([X]) x)

All union items must present in an ``each`` expression.

To determine the inner type and extract the inner value of a union value,
use a ``when`` expression.

.. code-block:: none

    when (foo) {
        Type1 PATTERN-1 => EXPR-1,
        Type2 PATTERN-2 => EXPR-2,
        Other PATTERN-3 => EXPR-3
    }

Similar to a ``when`` expression on an enum value,
a ``when`` expression on a union value
is also required to be exhaustive by default.
A default branch can be provided to ignore this restriction.

.. code-block:: none

    _ => EXPR  // a default branch

In addition, similar to enum values,
union values also can be converted to a ``Lens2`` type,
which can be used in ``if`` expressions,
as shown below.

.. code-block:: none

    if (let PATTERN = foo.(Type1)) 1-BLOCK
    else OTHERWISE-BLOCK

.. code-block:: none

    if (let v1 = foo.(Type1), let v2 = bar.(Type2)) 1-2-BLOCK
    if (let v2 = foo.(Type2), let v1 = bar.(Type1)) 2-1-BLOCK
    if (let foo = foo.(Other), let bar = bar.(Other)) OTHER-OTHER-BLOCK
    else OTHERWISE-BLOCK

Usage of Record
===============

Constructing a record value.

.. code-block:: none

    new FooBar(value1, value2)

.. code-block:: none

    new FooBar { Field1: value1, Field2: value2 }

Pattern matching on a record value.

.. code-block:: none

    let (a,b) = foobar

Getting the value of a field of a record value.

.. code-block:: none

    foobar.Field1

Constructing an updated record value with an updated field.

.. code-block:: none

    foobar.(Field1).Assign(new-value)

.. code-block:: none

    foobar.(Field1).Update({ old-value => { let new-value = f(old-value), new-value }})

Updating multiple fields.

.. code-block:: none

    foobar.(Field1).Assign(new-value-1)
          .(Field2).Assign(new-value-2)

Updating a field in an inner record.

.. code-block:: none

    foobar.(Field1).(Field1).Assign(new-value)

``foobar.(Field1)`` converts ``foobar`` to a ``Lens1`` type,
which is composable.

.. Tip::
    With a type named ``FooBar`` already defined,
    it is still okay to define a function named ``FooBar``,
    as a conventional constructor.

Record Observable Projection
----------------------------

Making a field projection from a record Observable.

.. code-block::

    $(Pair(1,'a'), Pair(2,'b'), Pair(3,'c')).First
    // emits 1, 2, 3

.. code-block::

    $(Pair(1,'a'), Pair(2,'b'), Pair(3,'c')).Second
    // emits "a", "b", "c"

The ``map`` operator can be used to achieve similar result,
but when using the ``map`` operator,
even if the value of the field isn't updated,
a value that is identical to the previous one is emitted anyways.

.. code-block::

    { let p = Pair(1,'foo'), $(p, p.(Second).Assign('bar')) | map({ pair => pair.First }) }
    // emits 1, 1

.. code-block::

    { let p = Pair(1,'foo'), $(p, p.(Second).Assign('bar')).First }
    // emits 1

Decorated Construction of Record
--------------------------------

Combining field Observables into a record Observable.

.. code-block::

    new:$ FooBar(o1, o2)

.. code-block::

    new:$ FooBar { Field1: o1, Field2: o2 }

Combining field Maybe-Observables into a record Maybe-Observable.

.. code-block::

    new:$? FooBar(o1, o2)

.. code-block::

    new:$? FooBar { Field1: o1, Field2: o2 }

Combining Hooks of fields into a Hook of record.

.. code-block::

    new:Hook FooBar { Field1: hook1, Field2: hook2 }

Usage of Interface
==================

Similar to other languages,
the following conversions can be done implicitly or explicitly.

* Converting values from concrete types to interface types.
* Converting values from interface types to more abstract interface types.

Dynamic Conversion
------------------

Trying to extract a value of a specific concrete type from an interface value,
i.e. converting the interface value to a ``Lens2`` type.

.. code-block:: none

    if (let concrete = abstract.(Concrete))

Trying to extract a value of a more concrete interface type from an interface value.

.. code-block:: none

    if (let more-concrete = ([Maybe[MoreConcrete]]) abstract)

SAM Interface
-------------

When assigning a value to a SAM interface type,
if the type of the value is assignable to the method type,
the value is converted to the interface type automatically.

When the method in a SAM interface has a ``Lambda`` type,
values of the interface type become directly callable.


