package index

import (
    "math"
)

// this will make a if-idf vector of those string slice, it need those strings
// to calc the tf and also need the df of words
func MakeMap(splited_string []string, df map[string]int) (m map[string]float64) {
    content_len := 0

    for _, word := range splited_string {
        m[word] ++
        content_len ++
    }

    for key := range m {
        m[key] = 1 + math.Log10(m[key])
    }

    return
}
