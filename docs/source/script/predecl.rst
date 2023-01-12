Namespace & Alias
+++++++++++++++++

A source file starts with a namespace declaration,
assigning the source file to a namespace.

.. code-block:: none

    namespace ::

.. code-block:: none

    namespace foo ::

The namespace declaration can be followed by some alias declarations.

.. code-block:: none

    using FancyName = Foo::Bar
    using Foo::Baz

.. code-block:: none

    using abc = namespace abcdefghijkl

Aliases are only effective inside the prior declared namespace,
i.e. the namespace that the current source file is assigned to.


