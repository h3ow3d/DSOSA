FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
COPY web ./web
RUN go build -o /out/dsovs ./cmd/dsovs

FROM alpine:3.20
RUN adduser -D -h /app app && mkdir -p /data && chown -R app:app /data
WORKDIR /app
COPY --from=builder /out/dsovs /app/dsovs
COPY --from=builder /src/web /app/web
USER app
EXPOSE 8080
ENV APP_ADDR=:8080 DATA_DIR=/data DSOVS_URL=https://owasp.org/www-project-devsecops-verification-standard/dist/dsovs.json AUTO_SYNC_CATALOGUE=false
ENTRYPOINT ["/app/dsovs"]
