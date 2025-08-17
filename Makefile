.PHONY: connect serve
.DEFAULT_GOAL := serve


serve:
	go run .

connect:
	ssh -o StrictHostKeyChecking=no -p 23234 mluna@localhost
