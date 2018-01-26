# Build container
FROM golang:alpine
COPY . /src
WORKDIR /src

# --build-arg
ARG SOURCE_COMMIT

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/echo -a -ldflags "-X main.src_commit=${SOURCE_COMMIT} -extldflags \"-static\"" .

# Minimum runtime container (can also be FROM scratch)
FROM alpine

# Copy built binary from build container
COPY --from=0 /build/echo /bin/echo

# Default command
CMD ["/bin/echo"]
