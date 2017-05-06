all:
	go build 
db:
	go build db.go
clean: 
	rm -rf db aba
