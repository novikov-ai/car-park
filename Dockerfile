FROM golang:1.19

WORKDIR /car-park

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o car-park ./cmd/service

EXPOSE 8080

ENV DATABASE_URL="<place-holder>"

CMD ["./car-park"]