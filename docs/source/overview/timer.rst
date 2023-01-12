Timer
=====

.. image :: stopwatch.png
    :scale: 62%

The following example program implements a simple stopwatch.

.. literalinclude :: stopwatch.rxsc
    :language: none

When writing a timer-related program,
the following APIs are available.

* ``SetTimeout(m)`` returns an Observable that
  starts a single-shot timer when subscribed,
  emitting a Null value ``m`` milliseconds later.
* ``SetInterval(m, n)`` returns an Observable that
  starts a ``n``-shots timer when subscribed,
  emitting increasing numbers up to ``n`` (1, 2, 3, ..., ``n``)
  with an interval of ``m`` milliseconds.
* ``SetInterval(m)`` returns an Observable that
  starts a long-running timer when subscribed,
  emitting increasing numbers (1, 2, 3, ...)
  with an interval of ``m`` milliseconds.
* ``Now`` is an Observable that
  emits the current time when it is subscribed.
* ``TimeOf(o)`` transforms each emission of the Observable ``o``
  into the time of each emission.
* ``o | with-time()`` transforms each emission ``e`` of the Observable ``o``
  into ``Pair(e,t)``, in which ``t`` is the time of each emission.
* ``(t2 -ms t1)`` calculates the time difference between ``t2`` and ``t1``,
  in milliseconds.

Timers created by ``SetTimeout`` and ``SetInterval``
are cleared automatically when the returned Observables are unsubscribed.

.. Caution::
    The behavior of ``SetInterval(m)`` is different from "interval(m)" in RxJS.
    The emission from ``SetInterval(m)`` starts with 1 instead of 0.
    That is, 1 is emitted roughly ``m`` milliseconds after subscription,
    2 is emitted roughly ``(2 * m)`` milliseconds after subscription,
    3 is emitted roughly ``(3 * m)`` milliseconds after subscription, etc.
    This rough correspondence is just for consistency,
    don't rely on it to measure time.


