#!/bin/bash
set -e

script() {
    if [ "${TRAVIS_PULL_REQUEST}" == "false" ] && [ "$TRAVIS_OS_NAME" != "windows" ];
    then
        COVERALLS_PARALLEL=true
        go get github.com/axw/gocov/gocov github.com/mattn/goveralls get golang.org/x/tools/cmd/cover
        $HOME/gopath/bin/goveralls -service=travis-ci -repotoken $COVERALLS_TOKEN
    fi

    if [ "$TRAVIS_OS_NAME" != "windows" ];
    then
        go get golang.org/x/lint/golint && golint ./...
        go vet
    fi

    go test -bench=.* -v ./...
}

"$@"