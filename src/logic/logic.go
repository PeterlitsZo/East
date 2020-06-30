package logic


import (
    "reflect"
    "fmt"

    "../units"
    "../parse"
)

// ---[ return the result of command ]-----------------------------------------
// it need a AST as a parameter and deal with it and then return the output as
// a object that can print.
func Run(AST *parse.AST, env *units.Env) interface{} {
    switch AST.Command {
    case "list":
        return "list"
    case "sreach":
        return "sreach"
    case "print":
        return AST.Value
    default:
        return "[Error] Unkown Command"
    }
}

// ---[ return bool sreach's result ]------------------------------------------
// return the AST's result.
func Sreach_result(list_ *parse.ExprList, all_DocID *units.DocList, wordmap map[string]*units.DocList) *units.DocList {
    result := &units.DocList{}
    // if the list_'s len is zero, then return the full units.DocList
    if list_ == nil || len(*list_) == 0 {
        result.Copy(all_DocID)
        return result
    }
    for _, expr_ptr := range *list_ {
        sub_result := EXPR_result(expr_ptr, all_DocID, wordmap)
        current := sub_result.Start
        for current != nil {
            // get the union of all sub_result
            result.AddDoc(current.DocID)
            current = current.Next
        }
    }
    return result
}

// return the result of expr:
// --------------------------
// return the expr's result. all atom is link by op 'and'
func EXPR_result(expr *parse.Expr, all_DocID *units.DocList, wordmap map[string]*units.DocList) *units.DocList {
    // if the expr's len is zero, then return empty units.DocList
    result := &units.DocList{}
    if len(*expr) == 0 {
        return result
    // else initial it by the first node
    } else {
        result.Copy(ATOM_result((*expr)[0], all_DocID, wordmap))
    }
    for _, atom_ptr := range *expr {
        sub_result := ATOM_result(atom_ptr, all_DocID, wordmap)
        current := result.Start
        // iter the sub-result
        for current != nil {
            // only need the node that the sub-result has
            if !sub_result.Has(current.DocID) {
                result.RemoveDoc(current.DocID)
            }
            current = current.Next
        }
    }
    return result
}

// return the result of atom:
// --------------------------
// atom's value is a interface, we need know that the type of the value. if the
// type is typeAst, then it need call AST, else it should be a string, so we
// need get the result by wordmap. if atom.not, then need to negate it by full
// DocID units.
func ATOM_result(atom *parse.Atom, all_DocID *units.DocList, wordmap map[string]*units.DocList) *units.DocList {
    result := &units.DocList{}
    switch v := atom.Value.(type) {
    case string:
        doclist_ptr, ok := wordmap[atom.Value.(string)]
        // if it have the key then copy else return empty list, else copy the 
        if ok {
            result.Copy(doclist_ptr)
        }
    case *parse.ExprList:
        sub_ast := atom.Value.(*parse.ExprList)
        result.Copy(Sreach_result(sub_ast, all_DocID, wordmap))
    default:
        // TODO: this is not OK( it look ugly )
        fmt.Println("Error! the type is ", reflect.TypeOf(v))
        return result
    }
    // negate the result
    if atom.Not {
        current := all_DocID.Start
        for current != nil {
            // does not have? make it have
            if !result.Has(current.DocID) {
                result.AddDoc(current.DocID)
            // does it have? remove it
            } else {
                result.RemoveDoc(current.DocID)
            }
            current = current.Next
        }
    }
    return result
}

