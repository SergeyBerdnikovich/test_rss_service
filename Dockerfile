FROM golang:1.14

RUN apt-get update -y && \
    apt-get install -y fswatch

RUN mkdir /myapp

WORKDIR /myapp

COPY . /myapp

RUN make build

CMD ["make", "run"]
