# API image: the api and migrate binaries on a distroless base (no shell, no
# package manager), running as a non-root user.

FROM golang:1.27-alpine AS build
WORKDIR /src
# Download modules in their own layer so code changes do not refetch them.
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY db ./db
COPY internal ./internal
# CGO off: a static binary that needs no libc in the runtime image.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api \
 && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /out/api /out/migrate /usr/local/bin/
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/api"]
