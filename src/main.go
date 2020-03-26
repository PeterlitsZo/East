package main

import (
    "fmt"
    "flag"
    "io/ioutil"
    "os"
)

func getFiles(dirpath string) (files []string, err error) {
    dir, err := ioutil.ReadDir(dirpath)
    if err != nil {
        return nil, err
    }
    for _, file := range dir {
        if !file.IsDir() {
            file_path := dirpath + string(os.PathSeparator) + file.Name()
            files = append(files, file_path)
        }
    }
    return files, nil
}

func main() {
    // parse args
    dirpath := flag.String("dirpath", "input", "the input files' path")
    command := flag.String("command", "", "the input files' path")
    flag.Parse()

    files, err := getFiles(*dirpath)
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }
    com, errmsg := Aim_Parser(*command)
    if errmsg != "" {
        fmt.Println("ERROR:", errmsg)
    }

    // read word and make word-map
    WordMap := make(map[string]*DocList)
    for docID, file := range files {
        file_byte, err := ioutil.ReadFile(file)
        if err != nil {
            fmt.Println("ERROR:", err)
            return
        }
        for _, word := range Split(string(file_byte)){
            _, ok := WordMap[word]
            if ok {
                doclist := WordMap[word]
                doclist.AddDoc(docID)
            } else {
                WordMap[word] = &DocList{}
                doclist := WordMap[word]
                doclist.AddDoc(docID)
            }
        }
    }

    // get the docTD
    result := DocList{}
    for _, aim := range com{
        docli := WordMap[aim.value]
        curre := docli.start
        if aim.aim == true{
            for curre != nil {
                result.AddDoc(curre.docID)
                curre = curre.next
            }
        } else {
            for curre != nil {
                result.RemoveDoc(curre.docID)
                curre = curre.next
            }
        }
    }

    fmt.Println("result:", result.Str())
    fmt.Println("-------------------------------------------------------------")

    // print the key and the value
    for key, value := range WordMap {
        fmt.Printf("%#20v (length=%3v) -> %v\n", key, value.length, value.Str())
    }

    return
}
