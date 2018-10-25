all:
	make -C src
	test -d ./bin/ || mkdir -p ./bin/
	@mv ./src/msgboard ./bin/
	@echo all done.

clean:
	@rm -f ./bin/msgboard
	@make -C src clean
