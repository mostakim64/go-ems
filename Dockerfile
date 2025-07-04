ARG GO_VERSION=1.24.2

FROM golang:${GO_VERSION}-alpine AS builder

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
RUN apk add --no-cache ca-certificates

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Import the code from the context.
COPY ./ ./

# Build the Go app
RUN CGO_ENABLED=0 GOFLAGS=-mod=vendor GOOS=linux go build -a -o /app .

######## Start a new stage from scratch #######
# Final stage: the running container.
FROM alpine:latest AS final

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/
# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import the open api documentation file from builder stage
# COPY --chown=nobody:nobody --from=builder /src/swagger.yaml swagger.yaml
# Import the compiled executable from the first stage.
COPY --from=builder /app /app

# Perform any further action as an unprivileged user.
USER nobody:nobody

EXPOSE 8080
# Run the compiled binary.
ENTRYPOINT ["/app"]