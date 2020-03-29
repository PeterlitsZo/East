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

func (list *DocList) AddDoc(docID string) {
    if list.length == 0 {
        list.start = &DocListNote{ docID: docID }
    } else {
        current := list.start
        // choose its order
        for current.next != nil && current.docID < docID {
            current = current.next
        }
        // the element need be different with the others
        if current.docID == docID {
            return
        }
        // insert the element
        next_backup := current.next
        current.next = &DocListNote {
            docID: docID,
            pre: current,
            next: next_backup,
        }
        if next_backup != nil {
            next_backup.pre = current
        }
    }
    list.length += 1
}

func (list *DocList) RemoveDoc(docID string) {
    current := list.start
    if current == nil {
        return
    }
    for current.next != nil && current.docID != docID {
        current = current.next
    }
    // now the current maybe is at the tail or point the corrent docID
    if current.docID == docID {
        if current.pre != nil && current.next != nil {
            current.pre.next = current.next
            current.next.pre = current.pre
        }
        if current.pre == nil && current.next != nil{
            // start node
            list.start = current.next
            current.next.pre = nil
        }
        if current.next == nil && current.pre != nil{
            // tail node
            current.pre.next = nil
        }
        if current.next == nil && current.pre == nil{
            list.start = nil
        }
    }
}

func (list *DocList) Has(docID string) bool {
    current := list.start
    for current != nil {
        if current.docID == docID {
            return true
        }
        current = current.next
    }
    return false
}

func (list *DocList) Str() string {
    if list.length == 0 {
        return "[ ]"
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
