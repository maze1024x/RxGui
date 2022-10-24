RESTful API Request
===================

.. image :: rest.png
    :scale: 62%

The following example program implements a dummy CURD app
based on the fake RESTful API provided by
`JSON Placeholder <https://jsonplaceholder.typicode.com/>`_.

.. literalinclude :: rest.km
    :language: none

We have out-of-the-box support for sending JSON requests to a RESTful API.
To send a request, use the following functions.

* HTTP GET: ``Get[Resp](endpoint, ReflectType)``
* HTTP POST: ``Post[Req,Resp](data, endpoint, ReflectType)``
* HTTP PUT: ``Put[Req,Resp](data, endpoint, ReflectType)``
* HTTP DELETE: ``Delete[Resp](endpoint, ReflectType)``

All the functions above accepts an optional final argument ``token``
representing X-Auth-Token.
They return Observables emitting a value of type ``Resp``
or an error that can be caught by the ``catch`` operator.

.. note::
    ``ReflectType`` is a dummy constant
    that implicitly converts to the type info of ``Resp``.
    This is necessary for JSON deserialization.

.. note::
    The example program uses ListEditView instead of ListView to display a list.
    This is because the fake API does not reflect data changes,
    so that we can't simply use ListView.
    When using a real API, ListView can be used as well.


