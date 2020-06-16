FROM golang:latest
EXPOSE 8000
RUN apt update && apt install -y libzmq5-dev python3-zmq
RUN go get github.com/pebbe/zmq4 && go get github.com/gorilla/websocket

