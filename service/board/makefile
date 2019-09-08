all:
	cd ../../proto/pb/ && ./gen.sh
	@go build -o board cmd/*.go
	@echo done.

clean:
	cd ../../proto/pb/ && ./clear.sh
	@go clean
	@rm -f ./board
