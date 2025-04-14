##
## Build
##
FROM --platform=$BUILDPLATFORM crazymax/goxx:1.23 AS base
ENV GO111MODULE=auto
ENV CGO_ENABLED=1
WORKDIR /go/src/hello
LABEL maintainer="Dojima"

FROM base AS build
ARG TARGETPLATFORM
RUN --mount=type=cache,sharing=private,target=/var/cache/apt \
  --mount=type=cache,sharing=private,target=/var/lib/apt/lists \
  goxx-apt-get install -y binutils gcc g++ pkg-config
RUN --mount=type=bind,source=. \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goxx-go build -o /app/remediate ./cmd/main.go
RUN chmod +x /app/remediate

##
## Deploy
##
FROM 908538848727.dkr.ecr.eu-west-3.amazonaws.com/mirrors/gcr.io/distroless/base

WORKDIR /app

COPY --from=build /app/remediate /app/remediate

ENV CONFIG_PATH ""

USER nonroot:nonroot

CMD ["/app/remediate"]

ARG BUILD_DATE
ARG VCS_TYPE=git
ARG VCS_URL
ARG VCS_REF
LABEL build-date=$BUILD_DATE \
    vcs-type=$VCS_TYPE \
    vcs-url=$VCS_URL \
    vcs-ref=$VCS_REF
