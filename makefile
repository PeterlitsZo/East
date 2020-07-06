viewer     = chrome.exe

SRC_F      = ./src
SRC        = $(SRC_F)/main.go
UNITS_F    = $(SRC_F)/units
UNITS      = $(UNITS_F)/version.go      $(UNITS_F)/file.go     $(UNITS_F)/list.go
PARSE_F    = $(SRC_F)/parse
PARSE      = $(PARSE_F)/parse.y         $(PARSE_F)/parse.go
ARGPARSE_F = $(SRC_F)/argparse
ARGPARSE   = $(ARGPARSE_F)/argparse.go
LOGIC_F    = $(SRC_F)/logic
LOGIC      = $(LOGIC_F)/main.go    		$(LOGIC_F)/sreach.go
INDEX_F    = $(SRC_F)/index
INDEX      = $(INDEX_F)/split.go

DOC        = ./doc/README.tex

main: $(SRC) $(UNITS) $(PARSE) $(ARGPARSE) $(LOGIC)
	go build -o main ./src

./src/parse/parse.go : ./src/parse/parse.y
	goyacc -o ./src/parse/parse.go ./src/parse/parse.y
	rm y.output

.PHONY: run
run: $(PARSE)
	@echo "------------------------------------------------"
	@go run ./src interactive

.PHONY: doc
doc: $(DOC)
	cd ./doc && lualatex ./README.tex && mv ./README.pdf ..

.PHONY: look
look:
	$(viewer) ./README.pdf

.PHONY: init
init:
	go get -u -v github.com/PeterlitsZo/argparse

.PHONY: clean
clean:
	-rm ./main
	-rm ./index.dict
	-rm ./src/parse.go

.PHONY: test
test:
	go test $(INDEX_F)
