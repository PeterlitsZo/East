package main

import (
    "io/ioutil"
    "os"
)

type File struct {
    path string
    name string
}

func GetFiles(dirpath string) (files []File, err error) {
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
