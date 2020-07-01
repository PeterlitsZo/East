package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"

    "./units"
    "./parse"
    "./argparse"
    "./logic"
)


func main() {
    // parse the argument
    pr := argparse.EastArgparse()

    // if East raise the error, then putout the usage
    if pr.Err != nil {
        fmt.Println(pr.Parser.Usage(pr.Err))
        return

    // if it need the version, then putout the version
    } else if pr.Version.Self.Happened() {
        fmt.Println(units.Version())
        return

    // if it need to make the index file, then make it
    } else if pr.Mkindex.Self.Happened() {
        wordmap, _ := units.GetWordMap(false, *pr.Mkindex.Dirpath)
        units.WriteWordMap(wordmap)
        return

    // if it need to run a command, then run it
    } else if pr.Run.Self.Happened() {
        WordMap, files_docID := units.GetWordMap(
            // need index file?
            *pr.Run.Useindex,
            // the input folder path
            *pr.Run.Dirpath,
        )

        comast := parse.GetAST(*pr.Run.Command)
        result, runresult := logic.Run(comast, &units.Env{files_docID, WordMap})
        // output the return
        if !runresult.NoOutput {
            fmt.Println("Result  :", result, "\n")
        }

        return

    // if it need to get into interactive mode, then make it
    } else if pr.Interactive.Self.Happened() {
        WordMap, files_docID := units.GetWordMap(
            // need index file?
            *pr.Interactive.Useindex,
            // the input folder path
            *pr.Interactive.Dirpath,
        )

        var comast *parse.AST

        fmt.Println("Enter `quit` for quit")
        fmt.Println("copyleft (C) Peterlits Zo <peterlitszo@outlook.com>")
        fmt.Println("Github: github.com/PeterlitsZo/East, version:", units.Version())
        fmt.Println("")
        for {
            reader := bufio.NewReader(os.Stdin)
            fmt.Print("Command > ")
            text, _ := reader.ReadString('\n')
            text = strings.Replace(text, "\n", "", -1)
            // parse the input as a AST
            comast = parse.GetAST(text)
            // use logicer to hold it.
            result, runresult := logic.Run(comast, &units.Env{files_docID, WordMap})
            // output the return
            if !runresult.NoOutput {
                fmt.Println("Result  :", result, "\n")
            }
            if runresult.NeedBreak {
                break
            }
        }

    }

    return

}

// copyleft: PeterlitsZo<peterlitszo@outlook.com>

