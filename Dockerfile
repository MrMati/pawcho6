# syntax=docker/dockerfile:1.7

FROM golang:1.24-alpine AS stage1

WORKDIR /src

ARG VERSION=1.0.0
ARG REPO_SSH_URL=git@github.com:MrMati/pawcho6.git

RUN apk add --no-cache git openssh-client

RUN mkdir -p -m 0700 /root/.ssh && \
    ssh-keyscan github.com >> /root/.ssh/known_hosts

COPY docker/start.sh /start.sh
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

RUN --mount=type=ssh git clone "$REPO_SSH_URL" /src/repo

RUN chmod +x /start.sh && \
    echo "$VERSION" > /src/VERSION && \
    go build -ldflags="-s -w -X main.version=$VERSION" -o /server /src/repo/app/main.go


FROM nginx:alpine AS stage2

WORKDIR /app

RUN apk add --no-cache wget

COPY --from=stage1 /start.sh /start.sh
COPY --from=stage1 /server /app/server
COPY --from=stage1 /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --quiet --spider http://127.0.0.1/ || exit 1

CMD ["/start.sh"]
