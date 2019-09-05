all:
	cd ../../proto/pb/ && ./gen.sh
	@go build -o uploader
	@echo done.

clean:
	cd ../../proto/pb/ && ./clear.sh
	@go clean
	@rm -f ./uploader
