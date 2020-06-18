FROM golang:latest
EXPOSE 8000
RUN apt update && apt install -y libzmq5-dev python3-zmq supervisor
RUN go get github.com/pebbe/zmq4 && go get github.com/gorilla/websocket

RUN mkdir -p /var/log/supervisor
RUN mkdir -p /src/zmq_proxy/static
WORKDIR /src/zmq_proxy
COPY main.go /src/zmq_proxy
COPY static /src/zmq_proxy/static
RUN go build
COPY simple_server.py /src/zmq_proxy
COPY example_services.conf /etc/supervisor/conf.d/example_services.conf
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
