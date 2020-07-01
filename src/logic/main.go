package logic

import (
    "../units"
    "../parse"
)

type RunResult struct {
    NeedBreak bool
    NoOutput  bool
}

// ---[ return the result of command ]-----------------------------------------
// it need a AST as a parameter and deal with it and then return the output as
// a object that can print.
func Run(AST *parse.AST, env *units.Env) (interface{}, RunResult) {
    // if the AST is a nil pointer, it is a error pointer.
    if AST == nil{
        return "<Error: nil AST>", RunResult{
            NeedBreak: false,
            NoOutput: true,
        }
    }

    // switch by the first part -- AST.Command
    switch AST.Command {
    case "list":
        return env.DocID, RunResult{
            NeedBreak: false,
            NoOutput: false,
        }

    case "sreach":
        var exprlist_ptr = AST.Value.(*parse.ExprList)
        return Sreach_result(exprlist_ptr, env.DocID, *env.WordMap), RunResult{
            NeedBreak: false,
            NoOutput: false,
        }

    case "print":
        return AST.Value, RunResult{
            NeedBreak: false,
            NoOutput: false,
        }

    case "quit":
        return nil, RunResult{
            NeedBreak: true,
            NoOutput: true,
        }

    case "empty":
        return nil, RunResult{
            NeedBreak: false,
            NoOutput: true,
        }

    default:
        return "<Error: Unkown Command>", RunResult{
            NeedBreak: false,
            NoOutput: false,
        }
    }
}

