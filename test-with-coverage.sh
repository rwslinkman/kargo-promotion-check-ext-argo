#!/bin/sh
rm coverage.out
rm coverage.html

go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
echo "HTML coverage report generated"
echo "Open at file://$(pwd)/coverage.html"