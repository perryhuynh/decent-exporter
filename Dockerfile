# syntax=docker/dockerfile:1.25

ARG GO_VERSION=1.26.4
ARG VERSION=dev
ARG REVISION=unknown

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG REVISION
WORKDIR /src
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -buildvcs=false \
    -trimpath \
    -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${REVISION} -X main.date=unknown" \
    -o /out/decent-exporter ./cmd/decent-exporter

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /out/decent-exporter /decent-exporter
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/decent-exporter"]
