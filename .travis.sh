#!/bin/bash
set -e

script() {
    if [ "${TRAVIS_PULL_REQUEST}" == "false" ] && [ "$TRAVIS_OS_NAME" != "windows" ];
    then
        COVERALLS_PARALLEL=true
        go get github.com/axw/gocov/gocov github.com/mattn/goveralls golang.org/x/tools/cmd/cover
        $HOME/gopath/bin/goveralls -service=travis-ci -repotoken $COVERALLS_TOKEN

        # add js coverage
        if [ "$TRAVIS_OS_NAME" == "linux" ];
        then
            export PATH="$PATH:$(go env GOROOT)/misc/wasm"
            GOOS=js GOARCH=wasm go test -covermode=count -coverprofile=profile.cov
            $HOME/gopath/bin/goveralls -coverprofile=profile.cov -service=travis-ci -repotoken $COVERALLS_TOKEN
        fi
    fi

    go get golang.org/x/lint/golint && golint ./...
    go vet
    go test -bench=.* -v ./...
}

"$@"