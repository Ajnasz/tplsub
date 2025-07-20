.PHONY: run
run:
	@echo '{"FirstName": "John", "LastName": "Doe"}' | go run . example.tmpl --
