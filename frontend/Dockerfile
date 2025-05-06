# syntax=docker/dockerfile:1
FROM node:slim AS builder

WORKDIR /app

COPY package.json ./

RUN npm install --force

COPY . .

# bump Node heap to 2 GB so SWC doesn’t get OOM-killed
ENV NODE_OPTIONS="--max-old-space-size=2048"

# build (with your next.config.js adjustments in place)
RUN npm run build

# --- Stage 2: Production image ---
FROM node:slim AS runner

# Set working directory
WORKDIR /app

# Set environment to production
ENV NODE_ENV=production

# Copy only package.json for production install
COPY package.json ./

# Install production dependencies
RUN npm install --omit=dev --force

# Copy build output and public assets from builder
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/next.config.js ./next.config.js

# Expose Next.js default port
EXPOSE 3000

# Start the Next.js application
CMD ["npm", "start"]
