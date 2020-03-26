package main

import (
    "regexp"
)

type token struct {
    name string
    value string
}

type aim struct {
    aim bool
    value string
}

func Bool_Tokenstream(s string) ([]token, string) {
    result := []token{}
    restr_map := map[string]string{
        "AND" : "^AND\\b",
        "NOT" : "^NOT\\b",
        "STR" : "^(\"[^\"]*\"|'[^']*')",
        "SPA" : "^\\s+",
    }
    re_map := map[string]*regexp.Regexp{}
    for key, value := range restr_map {
        re_map[key] = regexp.MustCompile(value)
    }
    for index := 0; index < len(s); {
        index_backup := index
        for token_key, token_re := range re_map {
            if token_re.Match([]byte(s[index:])) {
                result = append(result, token{
                    name: token_key,
                    value: string(token_re.Find([]byte(s[index:]))),
                })
                index += len(result[len(result)-1].value)
                break
            }
        }
        if index == index_backup {
            return nil, "do not match all"
        }
    }
    return result, ""
}

func Aim_Parser(s string) ([]aim, string){
    tokstream, errmsg := Bool_Tokenstream(s)
    if errmsg != "" {
        return nil, errmsg
    }
    and_flag := true
    not_flag := false
    result := []aim{}
    for _, tok := range tokstream {
        if tok.name == "SPA" {
            continue
        }
        if tok.name == "STR" {
            if and_flag {
                if not_flag {
                    result = append(result, aim{aim: false, value: tok.value})
                    not_flag = true
                } else {
                    result = append(result, aim{aim: true, value: tok.value})
                }
                and_flag = false
            } else {
                return nil, "need AND opertor"
            }
        }
        if tok.name == "AND" {
            and_flag = true
        }
        if tok.name == "NOT" {
            not_flag = true
        }
    }
    return result, ""
}
