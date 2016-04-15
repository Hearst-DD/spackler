package spackler

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const TEST_TIMEOUT = time.Duration(5 * time.Second)

func Test_No_Goroutines(t *testing.T) {
	s := New(false)

	assert.True(t, wait(s))
}

func Test_Stop(t *testing.T) {
	s1 := New(false)
	wg := sync.WaitGroup{}

	wg.Add(1)
	s1.Go(func(s2 *Caddy) {
		defer wg.Done()
	})

	wg.Wait()

	assert.False(t, wait(s1)) // Stop() required

	s1.Stop()

	assert.True(t, wait(s1))
}

func Test_SigChan(t *testing.T) {
	s1 := New(false)
	sigChan := s1.SigChan()

	s1.Go(func(s2 *Caddy) {
		return
	})
	close(sigChan) // same as Stop()

	assert.True(t, wait(s1))
}

func Test_Blocking(t *testing.T) {
	s1 := New(false)
	c1 := make(chan int)
	c2 := make(chan int)

	// create a blocked goroutine
	s1.Go(func(s2 *Caddy) {
		<-c1
	})
	s1.Stop()

	go func() {
		c2 <- 1
		s1.Wait()
		close(c2)
	}()

	<-c2 // ensure the waiting goroutine has started
	select {
	case <-c2:
		assert.True(t, false) // we should be blocked
	default:
		assert.True(t, true)
	}

	c1 <- 1 // unblock the spackler goroutine

	select {
	case <-c2:
		assert.True(t, true)
	case <-time.After(TEST_TIMEOUT):
		assert.True(t, false) // spackler should unblock
	}
}

func Test_Nested_Goroutines(t *testing.T) {
	s1 := New(false)
	x := 0

	s1.Go(func(s2 *Caddy) {
		s2.Go(func(s3 *Caddy) {
			x++
		})

		x++
	})
	s1.Stop()

	assert.True(t, wait(s1))
	assert.True(t, 2 == x)
}

func Test_While_Stopping(t *testing.T) {
	s1 := New(false)
	c := make(chan int)

	s1.Go(func(s2 *Caddy) {
		<-c
		err := s2.Go(func(s3 *Caddy) {
			return
		})
		assert.Nil(t, err)
	})

	s1.Stop()
	c <- 1

	err := s1.Go(func(s2 *Caddy) {
		return
	})
	assert.True(t, nil != err)
}

func Test_Ten_Goroutines(t *testing.T) {
	s1 := New(false)
	x := 0

	for i := 0; i < 10; i++ {
		s1.Go(func(s2 *Caddy) {
			x++
		})
	}
	s1.Stop()

	assert.True(t, wait(s1))
	assert.True(t, 10 == x)
}

func Test_Multiple_Nested_Goroutines(t *testing.T) {
	s1 := New(false)
	x := 0

	for i := 0; i < 10; i++ {
		s1.Go(func(s2 *Caddy) {
			for j := 0; j < 10; j++ {
				s2.Go(func(s3 *Caddy) {
					x++
				})
			}

			x++
		})
	}
	s1.Stop()

	assert.True(t, wait(s1))
	assert.True(t, 110 == x)
}

func Test_Looper_Zero_Duration(t *testing.T) {
	s1 := New(false)
	c := make(chan int)
	x := 0

	s1.Go(func(s2 *Caddy) {
		s2.Looper(0, false, func() {
			c <- 1
			x++
			c <- 1
		})
	})

	<-c // start loop func
	<-c // end loop func

	<-c // start loop func
	<-c // end loop func

	<-c       // start loop func
	s1.Stop() // broadcast quit
	<-c       // end loop func

	assert.True(t, wait(s1))
	assert.True(t, 3 == x)
}

func Test_Looper_NonZero_Duration(t *testing.T) {
	s1 := New(false)
	c := make(chan int)
	x := 0

	s1.Go(func(s2 *Caddy) {
		s2.Looper(1, false, func() {
			c <- 1
			x++
			c <- 1
		})
	})

	<-c // start loop func
	<-c // end loop func

	<-c // start loop func
	<-c // end loop func

	<-c       // start loop func
	s1.Stop() // broadcast quit
	<-c       // end loop func

	assert.True(t, wait(s1))
	assert.True(t, 3 == x)
}

func Test_Looper_RunImmediately(t *testing.T) {
	s1 := New(false)
	c := make(chan int)

	looperTime := time.Second * 3

	s1.Go(func(s2 *Caddy) {
		s2.Looper(looperTime, true, func() {
			close(c)
		})
	})

	select {
	case <-c:
		assert.True(t, true)
	case <-time.After(looperTime - 100):
		assert.True(t, false)
	}

	s1.Stop()
}

func Test_Looper_With_Goroutine(t *testing.T) {
	s1 := New(false)
	c := make(chan int)
	x := 0

	s1.Go(func(s2 *Caddy) {
		s2.Looper(0, false, func() {
			c <- 1
			s2.Go(func(s3 *Caddy) {
				x++
			})
			c <- 1
		})
	})

	<-c // start loop func
	<-c // end loop func

	<-c // start loop func
	<-c // end loop func

	<-c       // start loop func
	s1.Stop() // broadcast quit
	<-c       // end loop func

	assert.True(t, wait(s1))
	assert.True(t, 3 == x)
}

// true if Spackler.Wait() returns in time
func wait(s *Caddy) bool {
	c := make(chan int)

	go func() {
		s.Wait()
		close(c)
	}()

	pass := false
	select {
	case <-c:
		pass = true
	case <-time.After(TEST_TIMEOUT):
	}

	return pass
}
