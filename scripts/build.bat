@echo off

echo Building server...
go build -o bin/server.exe cmd/server/main.go

echo Building database tools...
go build -o bin/dbtools.exe cmd/tools/main.go

echo Build complete! 