FROM --platform=${BUILDPLATFORM} golang AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd/
COPY pkg pkg/

ARG TARGETPLATFORM
ENV GOOS=${TARGETPLATFORM%/*}
ENV GOARCH=${TARGETPLATFORM#*/}

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o /go/bin/llb ./cmd/llb

FROM scratch
COPY --from=builder /go/bin/llb /go/bin/llb
ENTRYPOINT ["/go/bin/llb"]