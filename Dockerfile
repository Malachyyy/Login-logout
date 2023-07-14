# Official golang image
FROM golang:1.20.6-alpine3.18

# Work directory for application
WORKDIR /app

# Copy dependencies and install them
COPY go.mod go.sum ./
RUN go mod download

# Copy go files
COPY *.go ./

# Make directory for html templates and copy into them
RUN mkdir /app/templates
COPY templates /app/templates

# Build docker image with the name webapp
RUN CGO_ENABLED=0 GOOS=linux go build -o /webapp

# Expose app with port 8080
EXPOSE 8080

# Run the app
CMD ["/webapp"]