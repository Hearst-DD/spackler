spackler
========

Spackler enables graceful application termination.  It allows running tasks to complete while preventing new tasks from starting.  Spackler accomplishes this by managing goroutines and exiting timer loops.  This can be of value, for example, in preserving data integrity.

Other features:
* Stop signal available for custom use such as exiting loops
* Custom registration for system signals (`SIGINT`, `SIGIO`, etc.)
* Programmatic termination

[![Build Status](https://travis-ci.org/Hearst-DD/spackler.svg?branch=master)](https://travis-ci.org/Hearst-DD/spackler) [![Coverage](http://gocover.io/_badge/github.com/Hearst-DD/spackler)](http://gocover.io/github.com/Hearst-DD/spackler) [![GoDoc](https://godoc.org/github.com/Hearst-DD/spackler?status.svg)](https://godoc.org/github.com/Hearst-DD/spackler)

Install
=======

```
go get github.com/Hearst-DD/spackler
```

And import:
```go

import "github.com/Hearst-DD/spackler"
```


##### Package Name

This package is named for Carl Spackler, the fictional golf caddy and greenskeeper portrayed by actor Bill Murray in the film, 'Caddyshack', who is entrusted with the graceful termination of gophers.
