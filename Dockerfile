FROM golang:1.22 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./temperatura-por-cep

FROM scratch
WORKDIR /app
COPY --from=build /app/temperatura-por-cep .
ENTRYPOINT [ "./temperatura-por-cep" ]