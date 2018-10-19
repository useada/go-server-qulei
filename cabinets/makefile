all:
	make -C src
	test -d ./bin/ || mkdir -p ./bin/
	@mv ./src/cabinets ./bin/
	@echo all done.

clean:
	@rm -f ./bin/cabinets
	@make -C src clean
