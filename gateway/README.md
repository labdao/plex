```
pip install localstack

pip install awscli-local

localstack start -d

localstack status services

awslocal sqs create-queue --queue-name job-queue


```

Utility

```
docker ps | grep localstack
docker restart <id_from_prev_step>

zip function.zip main

aws --endpoint-url=http://localhost:4566 lambda create-function \
    --function-name gateway \
    --handler main \
    --zip-file fileb://./function.zip \
    --runtime go1.x \
    --role arn:aws:iam::123456789012:role/dummyrole

```
