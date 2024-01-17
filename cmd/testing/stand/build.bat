set BUILD_DIR=bin
set GOOS=linux
set GOARCH=amd64

REM Build frontend
cd frontend
go build -o ../%BUILD_DIR%/frontend/main .
cd ..

REM Build http_server
cd http_server
go build -o ../%BUILD_DIR%/http_server/main .
cd ..

REM Build processor
cd processor
go build -o ../%BUILD_DIR%/processor/main .
cd ..

REM Build scheduler
cd scheduler
go build -o ../%BUILD_DIR%/scheduler/main .
cd ..

REM Build status
cd status
go build -o ../%BUILD_DIR%/status/main .
cd ..

echo Build completed successfully