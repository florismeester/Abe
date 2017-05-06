all:
	go build 
	mv abe bin/
	sha256sum bin/abe > bin/sha256sum.txt
db:
	go build db.go
clean: 
	rm -rf db aba
