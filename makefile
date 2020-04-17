main: ./src/parse.go ./src/list.go ./src/main.go ./src/split.go
	go get -u -v github.com/akamensky/argparse
	go build -o main ./src

./src/parse.go :
	goyacc -o ./src/parse.go ./src/parse.y
	rm y.output

.PHONY: clean
clean:
	rm ./src/parse.go
