# Alpine is a Lightweight Linux used for building containerized applications.

FROM golang:1.19-alpine AS build

# Setting the working directory
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY mongo_webserver.go ./
RUN CGO_ENABLED=0 go build -o /bin/mongo_webserver

FROM scratch
COPY --from=build /bin/mongo_webserver /bin/mongo_webserver
ENTRYPOINT ["/bin/mongo_webserver"]