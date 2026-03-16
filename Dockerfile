FROM golang:1.25-bookworm AS build
RUN apt-get update && apt-get install -y --no-install-recommends make
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make deb \
 && mkdir -p /tmp/pkg \
 && dpkg -x target/deb/xray-vpn_0.1.0_amd64.deb /tmp/pkg

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=build /tmp/pkg/usr/bin/xray-vpnd /usr/bin/xray-vpnd
COPY --from=build /tmp/pkg/usr/bin/xray-vpn  /usr/bin/xray-vpn
COPY --from=build /tmp/pkg/etc/xray-vpn/     /etc/xray-vpn/
RUN mkdir -p /run/xray-vpn /var/cache/xray-vpn
ENV XRAY_LOCATION_ASSET=/var/cache/xray-vpn
ENTRYPOINT ["/usr/bin/xray-vpnd"]
