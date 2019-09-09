all:
	cd ../../proto/pb/ && ./gen.sh
	@go build -o esearch cmd/*.go
	@echo done.


clean:
	cd ../../proto/pb/ && ./clear.sh
	@go clean
	@rm -f ./esearch
