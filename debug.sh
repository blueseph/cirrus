#!/bin/bash

$GOPATH/bin/dlv --headless debug --api-version 2 --listen 127.0.0.1:2345 -- up --stack cirrus