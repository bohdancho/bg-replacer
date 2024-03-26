FROM golang:1.22

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
RUN apt-get update && apt-get install -y libwebp-dev

COPY . .
RUN go build -o /usr/local/bin/app 

EXPOSE 8080

CMD ["app"]
