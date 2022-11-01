Installation
++++++++++++

Downloading Binary Files
========================

Binary files for x86-64 Linux and Windows are available at the
`GitHub Releases Page
<https://github.com/maze1024x/RxGui/releases>`_.

Note that the Linux version does not ship with Qt library files.
It assumes Qt5 is already installed.

.. Tip::
    Normally, Qt5 should be already installed on a Linux desktop.

Building From Source
====================

Build instructions differ by platform,
as described below in separate subsections.

On all the following platforms,
built files are stored in the ``build`` directory.

Linux
-----

* Build Dependencies: ``make``, ``go`` (or maybe ``golang``, at least 1.18),
  ``gcc``, ``qt5-base`` (or maybe ``qt5-qtbase``, ``qtbase5-dev``, etc.)
* Build Command: ``make``

.. Tip::
    On Debian stable, a standalone installation of Go is required.
    This is because current Debian stable (11 bullseye) only provides Go 1.15,
    which does not support generics.

.. Tip::
    On Debian-derived distros,
    even if Qt5 is already installed as a dependency of system programs,
    it is still necessary to install ``qtbase5-dev``
    to make sure all development files are available.

Windows (MSYS2 MinGW)
---------------------

`MSYS2
<https://www.msys2.org/>`_

* Build Dependencies: ``mingw-w64-x86_64-make``, ``mingw-w64-x86_64-go``,
  ``mingw-w64-x86_64-gcc``, ``mingw-w64-x86_64-qt5-base``
* Build Command: ``mingw32-make``

.. Tip::
    Building with the MSVC variant of Qt is possible,
    but it is NOT supported out-of-the-box by the ``Makefile``,
    and note that CGO may still require GCC.


