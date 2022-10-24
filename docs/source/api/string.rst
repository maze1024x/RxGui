String Manipulation
+++++++++++++++++++

API
===

~String
-------

.. code-block:: none

    type ~String interface { String String }

.. code-block:: none

    method Bool.String     String
    method Ordering.String String
    method Int.String   String
    method Float.String String
    method Char.String  String

.. code-block:: none

    function ParseInt { s String } Maybe[Int]
    function ParseFloat { s String } Maybe[Float]

Char
----

.. code-block:: none

    type Char(~String) native
    function Char { value Int } Char
    method Char.Int Int
    method Char.Utf8Size Int

.. code-block:: none

    operator == { a Char, b Char } Bool
    operator <  { a Char, b Char } Bool
    operator <> { a Char, b Char } Ordering

String
------

.. code-block:: none

    type String native

.. code-block:: none

    function variadic String { fragments List[~String] } String
    function StringFromChars { chars List[Char] } String

.. code-block:: none

    operator == { a String, b String } Bool    
    operator <  { a String, b String } Bool    
    operator <> { a String, b String } Ordering

.. code-block:: none

    function Quote { s String } String
    function Unquote { s String } Maybe[String]

.. code-block:: none

    method String.Empty Bool
    method String.Chars List[Char]
    method String.FirstChar Maybe[Char]
    method String.NumberOfChars Int
    method String.Utf8Size Int

.. code-block:: none

    operator shift { s String } Maybe[Pair[Char,String]]
    operator reverse { s String } String
    operator join  { l List[String], sep String } String
    operator split { s String, sep String } List[String]
    operator cut   { s String, sep String } Maybe[Pair[String,String]]
    operator has-prefix { s String, prefix String } Bool
    operator has-suffix { s String, suffix String } Bool
    operator trim-prefix { s String, prefix String } String
    operator trim-suffix { s String, suffix String } String
    operator trim { s String, chars List[Char] } String
    operator trim-left { s String, chars List[Char] } String
    operator trim-right { s String, chars List[Char] } String

RegExp
------

.. code-block:: none

    type RegExp(~String) native

.. code-block:: none

    method RegExp.String String

.. code-block:: none

    operator advance { s String, re RegExp } Maybe[Pair[String,String]]
    operator satisfy { s String, re RegExp } Bool
    operator replace { s String, re RegExp, f Lambda[String,String] } String

.. Note::
    When using a string literal in a context requiring a value of ``RegExp`` type,
    the string literal constructs a ``RegExp`` value if it is a valid regexp.

JSON Parser Example
===================

.. literalinclude :: json.km
    :language: none

.. Tip::
    The example program above uses a custom namespace ``json``.
    Running it requires specifying a command line argument ``--entry=json``.


