#!/bin/bash

function_name="linebotExample"

GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go
zip myFunction.zip bootstrap

if aws lambda list-functions | grep -q $function_name; then
  aws lambda update-function-code \
    --function-name $function_name \
    --zip-file fileb://myFunction.zip
else
  aws lambda create-function \
    --function-name $function_name \
    --runtime provided.al2023 \
    --role arn:aws:iam::828951707561:role/awslambdaBasicExecutionRole \
    --environment Variables="{LINE_CHANNEL_SECRET=YOUR_CHANNEL_SECRET,LINE_CHANNEL_ACCESS_TOKEN=YOUR_CHANNEL_ACCESS_TOKEN}" \
    --handler bootstrap \
    --zip-file fileb://myFunction.zip
fi
