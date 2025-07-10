#! /bin/sh
set -e

PROJ="leetcode-scrapper"
ORG_PATH="github.com/anwar-arif"
REPO_PATH="${ORG_PATH}/${PROJ}"

if ! [ -x "$(command -v go)" ]; then
    echo "go is not installed"
    exit 1
fi
if ! [ -x "$(command -v git)" ]; then
    echo "git is not installed"
    exit 1
fi
if [ -z "${GOPATH}" ]; then
    echo "set GOPATH"
    exit 1
fi

PATH="${PATH}:${GOPATH}/bin"

go mod verify
go mod vendor
go fmt ./...
go install -v -ldflags="-X ${REPO_PATH}/version.Version=${VERSION}" ./cmd/...