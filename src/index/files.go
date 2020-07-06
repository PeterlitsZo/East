package index

import (
    "io/ioutil"
    "os"
    // "errors"
)

type file struct {
    Name   string
    Vector map[string]float64
}

type Files []file

func getFilesContent(folderPath string) (result map[string]string, err error) {
    result = make(map[string]string)
    dir, err := ioutil.ReadDir(folderPath)
    if err != nil {
        return nil, err
    }
    for _, file := range dir {
        if !file.IsDir() {
            path := folderPath + string(os.PathSeparator) + file.Name()
            content, err := ioutil.ReadFile(path)
            result[file.Name()] = string(content)
            if err != nil {
                return nil, err
            }
        }
    }
    return
}

func getFiles(path string) (result Files) {
    ;
    return
}

