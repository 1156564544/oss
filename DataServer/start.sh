#!/bin/sh

STORAGE_ROOT=/data1 go run server.go :10001 & STORAGE_ROOT=/data2 go run server.go :10002 & STORAGE_ROOT=/data3 go run server.go :10003 & STORAGE_ROOT=/data4 go run server.go :10004