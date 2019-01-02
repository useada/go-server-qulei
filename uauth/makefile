all:
	make -C src
	test -d ./bin/ || mkdir -p ./bin/
	@mv ./src/uauth ./bin/
	@echo all done.

clean:
	@rm -f ./bin/uauth
	@make -C src clean
