SRC =       ./src/logic.go ./src/main.go
UNITS =     ./src/units/version.go ./src/units/split.go ./src/units/file.go
LIST =      ./src/list/list.go
PARSE =     ./src/parse/parse.y ./src/parse/parse.go
ARGPARSE =  ./src/argparse/argparse.go

main: $(SRC) $(UNITS) $(LIST) $(PARSE) $(ARGPARSE)
	go build -o main ./src

./src/parse/parse.go : ./src/parse/parse.y
	goyacc -o ./src/parse/parse.go ./src/parse/parse.y
	rm y.output

.PHONY: init
init:
	go get -u -v github.com/PeterlitsZo/argparse

.PHONY: clean
clean:
	rm ./main
	rm ./index.dict
	rm ./src/parse.go
