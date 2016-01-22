detox
=====

detox helps detoxify the output from programs by removing backspace
sequences (à la ``col -b``) and ANSI control codes (e.g. colour codes).

If no arguments are given, ``detox`` will read from standard
input. Otherwise, it detoxifies the files specified as arguments. If
the ``-i`` flag is given, ``detox`` will modify files in place; otherwise,
it emits the detoxified version to standard output.

This program is released under the MIT license; see the LICENSE file.

