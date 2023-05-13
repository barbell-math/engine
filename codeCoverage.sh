#!/bin/bash

go test ./... -covermode=count -coverprofile=count.out
go tool cover -func=count.out
go tool cover -html=count.out
