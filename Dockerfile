FROM arm64v8/golang:1.23-alpine as build

WORKDIR /build
COPY . /build
RUN go build -o . ./...

FROM arm64v8/alpine as release

COPY --from=build /build/rpi_exporter /rpi_exporter

EXPOSE 9090

CMD ["sh", "-c", "/rpi_exporter -addr=:9090"]
