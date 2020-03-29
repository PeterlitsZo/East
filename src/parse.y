%{
// ---[ need github.com/pingcap/parser ]-----------------------------------------------------------

package main

import (
    "log"
    "regexp"
    "strings"
    "unicode"
    "unicode/utf8"
)

const(
    Debug = 4
    ErrorVerbose = true
)

type typeAtom struct {
    not bool
    str string
}
type typeExpr []*typeAtom
type typeAst  []*typeExpr

var ast_result *typeAst

%}

%union{
    Atom *typeAtom
    Expr *typeExpr
    Ast  *typeAst
    Str  string
}

%token AND OR NOT
%token <Str> STR

%type  <Ast>  ast
%type  <Expr> expr
%type  <Atom> atom


%%

top         : ast {
                ast_result = $1
            }

ast         : expr OR ast {
                temp := append( *$3, $1 )
                $$ = &temp
            }
            | expr {
                temp := make( typeAst, 0 )
                temp = append( temp, $1 )
                $$ = &temp
            }
            ;

expr        : atom AND expr {
                temp := append( *$3, $1 )
                $$ = &temp
            }
            | atom {
                temp := make( typeExpr, 0 )
                temp = append( temp, $1 )
                $$ = &temp
            }
            ;

atom        : NOT STR {
                $$ = &typeAtom{true, $2}
            }
            | STR {
                $$ = &typeAtom{false, $1}
            }
            ;

%%

// ---[ Parser ]-----------------------------------------------------------------------------------

var re = map[int]*regexp.Regexp{
    AND: regexp.MustCompile(`^([aA][nN][dD]|&&)`),                  // e.g. and, AND, And, &&
    OR:  regexp.MustCompile(`^([oO][rR]|\|\|)`),                    // e.g. or,  OR,  Or,  ||
    NOT: regexp.MustCompile(`^([nN][oO][tT]|!)`),                   // e.g. not, NOT, Not, !
    STR: regexp.MustCompile(`^("([^"]*(\\")*)*"|'([^']*(\\')*)*')`),// e.g. "PETER", "\"", '\'', '"'
}

type GoLex struct {
    pos    int
    input  []byte
}

func getstring(org string) string {
    result := ""
    IN_ESPCAE_STATUS := false
    for _, v := range org[1:len(org)-1] {
        if IN_ESPCAE_STATUS {
            if v == '\\' || v == '"' || v == '\'' {
                result += string(v)
                IN_ESPCAE_STATUS = false
            } else {
                result += string('\\')
                result += string(v)
            }
        } else if v == '\\' {
            IN_ESPCAE_STATUS = true
        } else {
            result += string(v)
        }
    }
    return result
}



func (l *GoLex) Lex(lval *yySymType) int {
    for l.pos < len(l.input){
        r, n := utf8.DecodeRune(l.input[l.pos:])
        if unicode.IsSpace(r) {
            l.pos += n
            continue
        }
        switch{
        case re[AND].Match(l.input[l.pos:]):
            l.pos += len(re[AND].Find(l.input[l.pos:]))
            return AND
        case re[OR].Match(l.input[l.pos:]):
            l.pos += len(re[OR].Find(l.input[l.pos:]))
            return OR
        case re[NOT].Match(l.input[l.pos:]):
            l.pos += len(re[NOT].Find(l.input[l.pos:]))
            return NOT
        case re[STR].Match(l.input[l.pos:]):
            str_result := re[STR].Find(l.input[l.pos:])
            l.pos += len(str_result)
            lval.Str = getstring(string(str_result))
            return STR
        default:
            log.Println("can't match", "\"" + string(l.input[l.pos:]) + "\"")
            return 0
        }
    }
    return 0
}

func (l *GoLex) Error(s string) {
    log.Printf("syntax error: %s\n", s)
}

// ---[ AST ]--------------------------------------------------------------------------------------

func (expr *typeExpr) String() string {
    result := ""
    result += "---[ EXPR ]---\n"
    for _, v := range *expr {
        if v.not { result += "<NOT> "; }
        result += v.str
        result += "\n"
    }
    return result
}

func (ast *typeAst) String() string {
    result := ""
    result += "---[ AST ]---\n"
    for _, v := range *ast{
        for _, line := range strings.Split(v.String(), "\n") {
            result += "    " + line + "\n"
        }
    }
    return result
}

func getAST(s string) *typeAst {
    if s == "" {
        return nil
    }
    yyParse(&GoLex{input: []byte(s)})
    return ast_result
}

// ---[ test ]-------------------------------------------------------------------------------------

/*
    // test main function
    func main(){
        log.SetFlags(log.Ldate|log.Lshortfile)
        log.Println("parsing command ...")
        getAST([]byte("'me' AND 'bala' OR 'hey' AND 'hellp'"))
    }
*/
