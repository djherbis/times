script() {
    if [ "${TRAVIS_PULL_REQUEST}" == "false" ] && [ "$TRAVIS_OS_NAME" != "windows" ];
    then
        COVERALLS_PARALLEL=true
        go get github.com/axw/gocov/gocov github.com/mattn/goveralls
        if ! go get code.google.com/p/go.tools/cmd/cover; then go get golang.org/x/tools/cmd/cover;  fi
        $HOME/gopath/bin/goveralls -service=travis-ci -repotoken $COVERALLS_TOKEN
    fi

    go get golang.org/x/lint/golint && golint ./...
    go vet
    go test -bench=.* -v ./...
}

"$@"