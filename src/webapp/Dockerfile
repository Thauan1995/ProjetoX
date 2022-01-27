FROM golang:1.15.1-alpine3.12
RUN mkdir /app
ADD . /app
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go install
RUN go build -o router .
CMD ["/app/router"]