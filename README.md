Bridge
======

Bridge is command center for your local network. Designed with extensibility in mind, it allows flexible configuration for your tasks.

For example, with right combination of patterns, it is possible to use bridge for the following:
  * transfer files between computers on local networks;
  * exchange short messages;
  * execute arbitrary shell commands;
  * control services;
  * monitor state;
  * ... much more (provided there are plugins for it).

Security of data exchange and authentication are based on TLS and its certificates.

Currently it is possible to write plugins in Go or in [Gelo](http://code.google.com/p/gelo/); when reference Go implementation will be able to use shared libraries, it will be possible to use shared libraries written in any language.

Notes
-----

Currently the project is in very alpha state. It is not possible to use it for anything useful: plugin interfaces design is in progress, Gelo plugins are not possible now, security system is not implemented. However, this should change more or less soon.
