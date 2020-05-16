package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"

    "./units"
    "./list"
    "./parse"
    "./argparse"
    "./logic"
)


func main() {
    pr := argparse.EastArgparse()

    if pr.Err != nil {
        fmt.Println(pr.Parser.Usage(pr.Err))
        return

    } else if pr.Version.Self.Happened() {
        fmt.Println(units.Version())
        return

    } else if pr.Mkindex.Self.Happened() {
        files, _, err := units.GetFiles(*pr.Mkindex.Dirpath)
        if err != nil {
            fmt.Println("[ERROR]:", err)
        }

        wordmap := units.GetWordMap(files)
        units.WriteWordMap(wordmap)
        return

    } else if pr.Run.Self.Happened() {
        var WordMap *map[string]*list.DocList
        var files_docID *list.DocList
        if *pr.Run.Useindex {
            WordMap = units.GetWordMap_fromIndex("index.dict")
            _, files_docID, _ = units.GetFiles(*pr.Mkindex.Dirpath)
        } else {
            var files []units.File
            var err error
            files, files_docID, err = units.GetFiles(*pr.Mkindex.Dirpath)
            if err != nil {
                fmt.Println("[ERROR]:", err)
            }

            WordMap = units.GetWordMap(files)
        }

        var comast *parse.TypeAst
        comast = parse.GetAST(*pr.Run.Command)

        result := logic.AST_result(comast.Value.(*parse.TypeList), files_docID, *WordMap)
        fmt.Println("result:", result.Str())

        return

    } else if pr.Interactive.Self.Happened() {
        var WordMap *map[string]*list.DocList
        var files_docID *list.DocList

        if *pr.Run.Useindex {
            WordMap = units.GetWordMap_fromIndex("index.dict")
            _, files_docID, _ = units.GetFiles(*pr.Mkindex.Dirpath)
        } else {
            var files []units.File
            var err error
            files, files_docID, err = units.GetFiles(*pr.Mkindex.Dirpath)
            if err != nil {
                fmt.Println("[ERROR]:", err)
            }

            WordMap = units.GetWordMap(files)
        }

        var comast *parse.TypeAst
        fmt.Println("Enter `quit` for quit")
        fmt.Println("copyleft (C) Peterlits Zo <peterlitszo@outlook.com>")
        fmt.Println("Github: github.com/PeterlitsZo/East, version:", units.Version())
        fmt.Println("")
        for {
            reader := bufio.NewReader(os.Stdin)
            fmt.Print("Command > ")
            text, _ := reader.ReadString('\n')
            text = strings.Replace(text, "\n", "", -1)
            if text == "quit" {
                return
            }
            comast = parse.GetAST(text)
            result := logic.AST_result(comast.Value.(*parse.TypeList), files_docID, *WordMap)
            fmt.Println("result:", result.Str())
        }
        return

    }

    return
}

// copyleft: PeterlitsZo<peterlitszo@outlook.com>

