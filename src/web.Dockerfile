FROM golang:1.15.1-alpine3.12
RUN mkdir /app
ADD . /app
WORKDIR /app/webapp
RUN go mod download
RUN go install
RUN go build -o router .
ENTRYPOINT [ "./router" ]