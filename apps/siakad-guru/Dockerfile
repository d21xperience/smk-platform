# apps/siakad-guru/Dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
# Copy package files
COPY package.json yarn.lock ./
# Install dependencies
RUN yarn install --frozen-lockfile
# Copy source code
COPY . .
# Build static site
RUN yarn quasar build
# Production stage - Nginx
FROM nginx:alpine
# Copy built assets
COPY --from=builder /app/dist/spa /usr/share/nginx/html
# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
