#!/bin/bash
set -e

script() {
    if [ "${TRAVIS_PULL_REQUEST}" == "false" ];
    then
        COVERALLS_PARALLEL=true

        if [ ! -z "$JS" ];
        then
            bash js.cover.sh
        else
            go test -covermode=count -coverprofile=profile.cov
        fi

        go get github.com/axw/gocov/gocov golang.org/x/tools/cmd/cover
        go install github.com/mattn/goveralls@latest
        $GOPATH/bin/goveralls --coverprofile=profile.cov -service=travis-ci
    fi

    if [ -z "$JS" ];
    then
        go install honnef.co/go/tools/cmd/staticcheck@latest && $GOPATH/bin/staticcheck ./...
        go vet
        go test -bench=.* -v ./...
    fi
}

"$@"
