spackler
========

Spackler provides mechanisms for the graceful termination of an application which.

A mechanism is provided for tracking goroutines.  After a stop signal has been
received, untracked goroutines are prevented from starting new goroutines, but
tracked goroutines are allowed to start new tracked goroutines in order to complete
their tasks.

A mechanism is provided to execute recurring tasks on a fixed interval.  This
mechanism exits after a stop signal has been received.

Other features:
* Stop signal available for custom use such as the termination of for loops
* Custom registration for system signals (SIGINT, SIGIO, etc.)
* Programmatic termination

[![Build Status](https://travis-ci.org/Hearst-DD/spackler.svg?branch=master)](https://travis-ci.org/Hearst-DD/spackler) [![Coverage](http://gocover.io/_badge/github.com/Hearst-DD/spackler?v=1)](http://gocover.io/github.com/Hearst-DD/spackler?v=1) [![GoDoc](https://godoc.org/github.com/Hearst-DD/spackler?status.svg)](https://godoc.org/github.com/Hearst-DD/spackler)

Install
=======

```
go get github.com/Hearst-DD/spackler
```

And import:
```go

import l5g "github.com/Hearst-DD/spackler"
```


Operation
=========

When Go() or Looper() are called, Spackler starts a goroutine that blocks on
sigChan, which may be registered to receive system signals, such as SIGINT.  When
a signal is received, the goroutine closes stopChan, effectively broadcasting
"stop!" to any one reading from that channel.

Both Go() and Looper() read from stopChan.  Looper() exits when it reads from stopChan.
Go() conditionally reads from stopChan if it has been called from an untracked
goroutine.  If it reads from stopChan, it returns an error rather than creating
a new tracked goroutine.

Caddy is composed of reference fields that maintain system state, as well as a single
value field which indicates whether the current instance was created in a tracked
goroutine.  This value determines whether Go() should read stopChan.


Example
=======

```go
func doSomething(caddy *spackler.Caddy) {
	for i := 0; i < 10; i++{

		//return on stop signal
		select {
		case <-caddy.Stopping():
			return
		default:
		}

		// do something
	}

}

func doSomethingElse() { // do something else }

func main() {

	caddy := spackler.NewCaddy(true)

	caddy.Looper(1000, true, func() {
		caddy.Go(doSomething)
	})

	// you can wrap functions that don't need a Caddy in annonymous functions
	err := caddy.Go(func(c *spackler.Caddy) {
		doSomethingElse()
	})

	caddy.Wait()
}
```

On SIGINT, Looper will return.  The next instruction in main(), a call to Go(), will
fail and main() will block on Wait() until all running instances of doSomething()
have completed.  doSomething() may be interrupted by the stop signal at the top
of the for loop.


# Package Name

This package is named for Carl Spackler, the fictional golf caddy and greenskeeper
portrayed by actor Bill Murray in the film, 'Caddyshack', who is entrusted with
the graceful termination of gophers.
