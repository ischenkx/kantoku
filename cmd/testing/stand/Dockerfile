# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# ARG to specify the service name during build
ARG SERVICE_NAME

# Copy the compiled binary into the container
COPY ./bin/$SERVICE_NAME/main /app

# Specify the command to run on container startup
CMD ["./main"]