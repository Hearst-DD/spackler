spackler
========

Spackler enables graceful application termination.  It allows running tasks to complete while preventing new tasks from starting.  Spackler accomplishes this by managing goroutines and canceling scheduled tasks.  This can be of value, for example, in preserving data integrity.

Other features:
* Stop signal available for custom use such as loop termination
* Custom registration for system signals (`SIGINT`, `SIGIO`, etc.)
* Programmatic termination

[![Build Status](https://travis-ci.org/Hearst-DD/spackler.svg?branch=master)](https://travis-ci.org/Hearst-DD/spackler) [![Coverage](http://gocover.io/_badge/github.com/Hearst-DD/spackler?v=1)](http://gocover.io/github.com/Hearst-DD/spackler?v=1) [![GoDoc](https://godoc.org/github.com/Hearst-DD/spackler?status.svg)](https://godoc.org/github.com/Hearst-DD/spackler)

Install
=======

```
go get github.com/Hearst-DD/spackler
```

And import:
```go

import "github.com/Hearst-DD/spackler"
```

Use
===

#####`Go()`

This function provides goroutine tracking.  It is used in place of the `go` statement.  It calls the provided function passing in a `*Caddy` which may be used to make subsequent calls to `Go()` and `Looper()`.  


#####`Looper()`

This function provides cancelable, scheduled or continuous task execution.  It calls the provided function on the specified interval.  It does not start new goroutines, but the provided function may.


Operation
=========

`Caddy` is composed of reference fields that maintain system state, as well as a single value field. The value field indicates whether the current instance was created in a tracked goroutine.  This value determines whether `Go()` reads `stopChan`.

The first time `Go()` or `Looper()` is called, Spackler starts a goroutine that blocks on `sigChan`.  `sigChan` may be registered to receive system signals, such as `SIGINT`.  When a signal is received, this goroutine closes `stopChan`, effectively broadcasting "stop!" to any one reading from that channel.

Both `Go()` and `Looper()` read from `stopChan`.  `Looper()` exits when `stopChan` is closed. `Go()` conditionally reads from `stopChan` if it has been called from an untracked goroutine.  In this case, `Go()` returns an error if `stopChan` has been closed.


Example
=======

```go
func doSomething(caddy *spackler.Caddy) {
	for i := 0; i < 10; i++ {

		//return on stop signal
		select {
		case <-caddy.Stopping():
			return
		default:
		}

    // do things

    // this will never fail
    caddy.Go(func(c *spackler.Caddy) {
      // do more things
    })		
	}

}

func doSomethingElse() { /* do something else */ }

func main() {

	caddy := spackler.NewCaddy(true)

  // doSomething() in a tracked goroutine
	caddy.Looper(1000, true, func() {
		caddy.Go(doSomething)
	})

  // this will always fail
  err := caddy.Go(func(c *spackler.Caddy) {
		doSomethingElse()
	})

	caddy.Wait()
}
```

On `SIGINT`, `Looper()` will return.  The following call to `Go()` will fail because `main()` is not executing in a tracked goroutine and `doSomethingElse()` will never get called.  Then `main()` will block on `Wait()` until all tracked goroutines have completed.  Any running `doSomething()` threads will exit in the `select` statement.  If a stop signal is sent when a `doSomething()` thread is between `select` and `Go()`, the call to `Go()` will still succeed because `doSomething()` is running in a tracked goroutine.


##### Package Name

This package is named for Carl Spackler, the fictional golf caddy and greenskeeper portrayed by actor Bill Murray in the film, 'Caddyshack', who is entrusted with the graceful termination of gophers.
