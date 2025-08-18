.PHONY: connect serve reset
.DEFAULT_GOAL := serve

reset:
	goose reset

serve:
	MLSSH_HOST=0.0.0.0 go run .

connect:
	ssh -o StrictHostKeyChecking=no -p 23234 mluna@localhost
