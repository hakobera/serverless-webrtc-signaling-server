include .env
export $(shell sed 's/=.*//' .env)

GOBUILD_OPTS = GOOS=linux GOARCH=amd64 GO111MODULE=on

default: build

build:
	$(GOBUILD_OPTS) go build -o bin/onconnect ./onconnect
	$(GOBUILD_OPTS) go build -o bin/ondisconnect ./ondisconnect
	$(GOBUILD_OPTS) go build -o bin/echo ./echo
	$(GOBUILD_OPTS) go build -o bin/register ./register
	$(GOBUILD_OPTS) go build -o bin/broadcast ./broadcast

clean:
	rm -rf bin

test:
	GO111MODULE=on go test ./echo/ -test.v

package: build
	sam package \
	--template-file template.yaml \
	--output-template-file packaged.yaml \
	--s3-bucket ${SAM_BUCKET}

deploy: package
	sam deploy \
	--template-file packaged.yaml \
	--stack-name ${STACK_NAME} \
	--capabilities CAPABILITY_IAM \
	--parameter-overrides StageName=${STAGE} RoomsTableName=${ROOM_TABLE_NAME} ConnectionsTableName=${CONNECTION_TABLE_NAME} \
	--region ${AWS_REGION}

describe-stack:
	@aws cloudformation describe-stacks \
	--stack-name ${STACK_NAME} \
	--query 'Stacks[].Outputs' \
	--region ${AWS_REGION}

delete-stack:
	aws cloudformation delete-stack \
	--stack-name ${STACK_NAME} \
	--region ${AWS_REGION}

.PHONY: clean build test package deploy describe-stack delete-stack
