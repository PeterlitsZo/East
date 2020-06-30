package units

import (
    "fmt"
    "strings"
    "io/ioutil"
)

func _getWordMap_raw(files []File) (wordmap *map[string]*DocList) {
    wordmap_org := make(map[string]*DocList)
    wordmap = &wordmap_org
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
                (*wordmap)[word] = &DocList{}
                doclist := (*wordmap)[word]
                doclist.AddDoc(file.Name)
            }
        }
    }
    return
}


func _getWordMap_fromIndex(file string) (wordmap *map[string]*DocList) {
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
        // wordmap[name]       = DocList[ node, node, ... ]
        (*wordmap)[line_slice[0]] = &DocList{}
        for _, docID := range strings.Split(line_slice[2], " ") {
            (*wordmap)[line_slice[0]].AddDoc(docID)
        }
    }
    return
}


func GetWordMap(useindex bool, dirpath string) (
    WordMap *map[string]*DocList, files_docID *DocList,
) {
	if useindex {
		WordMap = _getWordMap_fromIndex("index.dict")
		_, files_docID, _ = GetFiles(dirpath)
	} else {
		var files []File
		var err error
		files, files_docID, err = GetFiles(dirpath)
		if err != nil {
			fmt.Println("[ERROR]", err)
		}

		WordMap = _getWordMap_raw(files)
	}
    return
}


func WriteWordMap(wordmap *map[string]*DocList) {
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

