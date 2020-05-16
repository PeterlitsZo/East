package units

import (
    "io/ioutil"
    "os"

    "../list"
)

type File struct {
    Path string
    Name string
}

func _getFiles(dirpath string) (files []File, err error) {
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

func GetFiles(path string) (files []File, files_docid *list.DocList, err error) {
    files_origin, err := _getFiles(path)
    files = files_origin
    if err != nil {
        return nil, nil, err
    }
    files_docid = &list.DocList{}
    for _, file := range files {
        files_docid.AddDoc(file.Name)
    }
    return
}


