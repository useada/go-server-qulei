all:
	make -C src
	test -d ./bin/ || mkdir -p ./bin/
	@mv ./src/esearch ./bin/
	@echo all done.

clean:
	@rm -f ./bin/esearch
	@make -C src clean
