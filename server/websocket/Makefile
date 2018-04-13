.PHONY: build doc fmt lint run test vendor_clean vendor_get

GOPATH := ${PWD}/_vendor:${GOPATH}
export GOPATH

default: build

build:
	go build -v -o ./bin/ws ./src/*.go

doc:
	godoc -http=:6060 -index -goroot .

fmt:
	go fmt ./

lint:
	golint .

launch:
	./bin/ws -addr=8081

run: build launch

test:
	@go test ./src/.

vendor_clean:
	rm -rf ./_vendor/src

# We have to set GOPATH to just the _vendor
# directory to ensure that `go get` doesn't
# update packages in our primary GOPATH instead.
# This will happen if you already have the package
# installed in GOPATH since `go get` will use
# that existing location as the destination.

vendor_get: vendor_clean
	GOPATH=${PWD}/_vendor go get -d -u -v \
		github.com/rs/cors \
		github.com/gorilla/mux \
		github.com/gorilla/handlers \
		github.com/gorilla/securecookie \
		golang.org/x/crypto/bcrypt \
		github.com/kylelemons/godebug/pretty \
		github.com/lib/pq \
		github.com/jmoiron/sqlx \
		github.com/dgrijalva/jwt-go \
		github.com/go-redis/redis \
		github.com/mailjet/mailjet-apiv3-go
