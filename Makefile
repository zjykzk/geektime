all:
	GOBIN=`pwd`/bin go install -v ./apps/*

install: all
	@echo

clean:
	rm -rf bin

.PHONY : all install clean
