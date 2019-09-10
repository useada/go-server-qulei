all:
	cd ../../proto/pb/ && ./gen.sh
	@go build -o uauth cmd/*.go
	@echo done.

clean:
	cd ../../proto/pb/ && ./clear.sh
	@go clean
	@rm -f ./uauth
