package main

import (
    "github.com/PeterlitsZo/argparse"
    "os"
)


type ParserRun struct {
     self       *argparse.Command
     useindex   *bool
     command    *string
     dirpath    *string
}

type ParserMkindex struct {
    self        *argparse.Command
    dirpath     *string
}

type ParserInteractive struct {
    self        *argparse.Command
    useindex    *bool
    dirpath     *string
}

type ParserVersion struct {
    self        *argparse.Command
}


type ParserResult struct {
    parser       *argparse.Parser
    run          ParserRun
    mkindex      ParserMkindex
    interactive  ParserInteractive
    version      ParserVersion
    err          error
}


func EastArgparse() (pr ParserResult) {
    // ---[ parse args ]-------------------------------------------------------

    pr.parser = argparse.NewParser("East", "sreach engine on file system")

    // command
    // the short command usage
    pr.run.self = pr.parser.NewCommand(
        "run", "the command to get the ID list (see README.pdf)",
        )
    pr.run.useindex = pr.run.self.Flag(
        "u", "useindex", &argparse.Options{
            Help: "use file 'index.dict' to find result",
        })
    pr.run.command = pr.run.self.String(
        "c", "command", &argparse.Options{
            Required: true,
            Help: "the command to run (see README.pdf for more infos)",
        })
    pr.run.dirpath = pr.run.self.String(
        "d", "dirpath", &argparse.Options{
            Required: false,
            Help: "the command to run (see README.pdf for more infos)",
            Default: "input",
        })

    // mkindex
    // to make a index file or not
    pr.mkindex.self = pr.parser.NewCommand(
        "mkindex", "use this flag to make index named 'index.dict'")
    pr.mkindex.dirpath = pr.mkindex.self.String(
        "d", "dirpath", &argparse.Options{
            Required: false,
            Help: "the command to run (see README.pdf for more infos)",
            Default: "input",
        })

    // interactive
    // to use interactive or not
    pr.interactive.self = pr.parser.NewCommand(
        "interactive", "make self under the interactive mode",
        )
    pr.interactive.useindex = pr.interactive.self.Flag(
        "u", "useindex", &argparse.Options{
            Help: "use file 'index.dict' to find result",
        })
    pr.interactive.dirpath = pr.interactive.self.String(
        "d", "dirpath", &argparse.Options{
            Required: false,
            Help: "the command to run (see README.pdf for more infos)",
            Default: "input",
        })

    // version
    // to show self's version
    pr.version.self = pr.parser.NewCommand(
        "version", "Show East's version")

    // parse the argument, and if there is error, then raise it out
    pr.err = pr.parser.Parse(os.Args)

    // return the value
    return
}
