FROM golang:1.27-alpine AS build

RUN apk add --no-cache ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/service ./cmd/${SERVICE}

FROM scratch

ARG SERVICE

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /out/service /service
COPY configs/${SERVICE} /configs/${SERVICE}

ENV CONFIG_PATH=/configs/${SERVICE}/prod.yaml

USER 65532:65532

ENTRYPOINT ["/service"]
