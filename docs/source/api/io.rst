I/O
+++

API
===

Bytes
-----

.. code-block:: none

    type Bytes native

Currently there is no API for the ``Bytes`` type.

Serialization
-------------

.. code-block:: none

    function Serialize[T] { v IntoReflectValue[T] } $[String]
    function Deserialize[T] { s String, t IntoReflectType[T] } $[T]

Time
----

.. code-block:: none

    type Time(~String) native

.. code-block:: none

    const Now $[Time]

.. code-block:: none

    method Time.String String

.. code-block:: none

    operator -ms { t Time, u Time } Int

.. code-block:: none

    function TimeOf[T] { o $[T] } $[Time]
    operator with-time[T] { o $[T] } $[Pair[T,Time]]

Process
-------

.. code-block:: none

    const Arguments List[String]  // without interpreter arguments and program name
    const Environment List[String]

File
----

.. code-block:: none

    type File(~String) native

.. code-block:: none

    method File.String String

.. code-block:: none

    operator == { f File, g File } Bool

.. code-block:: none

    function ReadTextFile { f File } $[String]
    function WriteTextFile { f File, text String } $[Null]

Config
------

.. code-block:: none

    function ReadConfig[T] { dir String, name String, default IntoReflectValue[T] } $[T]
    function WriteConfig[T] { dir String, name String, value IntoReflectValue[T] } $[Null]

Path of config file is as follows.

* Linux: $HOME/.config/``dir``/``name``
* Windows: %AppData%/``dir``/``name``

Request
-------

.. code-block:: none

    function Get[Resp] { endpoint String, t IntoReflectType[Resp], token String('') } $[Resp]
    function Post[Req,Resp] { data IntoReflectValue[Req], endpoint String, t IntoReflectType[Resp], token String('') } $[Resp]
    function Put[Req,Resp] { data IntoReflectValue[Req], endpoint String, t IntoReflectType[Resp], token String('') } $[Resp]
    function Delete[Resp] { endpoint String, t IntoReflectType[Resp], token String('') } $[Resp]
    function Subscribe[Resp] { endpoint String, t IntoReflectType[Resp], token String('') } $[Resp]


