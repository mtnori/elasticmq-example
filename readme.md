# Memo

## put message
```
aws sqs send-message \
--endpoint-url "http://localhost:9324" \
--queue-url "http://localhost:9324/000000000000/sample1" \
--message-body "SAMPLE"
```

## get approximate number of messages
```
aws sqs get-queue-attributes \
--endpoint-url "http://localhost:9324" \
--queue-url "http://localhost:9324/000000000000/sample1" \
--attribute-names ApproximateNumberOfMessages
```