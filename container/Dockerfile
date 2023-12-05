# Copyright 2023 The git-releaser Authors.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.21.4-alpine3.17 AS builder

ENV CGO_ENABLED=0
ARG VERSION
ARG COMMIT
ARG DATE
WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /workspace/git-releaser -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" ./

FROM gcr.io/distroless/static AS production

LABEL org.opencontainers.image.source="https://github.com/git-releaser/git-releaser" \
    org.opencontainers.image.url="https://github.com/git-releaser/git-releaser" \
    org.opencontainers.image.title="git-releaser" \
    org.opencontainers.image.vendor='The git-releaser Authors' \
    org.opencontainers.image.licenses='Apache-2.0'

WORKDIR /
COPY --from=builder /workspace/git-releaser .
USER 65532:65532

CMD ["/git-releaser", "update"]
