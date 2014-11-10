csgolangpinger
==============

A stand-alone executable that checks the latency of Counter-Strike: Global Offensive matchmaking servers. 
The title is intended as a play on words meant to be read as CS:GO(lang) Pinger.

The user interface is written using QML and the qml-go package available [here](http://github.com/go-qml/go-qml). 
The underlying mechanism to ping servers relies on the ping utilities built in to most operating systems,
so if for some reason those don't exist or they've been symlinked to some other utility, it won't work. 

Building and running this should be relatively easy and can _probably_ be done by doing the following:

```
  go get github.com/rue92/csgolangpinger
  go build csgolangpinger
```

I'm not entirely sure though because I'm too lazy to do it myself. Of course, you'll need whatever dependencies are necessary
for the go-qml package (like the qt-dev libraries).

Roadmap
-------

- Make the server list more extensible so if the IP addresses change then they can be changed without recompilation.
- Allow the ping column to be sortable, mostly for fun, though this isn't useful to most people since they'll be matchmaked to their region anyway.
- When (and if) Go ever supports ICMP natively without needing raw sockets (a feature Windows lacks unless an application is run as Admin), then the ping bit can be rewritten.
