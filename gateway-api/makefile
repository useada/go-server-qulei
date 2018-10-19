all:
	make -C src
	test -d ./bin/ || mkdir -p ./bin/
	@mv ./src/gateway-api ./bin/
	@echo all done.

clean:
	@rm -f ./bin/gateway-api
	@make -C src clean
