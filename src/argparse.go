package main

import (
    "github.com/PeterlitsZo/argparse"
    "os"
)


func EastArgparse() (parser *argparse.Parser, dirpath *string, command *string,
                mkindex *bool, useindex *bool, interactive *bool,
                version *bool, err error) {
    // ---[ parse args ]-------------------------------------------------------

    parser = argparse.NewParser("East", "sreach engine on file system")

    // initial the arg-parser
    // 
    // the input files' folder path
    dirpath = parser.String(
        "d", "dirpath",
        &argparse.Options{
            Required: false,
            Help: "the input files' folder path",
            Default: "input",
    })

    // the short command usage
    command = parser.String(
        "c", "command",
        &argparse.Options{
            Required: false,
            Help: "the command to get the ID list (see README.pdf)",
            Default: "",
    })

    // to make a index file or not
    mkindex = parser.Flag(
        "m", "mkindex",
        &argparse.Options{
            Help: "use this flag to make index named 'index.dict'",
    })

    // use the index file or not
    useindex = parser.Flag(
        "u", "useindex",
        &argparse.Options{
            Help: "use file 'index.dict' to find result",
    })

    // to use interactive or not
    interactive = parser.Flag(
        "i", "interactive",
        &argparse.Options{
            Help: "make self under the interactive mode",
    })

    // to show self's version
    version = parser.Flag(
        "v", "version",
        &argparse.Options{
            Help: "Show East's version",
    })

    // parse the argument, and if there is error, then raise it out
    err = parser.Parse(os.Args)

    // return the value
    return
}
