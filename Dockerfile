FROM golang:1.16.0-alpine AS base
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM base AS build
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /out/istio-ingress-migrate .

FROM alpine

COPY --from=build /out/istio-ingress-migrate /usr/bin/istio-ingress-migrate

ENTRYPOINT ["/usr/bin/istio-ingress-migrate"]
