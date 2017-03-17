rollbar
-------

`rollbar` is a Go Rollbar client that makes it easy to report errors to Rollbar
with stacktraces. Errors are sent to Rollbar asynchronously in a background
goroutine.

Because Go's `error` type doesn't include stack information from when it was set
or allocated, `rollbar` uses the stack information from where the error was
reported.

You may also want to look at:

* [stvp/roll](https://github.com/stvp/roll) - Simpler, synchronous (no
  background goroutine) with a nicer API.

Installation
=============

Assumming you have Go properly installed, as well as the proper workspace file structure, run the following command from your working project directory:

```
go get github.com/stvp/rollbar
```

Documentation
=============

[API docs on godoc.org](http://godoc.org/github.com/stvp/rollbar)

Usage
=====

```go
package main

import (
  // depending on your text edtitor / IDE / plugins,  some do not allow unused import statements, so make sure to add the import, make the rollbar calls below, and then save the file or else the line might automatically get removed
  "github.com/stvp/rollbar"
)

func main() {
  rollbar.Token = "MY_TOKEN"
  rollbar.Environment = "production" // defaults to "development"

  result, err := DoSomething()
  if err != nil {
    rollbar.Error(rollbar.ERR, err) // call this anywhere necessary in your code, along with rollbar.Wait() if needed
  }

  rollbar.Message("info", "Message body goes here")

  rollbar.Wait()
}
```

Running Tests
=============

Set up a dummy project in Rollbar and pass the access token as an environment
variable to `go test`:

    TOKEN=f0df01587b8f76b2c217af34c479f9ea go test

And verify the reported errors manually in the Rollbar dashboard.

Contributors
============

Thanks, all!

* @kjk
* @nazwa
* @ossareh
* @paulmach
* @Soulou
* @tike

