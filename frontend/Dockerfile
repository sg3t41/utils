
# Dependency installation stage
FROM node:20-slim AS deps
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install

# Build stage
FROM node:20-slim AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# Production stage
FROM node:20-slim AS runner
WORKDIR /app

ENV NODE_ENV=production

COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

EXPOSE 3000

CMD ["npm", "start"]
