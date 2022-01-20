FROM golang:latest AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o server main.go

FROM scratch
COPY static /
COPY --from=build /app/server /
EXPOSE 8080
ENTRYPOINT [ "/server" ]
