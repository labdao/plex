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
```
