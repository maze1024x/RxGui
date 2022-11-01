Observable & Subject
++++++++++++++++++++

API
===

Observable
----------

.. code-block:: none

    type $[T] native  // Observable

.. code-block:: none

    function variadic $[T] { items List[T] } $[T]
    function variadic return[T] { items List[T] } $[T]  // equivalent to "$"
    function variadic Merge[T] { items List[$[T]] } $[T]
    function variadic Concat[T] { items List[$[T]] } $[T]
    function StartWith[T] { v T, o $[T] } $[T]
    function WithChildContext[T] { o $[T] } $[T]
    function WithCancelTrigger[T] { cancel $[Null], o $[T] } $[T]
    function WithCancelTimeout[T] { ms Int, o $[T] } $[T]

.. code-block:: none

    function SetTimeout { ms Int } $[Null]
    function SetInterval { ms Int, n Int(-1) } $[Int]

.. code-block:: none

    const UUID $[String]
    function Random { supremum Int } $[Int]
    function Shuffle[T] { l List[T] } $[List[T]]

.. code-block:: none

    const NumCPU Int
    function Go[T] { k Lambda[Null,T] } $[T]
    function ForkJoin[T] { items List[$[T]], n Int(NumCPU) } $[List[T]]
    operator concurrent[T] { l List[$[T]], n Int(NumCPU) } $[T]
    operator concurrent-map[A,B] { o $[A], f Lambda[A,$[B]], n Int(NumCPU) } $[B]
    operator fork-join[T] { l List[$[T]], n Int(NumCPU) } $[List[T]]
    operator fork-join[A,B] { a $[A], b $[B], n Int(NumCPU) } $[Pair[A,B]]

.. code-block:: none

    method $.Result $[Result[T]]
    function Throw[T] { err Error } $[T]
    function Crash[T] { err Error } $[T]
    operator catch[T] { o $[T], f Lambda[Pair[Error,$[T]],$[T]] } $[T]
    operator retry[T] { o $[T], n Int(-1) } $[T]
    operator log-error[T] { o $[T] } $[T]

.. code-block:: none

    operator distinct-until-changed[T] { o $[T] } { T/== ==[T] } $[T]
    operator with-latest-from[T,X] { o $[T], x $[X] } $[Pair[T,X]]
    operator map-to-latest-from[X] { o $[Null], x $[X] } $[X]
    operator with-cycle[T,X] { o $[T], l List[X] } $[Pair[T,X]]
    operator with-index[T] { o $[T] } $[Pair[T,Int]]
    operator delay-subscription[T] { o $[T], ms Int } $[T]
    operator delay-values[T] { o $[T], ms Int } $[T]
    operator start-with[T] { o $[T], item T } $[T]
    operator end-with[T]   { o $[T], item T } $[T]
    operator throttle[T] { o $[T], f Lambda[T,$[Null]] } $[T]
    operator debounce[T] { o $[T], f Lambda[T,$[Null]] } $[T]
    operator throttle-time[T] { o $[T], ms Int } $[T]
    operator debounce-time[T] { o $[T], ms Int } $[T]
    operator complete-on-emit[T] { o $[T] } $[Null]
    operator skip[T] { o $[T], n Int } $[T]
    operator take[T] { o $[T], n Int } $[T]
    operator take-last[T] { o $[T] } $[T]
    operator take-last?[T] { o $[T] } $[Maybe[T]]
    operator take-while[T] { o $[T], f Lambda[T,Bool] }
    operator take-while?[T] { o $[Maybe[T]] } $[T]
    operator take-until[T] { o $[T], stop $[Null] } $[T]
    operator count[T] { o $[T] } $[Int]
    operator collect[T] { o $[T], n Int(-1) } $[List[T]]
    operator buffer-time[T] { o $[T], ms Int } $[List[T]]
    operator pairwise[T] { o $[T] } $[Pair[T,T]]
    operator buffer-count[T] { o $[T], n Int } $[Queue[T]]
    operator map[A,B] { o $[A], f Lambda[A,B] } $[B]
    operator map-to[A,B] { o $[A], v B } $[B]
    operator filter[T] { o $[T], f Lambda[T,Bool] } $[T]
    operator scan[A,B] { o $[A], v B, f Lambda[Pair[B,A],B] } $[B]
    operator reduce[A,B] { o $[A], v B, f Lambda[Pair[B,A],B] } $[B]
    operator combine-latest[A,B] { a $[A], b $[B] } $[Pair[A,B]]
    operator combine-latest[T] { l List[$[T]] } $[List[T]]
    operator await[A,B] { o $[A], k Lambda[A,$[B]] } $[B]
    operator await-noexcept[A,B] { o $[A], k Lambda[A,$[B]] } $[B]
    operator then[T] { o $[Null], k $[T] } $[T]
    operator with[T] { o $[T], bg $[Null] } $[T]
    operator and[T] { o $[T], bg $[Null] } $[T]
    operator auto-map[A,B] { o $[A], f Lambda[A,$[B]] } $[B]
    operator merge[T]  { l List[$[T]] } $[T]
    operator concat[T] { l List[$[T]] } $[T]
    operator merge[T]  { o1 $[T], o2 $[T] } $[T]
    operator concat[T] { o1 $[T], o2 $[T] } $[T]
    operator merge-map[A,B]   { o $[A], f Lambda[A,$[B]] } $[B]
    operator concat-map[A,B]  { o $[A], f Lambda[A,$[B]] } $[B]
    operator switch-map[A,B]  { o $[A], f Lambda[A,$[B]] } $[B]
    operator exhaust-map[A,B] { o $[A], f Lambda[A,$[B]] } $[B]

Subject
-------

.. code-block:: none

    type Subject[T] native

.. code-block:: none
    
    function variadic CreateSubject[T] { replay Int(0), items List[T] } $[Subject[T]]

.. code-block:: none

    method Subject.Values $[T]
    method Subject.$ $[T]  // equivalent to "Values"

.. code-block:: none

    operator plug[T] { s Subject[T], o $[T] } $[Null]
    operator push [T] { s Subject[T], v T } $[Null]
    operator << [T] { s Subject[T], o $[T] } $[Null]  // equivalent to "plug"
    operator <- [T] { s Subject[T], v T } $[Null]  // equivalent to "push"

.. code-block:: none

    function Multicast[T] { o $[T] } $[$[T]]
    function Loopback[T] { k Lambda[$[T],$[T]] } $[T]
    function SkipSync[T] { o $[T] } $[T]

.. code-block:: none

    method Subject.SampleValue $[T]
    function Sample[T] { o $[T] } $[T]


