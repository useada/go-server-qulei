all:
	cd ../../proto/pb/ && ./gen.sh
	make -C src
	test -d ./bin/ || mkdir -p ./bin/
	@mv ./src/uploader ./bin/
	@echo all done.

clean:
	cd ../../proto/pb/ && ./clear.sh
	@rm -f ./bin/uploader
	@make -C src clean
