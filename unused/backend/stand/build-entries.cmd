echo Hello World
set GOOS=linux
go build -o "client/bin/entry" client/main.go
go build -o "executor/bin/entry" executor/main.go
go build -o "depot/bin/entry" depot/main.go
go build -o "deps/bin/entry" deps/main.go
