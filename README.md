# audio-transcription-service [![Build Status](https://travis-ci.org/hack4impact/audio-transcription-service.svg?branch=master)](https://travis-ci.org/hack4impact/audio-transcription-service) [![GoDoc](https://godoc.org/github.com/hack4impact/audio-transcription-service?status.svg)](https://godoc.org/github.com/hack4impact/audio-transcription-service)

## Go set up the project

```
$ go get github.com/hack4impact/audio-transcription-service
$ cd $GOPATH/src/github.com/hack4impact/audio-transcription-service
```

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
$ ./audio-transcription-service
```

## Running tests

```
$ make test
```

## License
[MIT License](LICENSE.md)
