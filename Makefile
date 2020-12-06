GOC = go
MODULE = GoChessgameServer
OUTPUT = chessserver

all: linux windows

linux:
	GOOS=linux $(GOC) build -o ./bin/$(OUTPUT) $(MODULE)
windows:
	GOOS=windows $(GOC) build -o ./bin/$(OUTPUT).exe $(MODULE)
clear:
	rm -r ./bin/*
