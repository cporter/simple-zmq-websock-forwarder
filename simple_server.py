import zmq
import sys
import json
import time

def main():
    ctx = zmq.Context()
    sock = ctx.socket(zmq.PUB)
    sock.bind('tcp://*:5000')

    while True:
        sock.send_multipart([b'', json.dumps({'time' : time.time()}).encode('ascii')])
        time.sleep(1.0)

if '__main__' == __name__:
    main()        