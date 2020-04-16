package main

import (
    "fmt"
    "flag"
    "io/ioutil"
    "strings"
    "reflect"
    "os"
)

// ----------------------------------------------------------------------------
// ---[ UNITS FUNCTION ]-------------------------------------------------------
// ----------------------------------------------------------------------------

// ---[ return all files ]-----------------------------------------------------
// file struct
type File struct {
    path string
    name string
}

// getFiles: funtion that need a string then return all Files under that folder
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

// ---[ return bool sreach's result ]------------------------------------------
// return the AST's result.
func AST_result(ast *typeAst, all_docID *DocList, wordmap map[string]*DocList) *DocList {
    result := &DocList{}
    // if the ast's len is zero, then return the full DocList
    if ast == nil || len(*ast) == 0 {
        result.Copy(all_docID)
        return result
    }
    for _, expr_ptr := range *ast {
        sub_result := EXPR_result(expr_ptr, all_docID, wordmap)
        current := sub_result.start
        for current != nil {
            // get the union of all sub_result
            result.AddDoc(current.docID)
            current = current.next
        }
    }
    return result
}

// return the result of expr:
// --------------------------
// return the expr's result. all atom is link by op 'and'
func EXPR_result(expr *typeExpr, all_docID *DocList, wordmap map[string]*DocList) *DocList {
    // if the expr's len is zero, then return empty DocList
    result := &DocList{}
    if len(*expr) == 0 {
        return result
    // else initial it by the first node
    } else {
        result.Copy(ATOM_result((*expr)[0], all_docID, wordmap))
    }
    for _, atom_ptr := range *expr {
        sub_result := ATOM_result(atom_ptr, all_docID, wordmap)
        current := result.start
        // iter the sub-result
        for current != nil {
            // only need the node that the sub-result has
            if !sub_result.Has(current.docID) {
                result.RemoveDoc(current.docID)
            }
            current = current.next
        }
    }
    return result
}

// return the result of atom:
// --------------------------
// atom's value is a interface, we need know that the type of the value. if the
// type is typeAst, then it need call AST, else it should be a string, so we
// need get the result by wordmap. if atom.not, then need to negate it by full
// docID list.
func ATOM_result(atom *typeAtom, all_docID *DocList, wordmap map[string]*DocList) *DocList {
    result := &DocList{}
    switch v := atom.value.(type) {
    case string:
        doclist_ptr, ok := wordmap[atom.value.(string)]
        // if it have the key then copy else return empty list, else copy the 
        if ok {
            result.Copy(doclist_ptr)
        }
    case *typeAst:
        sub_ast := atom.value.(*typeAst)
        result.Copy(AST_result(sub_ast, all_docID, wordmap))
    default:
        // TODO: this is not OK( it look ugly )
        fmt.Println("Error! the type is ", reflect.TypeOf(v))
        return result
    }
    // negate the result
    if atom.not {
        current := all_docID.start
        for current != nil {
            // does not have? make it have
            if !result.Has(current.docID) {
                result.AddDoc(current.docID)
            // does it have? remove it
            } else {
                result.RemoveDoc(current.docID)
            }
            current = current.next
        }
    }
    return result
}

// ----------------------------------------------------------------------------
// ---[ MAIN FUNCTION ]--------------------------------------------------------
// ----------------------------------------------------------------------------

