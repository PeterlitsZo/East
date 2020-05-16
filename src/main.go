package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "bufio"
    "os"

    "./units"
    "./list"
    "./parse"
)


func getFiles(path string) (files []units.File, files_docid *list.DocList, err error) {
    files_origin, err := units.GetFiles(path)
    files = files_origin
    if err != nil {
        return nil, nil, err
    }
    files_docid = &list.DocList{}
    for _, file := range files {
        files_docid.AddDoc(file.Name)
    }
    return
}


func getWordMap(files []units.File) (wordmap *map[string]*list.DocList) {
    wordmap = new(map[string]*list.DocList)
    for _, file := range files {
        // read file
        file_byte, err := ioutil.ReadFile(file.Path)
        if err != nil {
            fmt.Println("error:", err)
            return
        }
        // split file
        for _, word := range units.Split(string(file_byte)){
            _, ok := (*wordmap)[word]
            if ok {
                // if word in the wordmap then just append it
                doclist := (*wordmap)[word]
                doclist.AddDoc(file.Name)
            } else {
                // if word not in the wordmap then initial a new doclist
                (*wordmap)[word] = &list.DocList{}
                doclist := (*wordmap)[word]
                doclist.AddDoc(file.Name)
            }
        }
    }
    return
}


func getWordMap_fromIndex(file string) (wordmap *map[string]*list.DocList) {
    // read file
    index_byte, err := ioutil.ReadFile(file)
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }
    index_str := string(index_byte)
    // read file's string line to line
    for _, line := range strings.Split(index_str, "\n") {
        // end of file
        if line == "" {
            break
        }
        line_slice := strings.Split(line, "\t")
        // wordmap[name]       = list.DocList[ node, node, ... ]
        (*wordmap)[line_slice[0]] = &list.DocList{}
        for _, docID := range strings.Split(line_slice[2], " ") {
            (*wordmap)[line_slice[0]].AddDoc(docID)
        }
    }
    return
}


func writeWordMap(wordmap *map[string]*list.DocList) {
    index_str := ""
    for key, value := range *wordmap{
        // format: 'key' 'value.length' '*value'
        // for:    'key' 'value.length' ........
        index_str += fmt.Sprintf("%v\t%v\t", key, value.Length)
        // for:    .................... '*value'
        curre := value.Start
        for curre != nil {
            index_str += curre.DocID + " "
            curre = curre.Next
        }
        index_str = index_str[:len(index_str)-1]
        // and a newline token
        index_str += "\n"
    }
    // write into file
    ioutil.WriteFile("index.dict", []byte(index_str), 0644)
    return
}

func main() {
    // ---[ parse the argument ]-----------------------------------------------
    pr := EastArgparse()

    if pr.err != nil {
        fmt.Println(pr.parser.Usage(pr.err))
        return

    } else if pr.version.self.Happened() {
        fmt.Println(units.Version())
        return

    } else if pr.mkindex.self.Happened() {
        files, _, err := getFiles(*pr.mkindex.dirpath)
        if err != nil {
            fmt.Println("[ERROR]:", err)
        }

        wordmap := getWordMap(files)
        writeWordMap(wordmap)
        return

    } else if pr.run.self.Happened() {
        var WordMap *map[string]*list.DocList
        var files_docID *list.DocList
        if *pr.run.useindex {
            WordMap = getWordMap_fromIndex("index.dict")
            _, files_docID, _ = getFiles(*pr.mkindex.dirpath)
        } else {
            var files []units.File
            var err error
            files, files_docID, err = getFiles(*pr.mkindex.dirpath)
            if err != nil {
                fmt.Println("[ERROR]:", err)
            }

            WordMap = getWordMap(files)
        }

        var comast *parse.TypeAst
        comast = parse.GetAST(*pr.run.command)

        result := AST_result(comast.Value.(*parse.TypeList), files_docID, *WordMap)
        fmt.Println("result:", result.Str())

        return

    } else if pr.interactive.self.Happened() {
        var WordMap *map[string]*list.DocList
        var files_docID *list.DocList

        if *pr.run.useindex {
            WordMap = getWordMap_fromIndex("index.dict")
            _, files_docID, _ = getFiles(*pr.mkindex.dirpath)
        } else {
            var files []units.File
            var err error
            files, files_docID, err = getFiles(*pr.mkindex.dirpath)
            if err != nil {
                fmt.Println("[ERROR]:", err)
            }

            WordMap = getWordMap(files)
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
            result := AST_result(comast.Value.(*parse.TypeList), files_docID, *WordMap)
            fmt.Println("result:", result.Str())
        }
        return

    }

    return
}

// copyleft: PeterlitsZo<peterlitszo@outlook.com>

