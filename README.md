# FreeSMS

[![Go Reference](https://pkg.go.dev/badge/github.com/tamboto2000/freesms.svg)](https://pkg.go.dev/github.com/tamboto2000/freesms)

FreeSMS is an API for sending SMS to all operator in Indonesia free of charge.

### Features
 - Free SMS Sending!
 - Send SMS with proxy

### Instalation

FreeSMS require Go v1.14 or up

```sh
$ GO111MODULE=on go get github.com/tamboto2000/freesms
```

### Example

```go
package main

import (
	"github.com/tamboto2000/freesms"
)

func main() {
	cl, err := freesms.NewClient()
	if err != nil {
		panic(err.Error())
	}

	// set proxy (optional)
	if err := cl.SetProxy("http://103.157.116.199:8080"); err != nil {
		panic(err.Error())
	}

    // Minimum chars is 15, maximum 122
	if err := cl.SendMsg("08xxxxxxx", "Lorem ipsum dolor sit amet, uuh, and stuff..."); err != nil {
		panic(err.Error())
	}
}

```

See [documentation](https://pkg.go.dev/github.com/tamboto2000/freesms) for more detailed info

License
---------

MIT