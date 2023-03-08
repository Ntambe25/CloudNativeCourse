# Alpine is a Lightweight Linux used for building containerized applications.
FROM golang:1.19-alpine AS build

# Setting the working directory
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY movieapi/ movieapi/
COPY movieserver/ movieserver/
WORKDIR /app/movieserver
RUN CGO_ENABLED=0 go build -o /bin/movieserver

FROM scratch
COPY --from=build /bin/movieserver /bin/movieserver
ENTRYPOINT ["/bin/movieserver"]