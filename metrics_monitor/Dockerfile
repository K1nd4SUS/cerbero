# Use the official Nginx image as the base image
FROM nginx:latest

# Copy the static content to the default Nginx web root
COPY src /usr/share/nginx/html

# Expose port 80 to the outside world
EXPOSE 80

# Command to start Nginx when the container runs
CMD ["nginx", "-g", "daemon off;"]
