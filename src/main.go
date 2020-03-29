package main

import (
    "fmt"
    "flag"
    "io/ioutil"
    "strings"
    "os"
)

type File struct {
    path string
    name string
}

func getFiles(dirpath string) (files []File, err error) {
    dir, err := ioutil.ReadDir(dirpath)
    if err != nil {
        return nil, err
    }
    for _, file := range dir {
        if !file.IsDir() {
            file_path := dirpath + string(os.PathSeparator) + file.Name()
            files = append(files, File{path: file_path, name: file.Name()})
        }
    }
    return files, nil
}

func main() {
    // ---[ parse args ]---------------------------------------------------------------------------
    dirpath := flag.String("dirpath", "input", "the input files' path")
    command := flag.String("command", "", "the command to get the ID list (see README.pdf)")
    mkindex := flag.Bool("mkindex", false, "use this flag to make index named 'index.dict'")
    useindex := flag.Bool("useindex", false, "use file 'index.dict' to find result")
    flag.Parse()

    files, err := getFiles(*dirpath)
    files_docID := []string{}
    for _, file := range files {
        files_docID = append(files_docID, file.name)
    }
    if err != nil { fmt.Println("ERROR:", err); return; }

    comast := getAST(*command)

    // ---[ read word and make word-map ]----------------------------------------------------------
    WordMap := make(map[string]*DocList)

    // only do not run if under flags: '--useindex ...' without '--mkindex'
    if !*useindex || *mkindex {
        for _, file := range files {
            file_byte, err := ioutil.ReadFile(file.path)
            if err != nil {
                fmt.Println("ERROR:", err)
                return
            }
            for _, word := range Split(string(file_byte)){
                _, ok := WordMap[word]
                if ok {
                    doclist := WordMap[word]
                    doclist.AddDoc(file.name)
                } else {
                    WordMap[word] = &DocList{}
                    doclist := WordMap[word]
                    doclist.AddDoc(file.name)
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
                index_str += curre.docID + " "
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
                WordMap[line_slice[0]].AddDoc(docID)
            }
        }
    }

    // if is without '--mkindex' or with '--useindex'
    if !*mkindex || *useindex {
        result := DocList{}
        if comast == nil {
            fmt.Println("Sorry but I need your command valid")
            return
        }
        // first loop: get the expr
        for _, expr := range *comast {
            expr_result := DocList{}
            for _, file_docID := range files_docID{
                expr_result.AddDoc(file_docID)
            }

            // second loop: get the atom
            for _, atom := range *expr{
                doclist, ok := WordMap[atom.str]
                if !ok {
                    doclist = &DocList{}
                }

                if atom.not {
                    current := doclist.start
                    for current != nil {
                        expr_result.RemoveDoc(current.docID)
                        current = current.next
                    }
                } else {
                    current := expr_result.start
                    for current != nil {
                        if !doclist.Has(current.docID) {
                            expr_result.RemoveDoc(current.docID)
                        }
                        current = current.next
                    }
                }
            }
            // end of second loop 

            current := expr_result.start
            for current != nil {
                result.AddDoc(current.docID)
                current = current.next
            }
        }
        fmt.Println("result:", result.Str())
        return
    }
    return
}
