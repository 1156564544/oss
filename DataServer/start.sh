#!/bin/sh

STORAGE_ROOT=/data1 go run server.go :10001 & STORAGE_ROOT=/data2 go run server.go :10002 & STORAGE_ROOT=/data3 go run server.go :10003 & STORAGE_ROOT=/data4 go run server.go :10004 & STORAGE_ROOT=/data5 go run server.go :10005 & STORAGE_ROOT=/data6 go run server.go :10006 & STORAGE_ROOT=/data7 go run server.go :10007 & STORAGE_ROOT=/data8 go run server.go :10008 