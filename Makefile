test:
	cd ./ovc
	go test -v
	goveralls  -service=travis-ci
	
