Coprocess & Pipe Protocol
=========================

.. image :: coproc.png
    :scale: 62%

The functions ``Get`` ``Post`` ``Put`` ``Delete`` can send requests via HTTP,
but they can also send requests via pipe, to a coprocess.

To send requests via pipe, the endpoint URL should start with ``pipe://stdio/``.
Request text is sent to standard output,
and response text is expected to be received from standard input.

Request text looks like this.

.. code-block:: none
    
    REQ 1 GET / "" 0

.. code-block:: none

    REQ 2 POST / "token" 2
    42

Response text looks like this.

.. code-block:: none

    OK 1 13
    "Hello World"

.. code-block:: none

    ERR 2 11
    bad request

If a returned Observable is unsubscribed
before its corresponding response is received,
a cancel signal is sent.

.. code-block:: none

    CANCEL 3

Pipe requests support an additional method called SUBSCRIBE,
which can have multiple OK responses over time.
To send a SUBSCRIBE request, call the ``Subscribe`` function,
which has a function signature identical to ``Get``.


