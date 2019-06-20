all:
	cd ../../proto/pb/ && ./gen.sh
	make -C src
	test -d ./bin/ || mkdir -p ./bin/
	@mv ./src/gateway-api ./bin/
	@echo all done.

clean:
	cd ../../proto/pb/ && ./clear.sh
	@rm -f ./bin/gateway-api
	@make -C src clean
