.PHONY: build doc fmt lint run test vendor_clean vendor_get

GOPATH := ${PWD}/_vendor:${GOPATH}
export GOPATH

default: build

fake:
	go run ../fakeData.go

build:
	go build -v -o ./bin/api ./src/*.go

doc:
	godoc -http=:6060 -index -goroot .

fmt:
	go fmt ./

lint:
	golint .

launch:
	./bin/api -port=8080

run: build launch

test:
	@go test ./src/.
	@go test ./src/routes/account/.
	@go test ./src/routes/mail/.
	@go test ./src/routes/profile/.
	@go test ./src/routes/user/.
	@go test ./src/routes/chat/.

test-d:
	@go test ./src/routes/$(TARGET)/.

test-verbose:
	@go test -v ./src/.
	@go test -v ./src/routes/account/.
	@go test -v ./src/routes/mail/.
	@go test -v ./src/routes/profile/.
	@go test -v ./src/routes/user/.
	@go test -v ./src/routes/chat/.

test-target:
	go test -v -run '$(TARGET)' ./src/...

cover:
	go test -cover ./src/...

ROUTES=src/routes/
cover-html:
	go test -coverprofile=cover.out ./src/...
	sed -i -e "s#.*\($(ROUTES).*\.go\)#\./\\1#" cover.out #Format cover.out
	go tool cover -html=cover.out -o coverage.html
	rm cover.out

cover-func:
	@go test -coverprofile=cover.out ./src/...
	@sed -i -e "s#.*\($(ROUTES).*\.go\)#\./\\1#" cover.out #Format cover.out
	@go tool cover -func=cover.out
	@rm cover.out

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
		github.com/mailjet/mailjet-apiv3-go \
		github.com/ip2location/ip2location-go \
		github.com/icrowley/fake
