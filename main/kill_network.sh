#!/bin/bash

function kill_network {
    echo "Killing all peer processes..."
    killall -9 p2pnet_auth
    killall -9 p2pnet_onion
    killall -9 p2pnet_nse
    killall -9 p2pnet_rps
    killall -9 p2pnet_gossip
}
function clean {
    echo "Cleaning up executables..."
    rm -f p2pnet_auth
    rm -f p2pnet_onion
    rm -f p2pnet_rps
    rm -f p2pnet_nse
    rm -f p2pnet_gossip
}

kill_network
clean
