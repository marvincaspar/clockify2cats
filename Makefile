VERSION=`cat version.txt`

run:
	go run main.go

release_local: 
	goreleaser build --snapshot --clean

release: 
	echo $(VERSION)
	git tag -s -a $(VERSION) -m "Release $(VERSION)"
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
	go test ./... -v -coverprofile=covprofile
	
	# ignore files from coverage, e.g. repository to fetch remote data
	grep -v -E -f .covignore covprofile > covprofile.filtered
	mv covprofile.filtered covprofile

	go tool cover -html="covprofile"