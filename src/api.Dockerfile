FROM golang:1.15.1-alpine3.12
RUN mkdir /api
ADD . /api
WORKDIR /api/site
RUN go mod download
RUN go install
RUN go build -o router .
ENTRYPOINT [ "./router" ]