# go-string

There are many string libraries. This one is mine.

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.6 so let's just assume you need [Go 1.8](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Example

```
package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-string/random"
	"log"
)

func main() {

	var ascii = flag.Bool("ascii", false, "")
	var length = flag.Int("length", 32, "")
	var chars = flag.Int("chars", 0, "")

	flag.Parse()

	opts := random.DefaultOptions()
	opts.ASCII = *ascii
	opts.Length = *length
	opts.Chars = *chars

	s, _ := random.String(opts)

	fmt.Println(s)
}
```

_Error handling omitted for the sake of brevity._
