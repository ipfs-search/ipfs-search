#!/bin/sh

BASEDIR=$(dirname $0)

cd $BASEDIR/..

go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
rm coverage.out
