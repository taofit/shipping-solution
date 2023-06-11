FROM golang:alpine 
# Add a work directory
WORKDIR /app
# Cache and install dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download
# Copy app files
COPY . .
# Install Reflex for development
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/shippingWeb.go
# Expose port
EXPOSE 8082
# Start app
CMD ["./shippingWeb"]