---
title: Future
---

## What is it?

Future is a wrapper around a dependency and some data storage - in order to resolve a future
you have to put data in it.

```go
type Future struct {
Id string
Dependency string
}

func (fut Future) Resolve(data []byte) error {
if resolved(fut.Id) {
return errors.New("already resolved")
}

return resolve(fut.Id, data)
}
```

Future is like a Promise in javascript - it simplifies asynchronous data processing.