# Use busybox as the base image
FROM busybox

# Set the working directory inside the container
WORKDIR /app

# Add the pre-built custom scheduler binary to the Docker image
ADD custom-scheduler /app/custom-scheduler

# Set the entrypoint to the custom scheduler binary
ENTRYPOINT ["./custom-scheduler"]
