# OK Run (golang)

ok, run your gofile. run your gofile when `go run` would not.

[![GoDoc](https://godoc.org/github.com/xta/okrun?status.svg)](https://godoc.org/github.com/xta/okrun)

### Features
`okrun` will automatically fix the following error(s):

* *imported and not used* - when your gofile contains an imported package that is not used, `okrun` will comment out the offending import line.

### Setup
    // get okrun
    go get github.com/xta/okrun

    // install okrun
    cd $GOPATH && go install github.com/xta/okrun/

### Usage
    okrun path/to/your/file.go

By using `okrun`, your `file.go` will attempt to run when `go run` would otherwise refuse.

### Example

With file *example.go*:

    package main

    import (
      "fmt"
      "log"
    )

    func main() {
      fmt.Println("hi")
    }

`go run` will not run *example.go*:

    go run example.go
    # command-line-arguments
    example.go:5: imported and not used: "log"

`okrun` will run *example.go*:

    okrun example.go
    hi
    
`okrun` will update your gofile & properly format it. after running `okrun`, *example.go* is now:

    package main

    import (
      "fmt"
      //  "log"
    )

    func main() {
      fmt.Println("hi")
    }


## About

by [Rex Feng](https://twitter.com/rexfeng) 2014

## License

MIT
