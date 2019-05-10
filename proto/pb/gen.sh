#!/bin/bash
protoc -I ./ ./*.proto --gofast_out=plugins=grpc:./
