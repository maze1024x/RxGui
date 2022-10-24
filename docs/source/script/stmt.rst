Statements
++++++++++

Entry Point
===========

.. code-block:: none

    entry BLOCK

``BLOCK`` should evaluate to an Observable
emitting Null values (or not emitting any value).
A program starts by subscribing to the Observable,
and exits when the Observable completes.
If the Observable throws an error, the program will crash.

Each namespace can have 0 or 1 entry point.
By default, the interpreter uses the entry point in the default namespace ``::``.
To override this behavior,
specify a command line option ``--entry=NAMESPACE``
to tell the interpreter to use the entry point in the namespace ``NAMESPACE::``.

Type Declaration
================

.. code-block:: none
    
    type Foo DEFINITION

.. code-block:: none

    type Foo[T,U] (Interface1,Interface2) DEFINITION

``DEFINITION`` can be a record, interface, union, or enum.

Record
------

.. code-block:: none

    record { Field1 Type1, Field2 Type2, Field3 Type3(DefaultValue) }

Interface
---------

.. code-block:: none

    interface { SingleAbstractMethod Type }

.. code-block:: none

    interface { AbstractMethod1 Type1, AbstractMethod2 Type2, AbstractMethod3 Type3 }

The abstraction capability of interfaces is restricted intentionally, as follows.

* For a generic interface,
  its implementations must also be generic and have identical parameters.
* For a non-generic interface,
  its implementations must also be non-generic.

i.e. types implement interfaces instead of instantiations of interfaces.

Union
-----

.. code-block:: none

    union { Type1, Type2, Type3 }

Unions are implicitly tagged, i.e. they are essentially tagged unions.

Enum
----

.. code-block:: none

    enum { A, B, C }

Enums cannot be generic.

Function Declaration
====================

Ordinary Function
-----------------

.. code-block:: none
    
    function FooBar { arg1 Type1, arg2 Type2 } ReturnType BLOCK

.. code-block:: none

    function variadic FooBar
        [T,U]
        { arg1 Type1, arg2 Type2(DefaultValue), arg3 List[Type3] }
        { ctx  Type4, T/op Type5[T] }
        ReturnType
        BLOCK

``ctx`` and ``T/op`` in the sample code above are implicit arguments.

If a function is variadic, its last explicit argument must have ``List`` type.

Operator
--------

The grammar of operators is almost identical to ordinary functions.
Just change the leading keyword ``function`` to ``operator``,
an ordinary function becomes an operator.

Note that the naming convention of operators is different from
the naming convention of ordinary functions.
``foo-bar`` should be used as an operator name instead of ``FooBar``.

Method
------

.. code-block:: none

    method Foo.Bar ReturnType BLOCK

There is an implicit binding ``this`` available in ``BLOCK``.

The behavior of methods is different from classic languages.
A method is considered to be a computed field (aka getter),
which should only depend on ``this``,
thus it doesn't take arguments.
When classic method behavior is desired, return a lambda.

Constant
--------

.. code-block:: none

    const FooBar ReturnType BLOCK

A constant is technically just a special function
that takes 0 arguments and only evaluates once.


