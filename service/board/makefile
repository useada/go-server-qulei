all:
	cd ../../proto/pb/ && ./gen.sh
	@go build -o board
	@echo done.

clean:
	cd ../../proto/pb/ && ./clear.sh
	@go clean
	@rm -f ./board
