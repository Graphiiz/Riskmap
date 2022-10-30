FROM golang:1.19.2-buster AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o server/

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/server /server

ENTRYPOINT ["/server"]