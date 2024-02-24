VERSION=`cat version.txt`

run:
	go run main.go

release_local: 
	goreleaser release --snapshot --rm-dist

release: 
	echo $(VERSION)
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
	goreleaser release --rm-dist

install:
	go install

clean:
	go clean
	rm -rf dist

fmt:
	go fmt $$(go list ./... | grep -v /vendor/)