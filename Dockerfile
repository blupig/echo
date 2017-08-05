# Minimum environment (can also be FROM scratch)
FROM alpine

# Copy built binary
COPY build/echo /app/echo

# Default command
CMD ["/app/echo"]
