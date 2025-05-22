FROM golang:1.24.3-alpine3.21 as build

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/etc/apk/cache apk add --update-cache make
WORKDIR /build
COPY . /build
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} make build

FROM alpine:3.21 as release

COPY --from=build /build/bin/rpi_exporter /rpi_exporter

EXPOSE 9090

CMD ["sh", "-c", "/rpi_exporter", "-addr", ":9090"]
