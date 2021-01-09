flags = "-X main.COMMIT=`git rev-parse --short HEAD` -X 'main.GOVERSION=`go version`'"
package = buntdb-cli
cov = .coverage
SRC = $(wildcard db/*.go cli/*.go ./*.go)

build : $(package)
test : $(cov)
cov: $(cov)
	go tool cover -html=$(cov)

$(package): $(SRC)
	go build -ldflags $(flags)

.DELETE_ON_ERROR:
$(cov): $(SRC)
	go test -cover -coverprofile=$(cov) ./...


.PHONY : clean
clean :
	-rm -f $(package) $(cov)