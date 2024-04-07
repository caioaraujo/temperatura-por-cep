FROM golang:1.22 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o temperatura-por-cep ./cmd/server
EXPOSE 8080

FROM scratch
WORKDIR /app
COPY --from=build /app/temperatura-por-cep .
EXPOSE 8080
ENTRYPOINT [ "./temperatura-por-cep" ]