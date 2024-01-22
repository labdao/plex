#!/bin/bash

docker buildx build --platform linux/amd64 -t quay.io/labdao/r-base-with-icodon:latest . --push