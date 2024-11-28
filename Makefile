vendor:
	cd src/ && go mod vendor

check: clean vendor
	cd src/ && go build

build: check
	sam build --use-container

clean:
	if [ -d "src/vendor" ]; then rm -r src/vendor; fi

deploy:
	sam deploy --guided --capabilities CAPABILITY_NAMED_IAM

delete:
	sam delete --stack-name $(STACK_NAME)

build_and_deploy: clean vendor check build deploy

test:
	cd src/ && go test ./...

coverage:
	cd src/ && go test -coverprofile=coverage.out ./...
	cd src/ && go tool cover -html=coverage.out -o coverage.html