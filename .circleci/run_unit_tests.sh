#/bin/bash

printenv

go test ./... -coverprofile=coverage.out
ls | grep *.out

if [[ "$CIRCLE_BRANCH" == "master" ]]; then
  curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
  chmod +x ./cc-test-reporter
  ./cc-test-reporter format-coverage ./coverage.out -t gocov
  ./cc-test-reporter upload-coverage
fi
