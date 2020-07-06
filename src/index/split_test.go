package index

import (
    "testing"
)

func TestSplit(t *testing.T) {
    result := Split("This is a testing text. This is really a testing text!")
    expected := []string{"This", "is", "a", "testing", "text", "This",
                         "is", "really", "a",  "testing", "text"}
    for index := range result {
        if result[index] != expected[index] {
            t.Errorf("I got %#v, but I want %#v", result, expected)
        }
    }
}
