Keyword & Identifier
++++++++++++++++++++

Strict Keywords
===============

.. code-block:: none

    const, if, else, when, each, let, new, =>, =

Conditional Keywords
====================

.. code-block:: none

    namespace, using, run, entry, type, function, operator, method,
    native, record, interface, union, enum, variadic

Identifier Rules
================

* An identifier must NOT be
  a comment (starting with ``//``)
  or a pragma (starting with ``#!``).
* An identifier must NOT be identical to a strict keyword.
* An identifier must NOT start with a digit (0-9).
* An identifier must NOT contain any blank or line break
  (Space, full-width Space, Tab, CR, LF).
* An identifier must NOT contain any of the following symbols.

.. code-block:: none

    { } [ ] ( ) . , : ; @ | & \ ' " `

Note that ``+ - * / % ^ ~ ! ? # $ > == <`` are all valid identifiers.

Naming Conventions
==================

* type/function/constant/field/method: ``UpperCamelCase``
* operator/argument/binding: ``kebab-case``

Exceptions:

* Field names should be consistent with JSON field names
  when they are being used in serialization.
* Sometimes it is better to name things with symbols
  (e.g. operator + - * /, etc.).
* For undocumented APIs, naming conventions are intentionally ignored
  to distinguish them from documented APIs.
* Other special cases.


