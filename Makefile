leatherman.xz: leatherman
	xz leatherman

leatherman: *.go
	go get -t ./...
	go get github.com/golang/lint/golint
	golint -set_exit_status ./...
	go vet ./...
	TZ=America/Los_Angeles go test
	go build -ldflags "-X main.version=$(TRAVIS_COMMIT)"
	strip leatherman
