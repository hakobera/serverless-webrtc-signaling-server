[![CircleCI](https://circleci.com/gh/hakobera/serverless-webrtc-signaling-server.svg?style=svg)](https://circleci.com/gh/hakobera/serverless-webrtc-signaling-server)

# serverless-webrtc-signaling-server

This is the code and template for the serverless-webrtc-signaling-server. There are three functions contained within the directories and a SAM teamplte that wires them up to a DynamoDB table and provides the minimal set of permissions needed to run the app.

## Requirements

* AWS CLI already configured with Administrator permission
* [Golang 1.12 or greater](https://golang.org)

## Setup process

### Installing dependencies

This repository use [GO Modules](https://github.com/golang/go/wiki/Modules) to manage dependencies.
Dependencies are automaticaly download when you build code by `make build` command.

### Building

Golang is a statically compiled language, meaning that in order to run it you have to build the executable target.

You can issue the following command in a shell to build it:

```shell
make build
```

### Local testing

**Unfortunatelly sam CLI does not support WebSocket local test. So you need to deploy to test functionality.**

## Packaging and deployment

First and foremost, we need a `S3 bucket` where we can upload our Lambda functions packaged as ZIP before we deploy anything - If you don't have a S3 bucket to store code artifacts then this is a good time to create one:

```bash
aws s3 mb s3://BUCKET_NAME
```

Next, create file named `.env` using `cp .env.template .env` command.
Then open the file in your editer and set S3 bucket name which you created to `SAM_BUCKET`.

```bash
# Must need to change
SAM_BUCKET=REPLACE_THIS_WITH_YOUR_S3_BUCKET_NAME

# Change if you want 
AWS_REGION=ap-northeast-1
STAGE=Dev
STACK_NAME=webrtc-signaling-server
ROOM_TABLE_NAME=Rooms
CONNECTION_TABLE_NAME=Connections
```

At last, run the following command to package our Lambda function to S3:

```shell
make deploy
```

> **See [Serverless Application Model (SAM) HOWTO Guide](https://github.com/awslabs/serverless-application-model/blob/master/HOWTO.md) for more details in how to get started.**

After deployment is complete you can run the following command to retrieve the API Gateway Endpoint URL:

```bash
make describe-stack
``` 
## Connect to API Gateway Endpoint URL using wscat (WebSocket CLI client)

```bash
$ npm instal -g wscat
$ wscat $(make describe-stack | jq -r '.[][] | select(.OutputKey == "WebSocketURI") | .OutputValue')
```

### Unit Testing

We use `testing` package that is built-in in Golang and you can simply run the following command to run our tests:

```shell
make test
```

## License

Apache License 2.0
