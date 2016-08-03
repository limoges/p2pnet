Peer-to-peer network & security - Summer 2016
===
by Julien Limoges

## Dependencies

The project requires a standard installation of the Go (https://golang.org/dl/).

Installation instructions can be found [here](https://golang.org/doc/install)

This project has the following dependencies:

    github.com/vaughan0/go-ini


## Building & running modules

Modules are found in the module/ directory. Each modules can be ran independently.
For example, to start the Onion Forwarding module, one can run:

    go run modules/p2p_onion.go
    go run modules/p2p_auth.go

Otherwise, you can build the executables using:

    go build modules/p2p_onion.go
    go build modules/p2p_auth.go

Modules may depend on each other to work properly.

## Generating the necessary RSA key
Simply run the following command and follow the instructions.
Note: Onion Authentication does not currently support password protected private
keys.

    ssh-keygen -t rsa -b 4096 -m PEM

## Supported features
### Onion Forwarding
- ONION_TUNNEL_BUILD
- ONION_TUNNEL_READY
- ONION_TUNNEL_INCOMING
- ONION_TUNNEL_DESTROY
- ONION_TUNNEL_DATA
- ONION_ERROR
- ONION_COVER

### Onion Authentication
- AUTH_SESSION_START
- AUTH_SESSION_HS1
- AUTH_SESSION_INCOMING_HS1
- AUTH_SESSION_HS2
- AUTH_SESSION_INCOMING_HS2
- AUTH_LAYER_ENCRYPT
- AUTH_LAYER_ENCRYPT_RESP
- AUTH_LAYER_DECRYPT
- AUTH_LAYER_DECRYPT_RESP
- AUTH_SESSION_CLOSE

### Extensions
- AUTH_HANDSHAKE_REQUEST
- AUTH_HANDSHAKE_RESPONSE

