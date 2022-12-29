GO_TEST_FLAGS = -v -c -coverpkg ./...

.PHONY: run test

run:
	go run ./cmd/bot

test:
	mkdir -p ./test/datasource
	go test ${GO_TEST_FLAGS} -o ./test/datasource/compiled ./pkg/datasource

test_db: test
	./test/datasource/compiled -test.run ^TestDB -test.v -test.count=1 -test.coverprofile ./test/datasource/db.coverage
