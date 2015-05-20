# pecolify

Transform your Golangs' data to peco

## Installation

Bundle your program

```
$ go get github.com/ysugimoto/pecolify
```

## Usage

Instanciate and pass the Golang's String-Slice `[]string`:

```
package main

import (
    "github.com/ysugimoto/pecolify"
    "fmt"
)

func main() {
    // Instanciate pecolify
    pf := pecolify.New()

    // Pass data to pecolify
    data := []string{
        "foo",
        "bar",
        "baz",
    }

    // pecolify!
    selected, err := pf.Transform(data)
    if err != nil {
        fmt.Printf("Error was occured: %v\n", err)
        return
    }

    fmt.Printf("Selected from peco: %s\n", selected)
}
```

See example : https://github.com/ysugimoto/pecolify/blob/master/example/main.go

## Author

Yoshiaki Sugimoto

## License

MIT

