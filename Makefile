test:
	go test -v ./ovc/...
	cd ./ovc && goveralls  -service=travis-ci
	
