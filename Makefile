VERSION=`cat version.txt`

run:
	go run main.go

release_local: 
	goreleaser build --snapshot --clean

release: 
	echo $(VERSION)
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
	goreleaser release --clean

install:
	go install

clean:
	go clean
	rm -rf dist

fmt:
	go fmt $$(go list ./... | grep -v /vendor/)

test:
	go test ./... -v -coverprofile=c.out
	go tool cover -html="c.out"