.PHONY: fake vendor_clean vendor_get

GOPATH := ${PWD}/_vendor:${GOPATH}
export GOPATH

fake:
	psql -f ../setup_db.sql > /dev/null
	go run fakeData.go

vendor_clean:
	rm -rf ./_vendor/src

vendor_get: vendor_clean
	GOPATH=${PWD}/_vendor go get -d -u -v \
		github.com/icrowley/fake
