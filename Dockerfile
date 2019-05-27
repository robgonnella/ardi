FROM golang:latest

WORKDIR /arduino

RUN mkdir -p data

ADD arduino-cli.yaml ./

RUN go get -u github.com/arduino/arduino-cli

RUN arduino-cli core update-index

RUN arduino-cli core install arduino:avr

ENTRYPOINT ["tail", "-f", "/dev/null"]
