package main

import (
    "github.com/akamensky/argparse"
    "os"
)

func EastArgparse() (parser *argparse.Parser, dirpath *string, command *string,
                mkindex *bool, useindex *bool, interactive *bool,
                version *bool, err error) {
    // ---[ parse args ]-------------------------------------------------------

    parser = argparse.NewParser("East", "sreach engine on file system")

    // initial the arg-parser
    dirpath = parser.String("d", "dirpath",
                             &argparse.Options{
                                 Required: false,
                                 Help: "the input files' folder path",
                                 Default: "input",
                             })
    command = parser.String("c", "command",
                             &argparse.Options{
                                 Required: false,
                                 Help: "the command to get the ID list (see README.pdf)",
                                 Default: "",
                             })
    mkindex = parser.Flag("m", "mkindex",
                           &argparse.Options{
                               Help: "use this flag to make index named 'index.dict'",
                           })
    useindex = parser.Flag("u", "useindex",
                            &argparse.Options{
                                Help: "use file 'index.dict' to find result",
                            })
    interactive = parser.Flag("i", "interactive",
                               &argparse.Options{
                                   Help: "make self under the interactive mode",
                               })
    version = parser.Flag("v", "version",
                           &argparse.Options{
                               Help: "Show East's version",
                           })
    err = parser.Parse(os.Args)

    // return the value
    return
}
