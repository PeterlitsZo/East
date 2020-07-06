package index

import (
    "math"
)

// this function will need those vectors of string, so that they can get the
// value of df of each word: type: map[string]int
func GetRawTf(content string) (vector map[string]int) {
    vector = make(map[string]int)
    splited_content := Split(content)
    for _, word := range splited_content {
        vector[word] += 1
    }
    return
}

func GetTf(content string) (vector map[string]float64) {
    vector = make(map[string]float64)
    raw_tf := GetRawTf(content)
    for word, number := range raw_tf {
        vector[word] = 1 + math.Log(float64(number))
    }
    return
}
