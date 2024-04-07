FROM golang:1.22.1-alpine3.19 as build-backend
WORKDIR /usr/src/app

ENV CGO_ENABLED=1
RUN apk add build-base

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum .
RUN go mod download && go mod verify
COPY . .
RUN go build -o /app

FROM alpine:3.19
COPY --from=build-backend /app /app
ENTRYPOINT ["/app"]
# CMD ["sleep", "infinity"]
EXPOSE 8080
