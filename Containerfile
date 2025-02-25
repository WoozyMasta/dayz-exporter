FROM docker.io/golang:1.24 AS build

ARG ARCH=amd64
ARG VERSION=0.0.0
ARG COMMIT=nil
ARG UPX_VERSION=4.2.4

# hadolint ignore=DL3008
RUN set -xeu; \
    apt-get update; \
    apt-get install -y --no-install-recommends xz-utils curl; \
    curl -#Lo upx.tar.xz \
        "https://github.com/upx/upx/releases/download/v$UPX_VERSION/upx-$UPX_VERSION-${ARCH}_linux.tar.xz"; \
    tar -xvf upx.tar.xz --strip-components=1 "upx-$UPX_VERSION-${ARCH}_linux/upx"; \
    chmod +x upx; \
    mv upx /usr/local/bin/upx; \
    rm -f upx.tar.xz; \
    go install github.com/tdewolff/minify/cmd/minify@latest

WORKDIR /src/dayz-exporter
COPY go.mod go.sum ./
COPY ./internal/vars ./internal/vars
RUN go mod download

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOFLAGS="-buildvcs=false -trimpath"
ENV GOARCH=$ARCH

COPY ./cli ./cli
COPY ./pkg ./pkg
COPY ./internal ./internal

RUN set -eux;\
    go mod tidy; \
    go generate ./...; \
    go build -ldflags="-s -w \
        -X 'internal/vars.Version=$VERSION' \
        -X 'internal/vars.Commit=$COMMIT' \
        -X 'internal/vars.BuildTime=$(date -uIs)' \
        -X 'internal/vars.URL=https://$(grep -Po 'module \K.*$' go.mod)'" \
      -o "./dayz-exporter" "cli/"*.go; \
    upx --lzma --best ./dayz-exporter; \
    upx -t ./dayz-exporter

FROM scratch
COPY --from=build /src/dayz-exporter/dayz-exporter /dayz-exporter
ENTRYPOINT ["/dayz-exporter"]
