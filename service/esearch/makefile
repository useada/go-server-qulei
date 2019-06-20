all:
	cd ../../proto/pb/ && ./gen.sh
	make -C src
	test -d ./bin/ || mkdir -p ./bin/
	@mv ./src/esearch ./bin/
	@echo all done.

clean:
	cd ../../proto/pb/ && ./clear.sh
	@rm -f ./bin/esearch
	@make -C src clean
