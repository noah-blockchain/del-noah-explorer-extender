# docker build -t noah-extender:latest -f DOCKER/Dockerfile .
# docker run -d -p 3000:3000 -p 127.0.0.1:9000:9000 --restart=always /node/node/:/root/noah/ node:latest

FROM golang:1.12-buster as builder

ENV APP_PATH /home/noah-extender

COPY . ${APP_PATH}
COPY ./DOCKER/cmd.sh /usr/local/bin/extenderstart.sh

WORKDIR ${APP_PATH}

RUN make create_vendor && \
    make build && \
    ls -la ./builds/

EXPOSE 9000
CMD ["/usr/local/bin/extenderstart.sh"]
STOPSIGNAL SIGTERM
