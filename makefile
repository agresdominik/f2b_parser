run:
	go run ./src/.

build:
	go build -o bin/parser ./src/.

clean:
	rm bin/*
