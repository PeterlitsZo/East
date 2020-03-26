package main

import (
    "fmt"
    "flag"
    "io/ioutil"
    "strconv"
    "strings"
    "os"
)

type file struct {
    path string
    name string
}

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
    // ---[ parse args ]---------------------------------------------------------------------------
    dirpath := flag.String("dirpath", "input", "the input files' path")
    command := flag.String("command", "", "the input files' path")
    mkindex := flag.Bool("mkindex", false, "use this flag to make index named 'index.dict'")
    useindex := flag.Bool("useindex", false, "use file 'index.dict' to find result")
    flag.Parse()

    files, err := getFiles(*dirpath)
    if err != nil { fmt.Println("ERROR:", err); return; }

    com, errmsg := Aim_Parser(*command)
    if errmsg != "" { fmt.Println("ERROR:", errmsg); return; }

    // ---[ read word and make word-map ]----------------------------------------------------------
    WordMap := make(map[string]*DocList)

    // only do not run if under flags: '--useindex ...' without '--mkindex'
    if !*useindex || *mkindex {
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
    }


    // ---[ get the docList for command ]----------------------------------------------------------
    if *mkindex {
        index_str := ""
        for key, value := range WordMap{
            index_str += fmt.Sprintf("%v\t%v\t", key, value.length)
            curre := value.start
            for curre != nil {
                index_str += strconv.Itoa(curre.docID) + " "
                curre = curre.next
            }
            index_str = index_str[:len(index_str)-1]
            index_str += "\n"
        }
        ioutil.WriteFile("index.dict", []byte(index_str), 0644)
    }

    // print the result: mkindex is flase or useindex is true
    // if under '--useindex', then use index.dict to build wordmap
    if *useindex {
        index_byte, err := ioutil.ReadFile("index.dict")
        if err != nil {
            fmt.Println("ERROR:", err)
            return
        }
        index_str := string(index_byte)
        for _, line := range strings.Split(index_str, "\n") {
            // end of file
            if line == "" {
                break
            }
            line_slice := strings.Split(line, "\t")
            // WordMap[name]       = DocList[ node, node, ... ]
            WordMap[line_slice[0]] = &DocList{}
            // fmt.Printf("%#v, %#v\n", line_slice, line)
            for _, docID := range strings.Split(line_slice[2], " ") {
                docid, err := strconv.Atoi(docID)
                if err != nil {
                    fmt.Println("ERROR:", err)
                    return
                }
                WordMap[line_slice[0]].AddDoc(docid)
            }
        }
    }

    // if is without '--mkindex'
    if !*mkindex {
        result := DocList{}
        for _, aim := range com{
            docli, ok := WordMap[aim.value]
            if !ok {
                // can match any thing
                continue
            }
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
        return
    }
}
