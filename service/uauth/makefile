all:
	cd ../../proto/pb/ && ./gen.sh
	make -C src
	test -d ./bin/ || mkdir -p ./bin/
	@mv ./src/uauth ./bin/
	@echo all done.

clean:
	cd ../../proto/pb/ && ./clear.sh
	@rm -f ./bin/uauth
	@make -C src clean
