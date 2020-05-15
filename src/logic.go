package main


import (
    "reflect"
    "fmt"
)

// ---[ return bool sreach's result ]------------------------------------------
// return the AST's result.
func AST_result(list *typeList, all_docID *DocList, wordmap map[string]*DocList) *DocList {
    result := &DocList{}
    // if the list's len is zero, then return the full DocList
    if list == nil || len(*list) == 0 {
        result.Copy(all_docID)
        return result
    }
    for _, expr_ptr := range *list {
        sub_result := EXPR_result(expr_ptr, all_docID, wordmap)
        current := sub_result.start
        for current != nil {
            // get the union of all sub_result
            result.AddDoc(current.docID)
            current = current.next
        }
    }
    return result
}

// return the result of expr:
// --------------------------
// return the expr's result. all atom is link by op 'and'
func EXPR_result(expr *typeExpr, all_docID *DocList, wordmap map[string]*DocList) *DocList {
    // if the expr's len is zero, then return empty DocList
    result := &DocList{}
    if len(*expr) == 0 {
        return result
    // else initial it by the first node
    } else {
        result.Copy(ATOM_result((*expr)[0], all_docID, wordmap))
    }
    for _, atom_ptr := range *expr {
        sub_result := ATOM_result(atom_ptr, all_docID, wordmap)
        current := result.start
        // iter the sub-result
        for current != nil {
            // only need the node that the sub-result has
            if !sub_result.Has(current.docID) {
                result.RemoveDoc(current.docID)
            }
            current = current.next
        }
    }
    return result
}

// return the result of atom:
// --------------------------
// atom's value is a interface, we need know that the type of the value. if the
// type is typeAst, then it need call AST, else it should be a string, so we
// need get the result by wordmap. if atom.not, then need to negate it by full
// docID list.
func ATOM_result(atom *typeAtom, all_docID *DocList, wordmap map[string]*DocList) *DocList {
    result := &DocList{}
    switch v := atom.value.(type) {
    case string:
        doclist_ptr, ok := wordmap[atom.value.(string)]
        // if it have the key then copy else return empty list, else copy the 
        if ok {
            result.Copy(doclist_ptr)
        }
    case *typeList:
        sub_ast := atom.value.(*typeList)
        result.Copy(AST_result(sub_ast, all_docID, wordmap))
    default:
        // TODO: this is not OK( it look ugly )
        fmt.Println("Error! the type is ", reflect.TypeOf(v))
        return result
    }
    // negate the result
    if atom.not {
        current := all_docID.start
        for current != nil {
            // does not have? make it have
            if !result.Has(current.docID) {
                result.AddDoc(current.docID)
            // does it have? remove it
            } else {
                result.RemoveDoc(current.docID)
            }
            current = current.next
        }
    }
    return result
}

