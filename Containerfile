FROM docker.io/golang:1.23 AS build

ARG ARCH=amd64
ARG VERSION=0.0.0
ARG COMMIT=nil
ARG UPX_VERSION=4.2.4

RUN set -xeu; \
    curl -#Lo upx.tar.xz \
        "https://github.com/upx/upx/releases/download/v$UPX_VERSION/upx-$UPX_VERSION-${ARCH}_linux.tar.xz"; \
    tar -xvf upx.tar.xz --strip-components=1 "upx-$UPX_VERSION-${ARCH}_linux/upx"; \
    chmod +x upx; \
    mv upx /usr/local/bin/upx; \
    rm -f upx.tar.xz

WORKDIR /src/dayz-exporter
COPY go.mod go.sum ./
RUN go mod download

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOFLAGS="-buildvcs=false -trimpath"
ENV GOARCH=$ARCH

COPY cli pkg .
RUN set -eux;\
    pkg="$(grep -Po 'module \K.*$' go.mod)/pkg/config"; \
    version="$pkg.Version=$VERSION"; \
    commit="$pkg.Commit=$COMMIT"; \
    date="$pkg.BuildTime=$(date -uIs)"; \
    go build -ldflags="-s -w -X '$version' -X '$commit' -X '$date'" \
      -o "./dayz-exporter" "cli/"*.go; \
    upx --lzma --best ./dayz-exporter; \
    upx -t ./dayz-exporter

FROM scratch
COPY --from=build /src/dayz-exporter/dayz-exporter /dayz-exporter
ENTRYPOINT ["/dayz-exporter"]
