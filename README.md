# Simple zmq->websock pubsub forwarder

Takes a zmq pub socket (there's a simple exmaple provided in simple_server.py)
and forwards all messages on to websocket subscribers.

## Prerequisites

You'll need libzmq-dev on the C side, and then python3-zmq for the python stuff.
There's a dockerfile that works if you're on Windows.

## Use

For the simple server:

    python3 simple_server.py

This will start publishing the time (with no subject) to port 5000.

For the webserver:

    go build -o zmq_fwd
    ./zmq_fwd

This will start a webserver on port 8000. Open a browser to http://localhost:8000
and an application that displays whatever the python producer sends it will start.