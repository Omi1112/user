#/bin/sh
golint ./... | tee .golint.txt
test ! -s .golint.txt
