all:
	GOBIN=`pwd`/bin go install -v ./apps/*

install: all
	@echo

clean:
	rm -rf bin

test:
	go test . -v

.PHONY : all install clean
