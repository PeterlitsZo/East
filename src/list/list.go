package list

type DocList struct {
    Length int
    Start *DocListNote
}

type DocListNote struct {
    DocID string
    Pre   *DocListNote
    Next  *DocListNote
}

// Add `DocID` into DocList
func (list *DocList) AddDoc(DocID string) {
    // []: this list should hold its increasing list
    // if its length is 0, then use a new note as its start
    if list.Length == 0 {
        list.Start = &DocListNote{ DocID: DocID }
    // there are already some nodes in its list
    } else {
        current := list.Start
        // choose its order
        for current.Next != nil && DocID > current.DocID {
            current = current.Next
        }
        // the element need be different with the others
        if DocID == current.DocID {
            return
        }
        // ---[ INSERT ]-------------------------------------------------------
        // insert the element after the tail element
        if DocID > current.DocID {
            current.Next = &DocListNote {
                DocID: DocID,
                Pre: current,
                Next: nil,
            }
        // (now DocID < current.DocID) insert the element before the current node
        } else if current.Pre != nil {
            Pre_backup := current.Pre
            current.Pre = &DocListNote {
                DocID: DocID,
                Pre: current.Pre,
                Next: current,
            }
            Pre_backup.Next = current.Pre
        // else it must be at the head of the list
        } else {
            current.Pre = &DocListNote {
                DocID: DocID,
                Pre: nil,           // the head of the list
                Next: current,      // link it to next node( current node )
            }
            list.Start = current.Pre
        }
    }
    list.Length += 1
}

// Remove `DocID` in DocID
func (list *DocList) RemoveDoc(DocID string) {
    current := list.Start
    // if there are not any node in its list
    if current == nil {
        return
    }
    for current.Next != nil && current.DocID != DocID {
        current = current.Next
    }
    // now the current maybe is at the tail or point the corrent DocID
    if current.DocID == DocID {
        // normal node
        if current.Pre != nil && current.Next != nil {
            current.Pre.Next = current.Next
            current.Next.Pre = current.Pre
        }
        // start node
        if current.Pre == nil && current.Next != nil{
            list.Start = current.Next
            current.Next.Pre = nil
        }
        // tail node
        if current.Next == nil && current.Pre != nil{
            current.Pre.Next = nil
        }
        // the only node
        if current.Next == nil && current.Pre == nil{
            list.Start = nil
        }
    }
    list.Length -= 1
}

// try to find `DocID` and then return the boolean result
func (list *DocList) Has(DocID string) bool {
    current := list.Start
    for current != nil {
        // get it and then return true
        if current.DocID == DocID {
            return true
        }
        current = current.Next
    }
    return false
}

// clear self to let self be a empty list
func (list *DocList) Clear() {
    list.Start = nil;
    list.Length = 0;
}

// make self's value is sample with other list
// if self's length is not 0, then try tto clear and then copy
func (list *DocList) Copy(other *DocList) {
    list.Clear()
    current := other.Start
    for current != nil {
        list.AddDoc(current.DocID)
        current = current.Next
    }
}

// return the string of DocList object
// -----------------------------------
// e.g. (DocList{"1", "2", "3", "4"} -> [1, 2, 3, 4]
func (list *DocList) Str() string {
    // if it is a empty DocList then return "[ ]"
    if list.Length == 0 {
        return "[ ]"
    // else it will be build by head "[", body and " ]"
    } else {
        current := list.Start
        result := "[ "
        for current.Next != nil {
            result += current.DocID
            result += ", "
            current = current.Next
        }
        result += current.DocID
        result += " ]"
        return result
    }
}

