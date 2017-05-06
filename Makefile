all:
	go build 
	mv abe bin/
db:
	go build db.go
clean: 
	rm -rf db aba
