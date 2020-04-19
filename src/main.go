package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "bufio"
    "os"

    "github.com/akamensky/argparse"
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

// ----------------------------------------------------------------------------
// ---[ MAIN FUNCTION ]--------------------------------------------------------
// ----------------------------------------------------------------------------

func main() {
    // ---[ parse args ]-------------------------------------------------------

    parser := argparse.NewParser("East", "sreach engine on file system")

    // initial the arg-parser
    dirpath := parser.String("d", "dirpath",
                             &argparse.Options{
                                 Required: false,
                                 Help: "the input files' folder path",
                                 Default: "input",
                             })
    command := parser.String("c", "command",
                             &argparse.Options{
                                 Required: false,
                                 Help: "the command to get the ID list (see README.pdf)",
                                 Default: "",
                             })
    mkindex := parser.Flag("m", "mkindex",
                           &argparse.Options{
                               Help: "use this flag to make index named 'index.dict'",
                           })
    useindex := parser.Flag("u", "useindex",
                            &argparse.Options{
                                Help: "use file 'index.dict' to find result",
                            })
    interactive := parser.Flag("i", "interactive",
                               &argparse.Options{
                                   Help: "make self under the interactive mode",
                               })
    err := parser.Parse(os.Args)
    if err != nil {
        fmt.Println(parser.Usage(err))
        return
    }

    // ---[ initial all variable ]---------------------------------------------

    // get all Files under the given folder
    files, err := getFiles(*dirpath)
    // files_docID is docID of all Files
    files_docID := &DocList{}
    for _, file := range files {
        files_docID.AddDoc(file.name)
    }
    if err != nil { fmt.Println("ERROR:", err); return; }

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

    var comast *typeAst
    // parser the command ast
    if *command != "" || !*interactive {
        comast = getAST(*command)

        // if is without '--mkindex' or with '--useindex'
        if !*mkindex || *useindex {
            if comast == nil {
                // empty string, default value
                fmt.Println("need help? use flag '-help' or read README for help")
                return
            }
            result := AST_result(comast, files_docID, WordMap)
            fmt.Println("result:", result.Str())
            // --------------------------------------------------------------------
        }
    // interactive mode!!!
    } else {
        fmt.Println("Enter `quit` for quit")
        fmt.Println("copyleft (C) Peterlits Zo <peterlitszo@outlook.com>")
        fmt.Println("Github: github.com/PeterlitsZo/East")
        fmt.Println("")
        for {
            reader := bufio.NewReader(os.Stdin)
            fmt.Print("Command > ")
            text, _ := reader.ReadString('\n')
            text = strings.Replace(text, "\n", "", -1)
            if text == "quit" {
                return
            }
            comast = getAST(text)
            if comast == nil {
                fmt.Println("need help? use flag '-help' or read README for help")
                continue
            }
            result := AST_result(comast, files_docID, WordMap)
            fmt.Println("result:", result.Str())
        }
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

// 
//                             _                _      _       _        
//  ___  ___   __      _____  | |__   __ _  ___| | __ (_)_ __ | |_ ___  
// / __|/ _ \  \ \ /\ / / _ \ | '_ \ / _` |/ __| |/ / | | '_ \| __/ _ \ 
// \__ \ (_) |  \ V  V /  __/ | |_) | (_| | (__|   <  | | | | | || (_) |
// |___/\___/    \_/\_/ \___| |_.__/ \__,_|\___|_|\_\ |_|_| |_|\__\___/ 
//                                                                      
//            _            
//  _ __ ___ (_)_ __   ___ 
// | '_ ` _ \| | '_ \ / _ \
// | | | | | | | | | |  __/
// |_| |_| |_|_|_| |_|\___|
//                         

// 
//              
//              
//              
//  _ _ _ _ _ _ 
// (_|_|_|_|_|_)
//              
