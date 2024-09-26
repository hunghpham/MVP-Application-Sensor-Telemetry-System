# Build the Svelte app
FROM node:16 as builder
WORKDIR /app
COPY . .
RUN npm install
RUN npm run build

# Serve the built app using a simple HTTP server
FROM node:16
WORKDIR /app
COPY --from=builder /app/dist /app
RUN npm install -g serve
CMD ["serve", "-s", ".", "-l", "5000"]
