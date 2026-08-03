FROM golang:{{ go_version }}-alpine AS build

WORKDIR /src

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/{{ project_slug }} ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 10001 app

COPY --from=build /out/{{ project_slug }} /usr/local/bin/{{ project_slug }}
COPY db/migrations /app/db/migrations

USER app
WORKDIR /app

EXPOSE {{ http_port }}

ENTRYPOINT ["/usr/local/bin/{{ project_slug }}"]
