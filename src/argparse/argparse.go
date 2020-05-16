package argparse

import (
    "github.com/PeterlitsZo/argparse"
    "os"
)


type ParserRun struct {
     Self       *argparse.Command
     Useindex   *bool
     Command    *string
     Dirpath    *string
}

type ParserMkindex struct {
    Self        *argparse.Command
    Dirpath     *string
}

type ParserInteractive struct {
    Self        *argparse.Command
    Useindex    *bool
    Dirpath     *string
}

type ParserVersion struct {
    Self        *argparse.Command
}


type ParsErresult struct {
    Parser       *argparse.Parser
    Run          ParserRun
    Mkindex      ParserMkindex
    Interactive  ParserInteractive
    Version      ParserVersion
    Err          error
}


func EastArgparse() (pr ParsErresult) {
    // ---[ parse args ]-------------------------------------------------------

    pr.Parser = argparse.NewParser("East", "sreach engine on file system")

    // command
    // the short command usage
    pr.Run.Self = pr.Parser.NewCommand(
        "Run", "the command to get the ID list (see README.pdf)",
        )
    pr.Run.Useindex = pr.Run.Self.Flag(
        "u", "UseIndex", &argparse.Options{
            Help: "use file 'index.dict' to find result",
        })
    pr.Run.Command = pr.Run.Self.String(
        "c", "command", &argparse.Options{
            Required: true,
            Help: "the command to Run (see README.pdf for more infos)",
        })
    pr.Run.Dirpath = pr.Run.Self.String(
        "d", "Dirpath", &argparse.Options{
            Required: false,
            Help: "the command to Run (see README.pdf for more infos)",
            Default: "input",
        })

    // Mkindex
    // to make a index file or not
    pr.Mkindex.Self = pr.Parser.NewCommand(
        "Mkindex", "use this flag to make index named 'index.dict'")
    pr.Mkindex.Dirpath = pr.Mkindex.Self.String(
        "d", "Dirpath", &argparse.Options{
            Required: false,
            Help: "the command to Run (see README.pdf for more infos)",
            Default: "input",
        })

    // Interactive
    // to use Interactive or not
    pr.Interactive.Self = pr.Parser.NewCommand(
        "Interactive", "make Self under the Interactive mode",
        )
    pr.Interactive.UseIndex = pr.Interactive.Self.Flag(
        "u", "UseIndex", &argparse.Options{
            Help: "use file 'index.dict' to find result",
        })
    pr.Interactive.Dirpath = pr.Interactive.Self.String(
        "d", "Dirpath", &argparse.Options{
            Required: false,
            Help: "the command to Run (see README.pdf for more infos)",
            Default: "input",
        })

    // Version
    // to show Self's Version
    pr.Version.Self = pr.Parser.NewCommand(
        "Version", "Show East's Version")

    // parse the argument, and if there is Error, then raise it out
    pr.Err = pr.Parser.Parse(os.Args)

    // return the value
    return
}
