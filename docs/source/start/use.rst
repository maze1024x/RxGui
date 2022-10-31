Usage
+++++

Writing Code
============

There is a VSCode extension for RxScript available in the marketplace.

`RxScript - Visual Studio Marketplace
<https://marketplace.visualstudio.com/items?itemName=rxgui.rxscript>`_

Note that the "Semantic Highlighting" configuration should be enabled
for the syntax highlighting feature to work.

It is still naive, no other configuration is required.

Run One File
============

To use the framework (interpreter),
run the executable ``rxgui`` or ``rxgui.exe``
with a RxScript source file as the argument.

.. code-block:: console

        $ build/rxgui foo.rxsc

A source file can have a shebang at its beginning.

Run Multiple Files
==================

It is also possible to run multiple source files
specified by a single manifest file.

A manifest file must have a name ending with ``.manifest.json``
and have a JSON content in the following format.

.. code-block:: json

    { "ProjectFiles": [], "DependencyFiles": [] }

``DependencyFiles`` can be also manifest files.

The interpreter executable can run a manifest file directly.

.. code-block:: console

        $ build/rxgui foo.manifest.json

A manifest file can also have a shebang at its beginning.

REPL and Debugging
==================

The interpeter has a companion program called *Naive Debugger*,
which provides the features of a REPL and a debug log viewer.

To enable it, run the interpreter with a ``--debug`` option.

It is also possible to start a standalone REPL.
When running the interpreter without arguments,
a simple prompt is shown as below.

.. code-block:: console

    $ build/rxgui
    Input a source file path or press Enter to start REPL:

A standalone REPL will pop up
after the Enter(Return) key is pressed without any input.

