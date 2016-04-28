# transcribe4all [![Build Status](https://travis-ci.org/hack4impact/transcribe4all.svg?branch=master)](https://travis-ci.org/hack4impact/transcribe4all) [![GoDoc](https://godoc.org/github.com/hack4impact/transcribe4all?status.svg)](https://godoc.org/github.com/hack4impact/transcribe4all)

## Go set up the project

```
$ go get github.com/hack4impact/transcribe4all
$ cd $GOPATH/src/github.com/hack4impact/transcribe4all
```
To set up Sphinx for transcription read the following [instructions.](Sphinx/README.md)

## Dependency management

If you are using 1.6 or 1.5 with GO15VENDOREXPERIMENT then the app should just work. If you are using Go < 1.5, you can run

```
$ go get github.com/tools/godep
$ godep restore
```

If you add new dependencies to the app, run

```
$ godep save ./...
```

## Running the app

```
$ go build
$ ./transcribe4all
```

## Tests

Install the testing dependencies

```
$ go get -u github.com/qur/withmock
$ go get -u github.com/qur/withmock/mocktest
$ go get -u golang.org/x/tools/cmd/goimports
```

Run the tests

```
$ make test
```

## License
[MIT License](LICENSE.md)
