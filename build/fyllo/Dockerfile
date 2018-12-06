# // Copyright (c) 2018 soren yang
# //
# // Licensed under the MIT License
# // you may not use this file except in complicance with the License.
# // You may obtain a copy of the License at
# //
# //     https://opensource.org/licenses/MIT
# //
# // Unless required by applicable law or agreed to in writing, software
# // distributed under the License is distributed on an "AS IS" BASIS,
# // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# // See the License for the specific language governing permissions and
# // limitations under the License.

# This file is mult-stage build for fyllo Docker image

# build stage
FROM alpine:3.7 AS build-env
ADD . /go/src/github.com/lsytj0413/fyllo
WORKDIR /go/src/github.com/lsytj0413/fyllo
ENV GOPATH /go
RUN export PATH=/go/bin:$PATH \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    ## install musl-dev for gcc
    && apk add --no-cache git go musl-dev make \
    && make build-local \
    # && go build -o /go/src/fyllo -v github.com/lsytj0413/fyllo/cmd/fyllo \
    && mv ./bin/fyllo /go/src/fyllo \
    && rm -rf /var/cache/apk/*

# image stage
FROM alpine:3.7
LABEL MAINTAINER sorenyang@foxmail.com

ENV GIN_MODE release
RUN apk add -U tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=build-env /go/src/fyllo /usr/local/bin/fyllo
EXPOSE 80
WORKDIR /usr/local/bin
CMD ["fyllo"]