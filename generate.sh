#!/bin/bash

protoc --go_out=./ --go_opt=paths=source_relative greet/greetpb/greet.proto
protoc --go-grpc_out=./ --go-grpc_opt=paths=source_relative greet/greetpb/greet.proto

protoc --go_out=./ --go_opt=paths=source_relative calculator/calculatorpb/calculator.proto
protoc --go-grpc_out=./ --go-grpc_opt=paths=source_relative calculator/calculatorpb/calculator.proto

protoc --go_out=./ --go_opt=paths=source_relative blog/blogpb/blog.proto
protoc --go-grpc_out=./ --go-grpc_opt=paths=source_relative blog/blogpb/blog.proto

