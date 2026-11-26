# Web image: builds the SPA and serves it with Caddy, which also proxies /api
# to the API container. Build context is the repo root (see docker-compose).

FROM node:22-alpine AS build
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM caddy:2-alpine
COPY deploy/Caddyfile /etc/caddy/Caddyfile
COPY --from=build /web/dist /srv
