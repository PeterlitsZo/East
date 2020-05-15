main: ./src/parse.go ./src/list.go ./src/main.go ./src/split.go ./src/argparse.go
	go build -o main ./src

./src/parse.go : ./src/parse.y
	goyacc -o ./src/parse.go ./src/parse.y
	rm y.output

.PHONY: init
init:
	go get -u -v github.com/PeterlitsZo/argparse

.PHONY: clean
clean:
	rm ./main
	rm ./index.dict
	rm ./src/parse.go
