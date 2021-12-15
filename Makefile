.PHONY: clean test

@:
	@make test
	@make clean

clean:
	@cd clientlib && make clean
	@cd shared && make clean
	@cd server && make clean
	@cd client && make clean
	@cd clientwww && make clean

test:
	@cd clientlib && make test
	@cd server && make test