SRC = ./src/argparse.go ./src/file.go ./src/list.go ./src/logic.go \
	  ./src/main.go ./src/parse.go
UNITS = ./src/units/version.go ./src/units/split.go

main: $(SRC) $(UNITS)
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
