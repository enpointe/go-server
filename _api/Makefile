
all:
	echo "Generating index.html"
	openapi-generator generate -i openapi.yaml -g html
	echo "Copying generated to server"
	cp index.html ../cmd/server/public

clean:
	rm -f index.html
