#!/bin/sh
export HTTP_PORT=":8080"
export MYSQL_DSN="user:password@/dbname?charset=utf8&parseTime=True&loc=Local"
go run main.go