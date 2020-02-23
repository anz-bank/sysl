@echo off

go test -coverprofile=coverage.txt -covermode=atomic .\...

exit /b %errorlevel%
