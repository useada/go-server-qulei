all:
	make -C src
	test -d ./bin/ || mkdir -p ./bin/
	@mv ./src/uploader ./bin/
	@echo all done.

clean:
	@rm -f ./bin/uploader
	@make -C src clean
