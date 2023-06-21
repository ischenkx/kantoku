echo Building Test Stand Binaries...
set GOOS=linux
go build -o "client/bin/entry" client/main.go
go build -o "executor/bin/entry" executor/main.go
go build -o "depot/bin/entry" depot/main.go
go build -o "deps/bin/entry" deps/main.go
go build -o "taskdep/bin/entry" taskdep/main.go
go build -o "futdep/bin/entry" futdep/main.go
