build:
	export GOPATH=`pwd`; go install -gcflags "-N -l" main
	./bin/main

win:
	export GOPATH=`pwd` && export CGO_ENABLED=0 && export GOARCH=386 && export GOOS=windows && go build -o ./bin/taoke.exe main
	mkdir taoke
	cp -rf bin taoke
	cp -rf conf taoke
	mkdir taoke/log
	echo "@echo off" > taoke/run.bat
	echo "bin\\\\taoke.exe" >> taoke/run.bat
	zip -r taoke.zip taoke/*
	rm -rf taoke

linux:
	export GOPATH=`pwd` && export CGO_ENABLED=0 && export GOARCH=amd64 && export GOOS=linux && go build -o ./bin/taoke main
	scp ./bin/taoke lizi@10.232.4.31:~
