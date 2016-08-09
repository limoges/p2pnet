#!/bin/bash

function build_executables {
    echo "Building executables..."
    go build p2pnet_auth.go    
    go build p2pnet_onion.go   
    go build p2pnet_nse.go     
    go build p2pnet_rps.go     
    go build p2pnet_gossip.go  
}

function launch_peer {
    echo "Launching peer based on ${1}."
    ./p2pnet_auth      -f $1   > bootstrap.log 2>&1 & 
    ./p2pnet_onion     -f $1   > default.log   2>&1 &
    ./p2pnet_nse       -f $1   > peer1.log     2>&1 &
    ./p2pnet_rps       -f $1   > peer2.log     2>&1 &
    ./p2pnet_gossip    -f $1   > peer3.log     2>&1 &
}

function build_network {
    echo "Building network..."
    launch_peer configs/bootstrap.ini 
    launch_peer configs/default.ini   
    launch_peer configs/peer1.ini     
    launch_peer configs/peer2.ini     
    launch_peer configs/peer3.ini     
}

build_executables
build_network
