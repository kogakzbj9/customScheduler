# Use golang:1.22 as the base image
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy the project files into the Docker image
COPY . .

# Build the custom scheduler binary inside the Docker image
RUN go build -o custom-scheduler .

# Set the entrypoint to the custom scheduler binary
ENTRYPOINT ["./custom-scheduler"]
