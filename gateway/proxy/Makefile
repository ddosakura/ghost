tls:
	openssl genrsa -out server.key 2048
	openssl req -new -x509 -key server.key -out server.crt -days 365
	# Common Name (e.g. server FQDN or YOUR name) []:localhost

clean:
	rm *.crt *.key
