FROM golang:latest as builder

WORKDIR /go/src/github.com/kousuke1201abe/go-haiku
RUN apt-get update \
  && apt-get install -y python3 python3-pip git curl wget make xz-utils file sudo unzip \
  && apt-get install -y mecab libmecab-dev mecab-ipadic-utf8 \
  && apt clean
RUN git clone --depth 1 https://github.com/neologd/mecab-ipadic-neologd.git
WORKDIR /go/src/github.com/kousuke1201abe/go-haiku/mecab-ipadic-neologd
RUN ./bin/install-mecab-ipadic-neologd -n -y

WORKDIR /go/src/github.com/kousuke1201abe/go-haiku
RUN go get github.com/slack-go/slack
RUN go get github.com/slack-go/slack/slackevents

ENV CGO_LDFLAGS="-L/path/to/lib -lmecab -lstdc++"
ENV CGO_CFLAGS="-I/path/to/include"
RUN go get github.com/shogo82148/go-mecab
COPY . .
RUN go build main.go
CMD ./main $PORT
