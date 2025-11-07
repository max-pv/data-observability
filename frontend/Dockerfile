# Build stage
FROM node:18-alpine AS build

WORKDIR /app

# Copy package manifests first for better build caching
COPY package*.json ./

# Install build dependencies and project dependencies
RUN apk add --no-cache python3 make g++ \
	&& npm install

# Copy the rest of the source code
COPY . .

# Build the app
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy built files from build stage
COPY --from=build /app/build /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
