#!/bin/bash
set -e

script() {
    if [ "${TRAVIS_PULL_REQUEST}" == "false" ] && [ "$TRAVIS_OS_NAME" != "windows" ];
    then
        COVERALLS_PARALLEL=true
        go get github.com/axw/gocov/gocov github.com/mattn/goveralls golang.org/x/tools/cmd/cover
        # go test -covermode=count -coverprofile=profile.cov

        # add js coverage
        if [ "$TRAVIS_OS_NAME" == "linux" ];
        then
            bash js.cover.sh
        fi

        # PROFILES=`ls -dm profile.cov*`
        # PROFILES=${PROFILES// /}
        $HOME/gopath/bin/goveralls -coverprofile=profile.cov.js -service=travis-ci -repotoken $COVERALLS_TOKEN
    fi

    go get golang.org/x/lint/golint && golint ./...
    go vet
    go test -bench=.* -v ./...
}

"$@"