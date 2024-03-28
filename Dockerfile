FROM node:20-slim AS build-frontend

RUN npm i -g pnpm@8.6.3
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

WORKDIR /app
COPY ./frontend .

RUN pnpm install --frozen-lockfile
RUN pnpm install -g @angular/cli
RUN ng build --configuration production

FROM golang:1.22.1-alpine3.19 as build-backend
WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY ./backend/go.mod ./backend/go.sum .
RUN go mod download && go mod verify
COPY ./backend .
RUN go build -o /server

FROM scratch
WORKDIR /app
COPY --from=build-backend /server /app/server
COPY --from=build-frontend /app/dist/browser /app/static
ENTRYPOINT ["/app/server"]
EXPOSE 8080
