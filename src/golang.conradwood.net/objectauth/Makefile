.PHONY:	client server
all:	client server
client:
	cd client && go install objectauth-client.go
server:
	cd server && go install objectauth-server.go `ls -1 *.go|grep -v objectauth-server.go`
