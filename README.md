times 
==========

[![GoDoc](https://godoc.org/github.com/djherbis/times?status.svg)](https://godoc.org/github.com/djherbis/times)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.txt)
[![Build Status](https://travis-ci.org/djherbis/times.svg?branch=master)](https://travis-ci.org/djherbis/times)
[![Coverage Status](https://coveralls.io/repos/djherbis/times/badge.svg?branch=master)](https://coveralls.io/r/djherbis/times?branch=master)

Usage
------------
File Times for #golang

Go has a hidden time functions for most platforms, this repo makes them accessible.

```go
package main

import (
  "log"

  "github.com/djherbis/times"
)

func main() {
  t, err := times.Stat("myfile")
  if err != nil {
    log.Fatal(err.Error())
  }

  log.Println(t.AccessTime())
  log.Println(t.ModTime())

  if ct, ok := t.ChangeTime(); ok {
    log.Println(ct)
  }

  if bt, ok := t.BirthTime(); ok {
    log.Println(bt)
  }
}
```

Installation
------------
```sh
go get github.com/djherbis/times
```
