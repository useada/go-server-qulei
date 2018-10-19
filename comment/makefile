all:
	make -C src
	test -d ./bin/ || mkdir -p ./bin/
	@mv ./src/comment ./bin/
	@echo all done.

clean:
	@rm -f ./bin/comment
	@make -C src clean
