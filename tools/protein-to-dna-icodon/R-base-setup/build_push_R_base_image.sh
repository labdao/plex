#!/bin/bash

docker buildx build --platform linux/amd64 -t docker.io/supraja968/my_r_base_image:latest . --push