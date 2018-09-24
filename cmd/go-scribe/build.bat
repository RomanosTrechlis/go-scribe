echo Getting dependencies
rem go get -v
echo Starting compilation for windows
@echo off
git describe --always --long > version.txt
set /p version= < version.txt
set GOOS=windows
go build -i -v -ldflags "-X 'main.version=%version%'"
go install -v -ldflags "-X 'main.version=%version%'"
@echo on
echo Starting compilation for linux
@echo off
set GOOS=linux
go build -i -v -ldflags "-X 'main.version=%version%'"
set GOOS=windows