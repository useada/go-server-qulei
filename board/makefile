all:
	make -C src
	test -d ./bin/ || mkdir -p ./bin/
	@mv ./src/board ./bin/
	@echo all done.

clean:
	@rm -f ./bin/board
	@make -C src clean