func main() {
    // ---[ parse args ]-------------------------------------------------------

    // initial the arg-parser
    dirpath := flag.String("dirpath", "input", "the input files' path")
    command := flag.String("command", "",
                           "the command to get the ID list (see README.pdf)")
    mkindex := flag.Bool("mkindex", false,
                         "use this flag to make index named 'index.dict'")
    useindex := flag.Bool("useindex", false,
                          "use file 'index.dict' to find result")
    flag.Parse()

    // ---[ initial all variable ]---------------------------------------------

    // get all Files under the given folder
    files, err := getFiles(*dirpath)
    // files_docID is docID of all Files
    files_docID := &DocList{}
    for _, file := range files {
        files_docID.AddDoc(file.name)
    }
    if err != nil { fmt.Println("ERROR:", err); return; }

    // parser the command ast
    comast := getAST(*command)

    // ---[ read word and make word-map ]--------------------------------------
    WordMap := make(map[string]*DocList)
    // TODO: will it a good idea to make it be a higher function to help wordmap
    //       add a file's name?

    // only do *NOT* run if under flags: '--useindex ...' without '--mkindex'
    // : make a index(WordMap) by the files under the floder
    if !*useindex || *mkindex {
        for _, file := range files {
            // read file
            file_byte, err := ioutil.ReadFile(file.path)
            if err != nil {
                fmt.Println("ERROR:", err)
                return
            }
            // split file
            for _, word := range Split(string(file_byte)){
                _, ok := WordMap[word]
                if ok {
                    // if word in the WordMap then just append it
                    doclist := WordMap[word]
                    doclist.AddDoc(file.name)
                } else {
                    // if word not in the WordMap then initial a new DocList
                    WordMap[word] = &DocList{}
                    doclist := WordMap[word]
                    doclist.AddDoc(file.name)
                }
            }
        }
    }


    // ---[ get the docList for command ]--------------------------------------
    // if use '--mkindex'
    // : it will turn map WordMap into a file.
    if *mkindex {
        index_str := ""
        for key, value := range WordMap{
            // format: 'key' 'value.length' '*value'
            // for:    'key' 'value.length' ........
            index_str += fmt.Sprintf("%v\t%v\t", key, value.length)
            // for:    .................... '*value'
            curre := value.start
            for curre != nil {
                index_str += curre.docID + " "
                curre = curre.next
            }
            index_str = index_str[:len(index_str)-1]
            // and a newline token
            index_str += "\n"
        }
        // write into file
        ioutil.WriteFile("index.dict", []byte(index_str), 0644)
    }

    // if under '--useindex'
    // then use index.dict to build wordmap print the result: mkindex is flase
    // or useindex is true
    if *useindex && !*mkindex {
        // read file
        index_byte, err := ioutil.ReadFile("index.dict")
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
            // WordMap[name]       = DocList[ node, node, ... ]
            WordMap[line_slice[0]] = &DocList{}
            for _, docID := range strings.Split(line_slice[2], " ") {
                WordMap[line_slice[0]].AddDoc(docID)
            }
        }
    }

    // ---[ main part ]--------------------------------------------------------
    // if is without '--mkindex' or with '--useindex'
    if !*mkindex || *useindex {
        if comast == nil {
            // empty string, default value
            fmt.Println("need help? use flag '-help' for help")
            return
        }
        result := AST_result(comast, files_docID, WordMap)
        fmt.Println("result:", result.Str())
        // --------------------------------------------------------------------
    }
    // ------------------------------------------------------------------------
    return
}

// copyleft: PeterlitsZo<peterlitszo@outlook.com>

//
//                                    ___ 
//   ___ _ __ ___  ___ _ __   ___ _ _|__ \
//  / __| '__/ _ \/ _ \ '_ \ / _ \ '__|/ /
// | (__| | |  __/  __/ |_) |  __/ |  |_| 
//  \___|_|  \___|\___| .__/ \___|_|  (_) 
//                    |_|                 
//

//                                                                             
//   __ ___      ____      ____      ____      ____      ____      ____      __
//  / _` \ \ /\ / /\ \ /\ / /\ \ /\ / /\ \ /\ / /\ \ /\ / /\ \ /\ / /\ \ /\ / /
// | (_| |\ V  V /  \ V  V /  \ V  V /  \ V  V /  \ V  V /  \ V  V /  \ V  V / 
//  \__,_| \_/\_/    \_/\_/    \_/\_/    \_/\_/    \_/\_/    \_/\_/    \_/\_/  
//                                                                             
//                        
//  _ __ ___   __ _ _ __  
// | '_ ` _ \ / _` | '_ \ 
// | | | | | | (_| | | | |
// |_| |_| |_|\__,_|_| |_|
//                        
