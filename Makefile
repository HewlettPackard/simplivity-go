test:
	cd $GOPATH/src/github.com/HewlettPackard/simplivity-go/ovc/ && goveralls  -service=travis-ci
