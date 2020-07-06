package index

import (
    "strings"
)

func Split(text string) []string {
    raw_slice := strings.Split(text, " ")
    result := []string{}
    for _, item := range raw_slice {
        item = strings.Trim(item, " ,.!~?\"'()-\r\n")
        if item != "" {
            result = append(result, item)
        }
    }
    return result
}
