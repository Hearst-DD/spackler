// Package spackler enables graceful application termination.  It allows running tasks to complete while preventing new tasks from starting.  Spackler accomplishes this by managing goroutines and exiting timer loops.  This can be of value, for example, in preserving data integrity.
//
//
// Operation
//
// Caddy is composed of reference fields that maintain a common system state, as well as a single value field. The value field indicates whether the current instance was created in a tracked goroutine.  This value determines whether Go() reads stopChan.
//
// The first time Go() or Looper() is called, Spackler starts a goroutine that blocks on sigChan.  sigChan may be registered to receive system signals, such as SIGINT.  When a signal is received, this goroutine closes stopChan, effectively broadcasting "stop!" to any one reading from that channel.
//
// Both Go() and Looper() read from stopChan.  Looper() exits when stopChan is closed. Go() conditionally reads from stopChan if it has been called from an untracked goroutine.  In this case, Go() returns an error if stopChan has been closed.
//
//
// Example
//
// The following example illustrates the main features of this package.
//
//  func doSomething(caddy *spackler.Caddy) {
//    for i := 0; i < 10; i++ {
//
//      //return on stop signal
//      select {
//      case <-caddy.Stopping():
//        return
//      default:
//      }
//
//      // do things
//
//      // this will never fail
//      caddy.Go(func(c *spackler.Caddy) {
//        // do more things
//      })
//    }
//
//  }
//
//  func doSomethingElse() {
//    //do something else
//  }
//
//  func main() {
//
//    caddy := spackler.NewCaddy(true)
//
//    // doSomething() in a tracked goroutine
//    caddy.Looper(1000, true, func() {
//      caddy.Go(doSomething)
//    })
//
//    // this will always fail
//    err := caddy.Go(func(c *spackler.Caddy) {
//      doSomethingElse()
//    })
//
//    caddy.Wait()
//  }
//
// On SIGINT, Looper() will return.  The following call to Go() will fail because main() is not executing in a tracked goroutine and doSomethingElse() will never get called.  Then main() will block on Wait() until all tracked goroutines have completed.  Any running doSomething() threads will exit in the select statement.  If a stop signal is sent when a doSomething() thread is between select and Go(), the call to Go() will still succeed because doSomething() is running in a tracked goroutine.
//
package spackler
