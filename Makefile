%PHONEY: connect


connect:
	ssh -o StrictHostKeyChecking=no -p 23234 mluna@localhost
