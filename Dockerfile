FROM golang:1.18

RUN mkdir /app
WORKDIR /app

# setup dependencies, cached
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go install github.com/cosmtrek/air@latest

CMD ["air", "-c", ".air.toml"]
