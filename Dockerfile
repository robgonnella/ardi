FROM golang:latest

WORKDIR /arduino-template

RUN mkdir -p data

RUN go get -u github.com/arduino/arduino-cli

ADD arduino-cli.yaml ./

RUN arduino-cli core update-index

RUN arduino-cli core install arduino:avr

ENTRYPOINT ["tail", "-f", "/dev/null"]
