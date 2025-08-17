.PHONY: connect serve
.DEFAULT_GOAL := serve


serve:
	MLSSH_HOST=0.0.0.0 go run .

connect:
	ssh -o StrictHostKeyChecking=no -p 23234 mluna@localhost
