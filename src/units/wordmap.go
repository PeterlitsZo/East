package units

import (
    "fmt"
    "strings"
    "io/ioutil"

    "../list"
)

func GetWordMap(files []File) (wordmap *map[string]*list.DocList) {
    wordmap = new(map[string]*list.DocList)
    for _, file := range files {
        // read file
        file_byte, err := ioutil.ReadFile(file.Path)
        if err != nil {
            fmt.Println("error:", err)
            return
        }
        // split file
        for _, word := range Split(string(file_byte)){
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


func GetWordMap_fromIndex(file string) (wordmap *map[string]*list.DocList) {
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
        line_slice := strings.Split(line, "	")
        // wordmap[name]       = list.DocList[ node, node, ... ]
        (*wordmap)[line_slice[0]] = &list.DocList{}
        for _, docID := range strings.Split(line_slice[2], " ") {
            (*wordmap)[line_slice[0]].AddDoc(docID)
        }
    }
    return
}


func WriteWordMap(wordmap *map[string]*list.DocList) {
    index_str := ""
    for key, value := range *wordmap{
        // format: 'key' 'value.length' '*value'
        // for:    'key' 'value.length' ........
        index_str += fmt.Sprintf("%v	%v	", key, value.Length)
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

