package main

type DocList struct {
    length int
    start *DocListNote
}

type DocListNote struct {
    docID string
    pre   *DocListNote
    next  *DocListNote
}

// Add `docID` into DocList
func (list *DocList) AddDoc(docID string) {
    // []: this list should hold its increasing list
    // if its length is 0, then use a new note as its start
    if list.length == 0 {
        list.start = &DocListNote{ docID: docID }
    // there are already some nodes in its list
    } else {
        current := list.start
        // choose its order
        for current.next != nil && docID > current.docID {
            current = current.next
        }
        // the element need be different with the others
        if docID == current.docID {
            return
        }
        // ---[ INSERT ]-------------------------------------------------------
        // insert the element after the tail element
        if docID > current.docID {
            current.next = &DocListNote {
                docID: docID,
                pre: current,
                next: nil,
            }
        // (now docID < current.docID) insert the element before the current node
        } else if current.pre != nil {
            pre_backup := current.pre
            current.pre = &DocListNote {
                docID: docID,
                pre: current.pre,
                next: current,
            }
            pre_backup.next = current.pre
        // else it must be at the head of the list
        } else {
            current.pre = &DocListNote {
                docID: docID,
                pre: nil,           // the head of the list
                next: current,      // link it to next node( current node )
            }
            list.start = current.pre
        }
    }
    list.length += 1
}

// Remove `docID` in docID
func (list *DocList) RemoveDoc(docID string) {
    current := list.start
    // if there are not any node in its list
    if current == nil {
        return
    }
    for current.next != nil && current.docID != docID {
        current = current.next
    }
    // now the current maybe is at the tail or point the corrent docID
    if current.docID == docID {
        // normal node
        if current.pre != nil && current.next != nil {
            current.pre.next = current.next
            current.next.pre = current.pre
        }
        // start node
        if current.pre == nil && current.next != nil{
            list.start = current.next
            current.next.pre = nil
        }
        // tail node
        if current.next == nil && current.pre != nil{
            current.pre.next = nil
        }
        // the only node
        if current.next == nil && current.pre == nil{
            list.start = nil
        }
    }
    list.length -= 1
}

// try to find `docID` and then return the boolean result
func (list *DocList) Has(docID string) bool {
    current := list.start
    for current != nil {
        // get it and then return true
        if current.docID == docID {
            return true
        }
        current = current.next
    }
    return false
}

// clear self to let self be a empty list
func (list *DocList) Clear() {
    list.start = nil;
    list.length = 0;
}

// make self's value is sample with other list
// if self's length is not 0, then try tto clear and then copy
func (list *DocList) Copy(other *DocList) {
    list.Clear()
    current := other.start
    for current != nil {
        list.AddDoc(current.docID)
        current = current.next
    }
}

// return the string of DocList object
// -----------------------------------
// e.g. (DocList{"1", "2", "3", "4"} -> [1, 2, 3, 4]
func (list *DocList) Str() string {
    // if it is a empty DocList then return "[ ]"
    if list.length == 0 {
        return "[ ]"
    // else it will be build by head "[", body and " ]"
    } else {
        current := list.start
        result := "[ "
        for current.next != nil {
            result += current.docID
            result += ", "
            current = current.next
        }
        result += current.docID
        result += " ]"
        return result
    }
}

