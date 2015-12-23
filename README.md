# Lucky 13

This is a (WIP) research implementation of the Lucky 13 attack on TLS. Ideally,
this code will be able to be pointed at a given TLS server and tell you if that
implementation is vulnerable to Lucky 13

## How to use it

For now, if you go to `cmd/lucky13` and `cmd/tls-server`, and `$ go build`, you
should have `lucky13` and `tls-server` executables. The server is simply a TLS
server implemented in golang for testing purposes. Run the server first, then
the `lucky13` client, and it will show you timing information.

NOTE: the actual attack is not yet implemented, currently this implementation is
collecting timing information to gauge the feasibility of the attack.
