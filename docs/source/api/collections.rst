Collections
+++++++++++

API
===

List
----

.. code-block:: none

    type List[T] native

.. code-block:: none

    function variadic List[T] { items List[T] } List[T]
    function variadic ListConcat[T] { lists List[List[T]] } List[T]

.. code-block:: none

    function Cons[T] { head T, tail List[T] } List[T]
    function Count { n Int } List[Int]

.. code-block:: none

    method List.Empty Bool
    method List.First Maybe[T]
    method List.Length Int
    method List.Seq Seq[T]
    method List.$ $[T]

.. code-block:: none

    operator shift[T]   { l List[T] } Maybe[Pair[T,List[T]]]
    operator prepend[T] { l List[T], item T } List[T]
    operator reverse[T] { l List[T] } List[T]
    operator sort[T] { l List[T] } { T/< <[T] } List[T]
    operator take[T] { l List[T], n Int } List[T]
    operator with-index[T] { l List[T] } List[Pair[T,Int]]
    operator map[A,B]  { l List[A], f Lambda[A,B] } List[B]
    operator map?[A,B] { l List[A], f Lambda[A,Maybe[B]] } List[B]
    operator map*[A,B] { l List[A], f Lambda[A,List[B]] } List[B]
    operator filter[T] { l List[T], f Lambda[T,Bool] } List[T]
    operator scan[A,B] { l List[A], v B, f Lambda[Pair[B,A],B] } List[B]
    operator fold[A,B] { l List[A], v B, f Lambda[Pair[B,A],B] } B

Seq
---

.. code-block:: none

    type Seq[T] native

.. code-block:: none

    function variadic Seq[T] { items List[T] } Seq[T]

.. code-block:: none

    method Seq.Empty Bool
    method Seq.Last Maybe[T]
    method Seq.Length Int
    method Seq.List List[T]

.. code-block:: none

    operator append[T] { s Seq[T], item T } Seq[T]
    operator append?[T] { s Seq[T], item? Maybe[T] } Seq[T]
    operator append*[T] { s Seq[T], items List[T] } Seq[T]
    operator sort[T] { s Seq[T] } { T/< <[T] } Seq[T]
    operator filter[T] { s Seq[T], f Lambda[T,Bool] } Seq[T]

Queue
-----

.. code-block:: none

    type Queue[T] native

.. code-block:: none

    function variadic Queue[T] { items List[T] } Queue[T]

.. code-block:: none

    method Queue.Empty Bool
    method Queue.Size  Int
    method Queue.First Maybe[T]
    method Queue.List List[T]

.. code-block:: none

    operator shift[T]  { q Queue[T] } Maybe[Pair[T,Queue[T]]]
    operator append[T] { q Queue[T], item T } Queue[T]

Heap
----

.. code-block:: none

    type Heap[T] native

.. code-block:: none

    function variadic Heap[T] { items List[T] } { T/< <[T] } Heap[T]

.. code-block:: none

    method Heap.Empty Bool
    method Heap.Size  Int
    method Heap.First Maybe[T]
    method Heap.List List[T]

.. code-block:: none

    operator shift[T]  { h Heap[T] } Maybe[Pair[T,Heap[T]]]
    operator insert[T] { h Heap[T], item T } Heap[T]

Set
---

.. code-block:: none

    type Set[T] native

.. code-block:: none

    function variadic Set[T] { items List[T] } { T/<> <>[T] } Set[T]

.. code-block:: none

    method Set.Empty Bool
    method Set.Size  Int
    method Set.List List[T]

.. code-block:: none

    operator has[T] { s Set[T], item T } Bool
    operator delete[T] { s Set[T], item T } Set[T]
    operator insert[T] { s Set[T], item T } Set[T]

Map
---

.. code-block:: none

    type Map[K,V] native

.. code-block:: none

    function variadic Map[K,V] { entries List[Pair[K,V]] } { K/<> <>[K] } Map[K,V]

.. code-block:: none

    method Map.Empty Bool
    method Map.Size  Int
    method Map.Keys    List[K]
    method Map.Values  List[V]
    method Map.Entries List[Pair[K,V]]

.. code-block:: none

    operator has[K,V] { m Map[K,V], key K } Bool
    operator lookup[K,V] { m Map[K,V], key K } Maybe[V]
    operator delete[K,V] { m Map[K,V], key K } Map[K,V]
    operator insert[K,V] { m Map[K,V], pair Pair[K,V] } Map[K,V]


