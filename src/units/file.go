package units

import (
    "io/ioutil"
    "os"
)

type File struct {
    Path string
    Name string
}

func GetFiles(dirpath string) (files []File, err error) {
    dir, err := ioutil.ReadDir(dirpath)
    if err != nil {
        return nil, err
    }
    for _, file := range dir {
        if !file.IsDir() {
            file_path := dirpath + string(os.PathSeparator) + file.Name()
            files = append(files, File{Path: file_path, Name: file.Name()})
        }
    }
    return files, nil
}
