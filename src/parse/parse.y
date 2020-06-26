%{
// ----------------------------------------------------------------------------
// ---[ begin of source file head ]--------------------------------------------

// ---[ need github.com/pingcap/parser ]---------------------------------------

package parse

import (
    "log"
    "regexp"
    "unicode"
    "unicode/utf8"
)

const (
    Debug = 4
    ErrorVerbose = true
)

// the part of parser's type that it want:
type TypeAst struct {
    Command string
    Value interface{}
}

// if the Command is "sreach"
type TypeAtom struct {
    Not bool
    Value interface{}
}
type TypeExpr []*TypeAtom
type TypeList []*TypeExpr

var ast_result *TypeAst

// ---[ end of source file head ]----------------------------------------------
// ----------------------------------------------------------------------------
%}

// the part of parser's type that it want:
%union{
    Atom *TypeAtom
    Expr *TypeExpr
    List *TypeList
    Ast  *TypeAst
    Str  string
}

%token SREACH LIST PRINT
%token AND OR NOT '(' ')'
%token <Str> STR

%type  <Ast>  ast
%type  <List> sreach_word
%type  <Expr> expr
%type  <Atom> atom


%%

top         : ast {
                ast_result = $1
            }

ast         : SREACH sreach_word {
                $$ = &TypeAst{"sreach", $2}
            }
            | LIST {
                $$ = &TypeAst{"list", nil}
            }
            | PRINT STR{
                $$ = &TypeAst{"print", $2}
            }

sreach_word : expr OR sreach_word {
                temp := append( *$3, $1 )
                $$ = &temp
            }
            | expr {
                temp := make( TypeList, 0 )
                temp = append( temp, $1 )
                $$ = &temp
            }
            ;

expr        : atom AND expr {
                temp := append( *$3, $1 )
                $$ = &temp
            }
            | atom {
                temp := make( TypeExpr, 0 )
                temp = append( temp, $1 )
                $$ = &temp
            }
            ;

atom        : NOT STR {
                $$ = &TypeAtom{true, $2}
            }
            | STR {
                $$ = &TypeAtom{false, $1}
            }
            | NOT '(' sreach_word ')' {
                $$ = &TypeAtom{true, $3}
            }
            | '(' sreach_word ')' {
                $$ = &TypeAtom{false, $2}
            }
            ;

%%
// ----------------------------------------------------------------------------
// ---[ begin of the tail source file ]----------------------------------------

// ---[ Parser ]---------------------------------------------------------------

// the regex machine
var re = map[int]*regexp.Regexp{
    // e.g. sreach, Sreach, SREACH
    SREACH: regexp.MustCompile(`^[sS][rR][eE][aA][cC][hH]`),
    // e.g. list
    LIST:   regexp.MustCompile(`^[lL][iI][sS][tT]`),
    // e.g print
    PRINT:  regexp.MustCompile(`^[pP][rR][iI][nN][tT]`),
    // e.g. and, AND, And, &&
    AND:    regexp.MustCompile(`^([aA][nN][dD]|&&)`),
    // e.g. or, OR, Or, ||
    OR:     regexp.MustCompile(`^([oO][rR]|\|\|)`),
    // e.g. not NOT, Not, !
    NOT:    regexp.MustCompile(`^([nN][oO][tT]|!)`),
    // e.g. "PETER", "\"", '\'', '"'
    STR:    regexp.MustCompile(`^("(\\"|[^"])*"|'(\\'|[^'])*')`),
}

// the struct of input (member input is the string of its input)
// poos hold its postion
type GoLex struct {
    pos    int
    input  []byte
}

// get the string that command want to mean:
// -----------------------------------------
//     e.g. "\"" -> "(a string), '\'' -> '(a string), 'abc' -> abc(a string),
// '"' -> "(a string), "string" -> string(a string)
func getstring(org string) string {
    result := ""
    // it is a little status machine( IN ESPCAE: '"...\', NOT IN ESPCE: '"...'(
    // mormal string that do not have char '\')
    IN_ESPCAE_STATUS := false
    for _, v := range org[1:len(org)-1] {
        if IN_ESPCAE_STATUS {
            // if now in espcae status, it will espcae char '\', '\"', and "\'"
            if v == '\\' || v == '"' || v == '\'' {
                result += string(v)
                // back to normal status
                IN_ESPCAE_STATUS = false
            } else {
            // it do not matter, we do not need to espcae it, so we back the
            // normal string( e.g. "\a" -> "\a" )
                result += string('\\')
                result += string(v)
            }
        } else if v == '\\' {
            // if match the char '\', then change self's status
            IN_ESPCAE_STATUS = true
        } else {
            // in normal status
            result += string(v)
        }
    }
    return result
}



// lexer: get a long command then return the token stream that the yaccer need
// ---------------------------------------------------------------------------
//     lval(type: *yySymType): the atom of the token stream, which is need for
// yacc
//     l(type: *GoLex): the input string(and a interger `pos` to remember its
// position
func (l *GoLex) Lex(lval *yySymType) int {
    for l.pos < len(l.input){
        // get the next Rune(part of string) and its length
        r, n := utf8.DecodeRune(l.input[l.pos:])
        // if it is a space, then skip it.
        if unicode.IsSpace(r) {
            l.pos += n
            continue
        }
        switch{
        case re[SREACH].Match(l.input[l.pos:]):
            l.pos += len(re[SREACH].Find(l.input[l.pos:]))
            return SREACH

        case re[LIST].Match(l.input[l.pos:]):
            l.pos += len(re[LIST].Find(l.input[l.pos:]))
            return LIST

        case re[PRINT].Match(l.input[l.pos:]):
            l.pos += len(re[LIST].Find(l.input[l.pos:]))
            return PRINT

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
            // let itself has the value that it want.
            lval.Str = getstring(string(str_result))
            return STR

        case string(l.input[l.pos:l.pos+1]) == "(":
            l.pos += len("(")
            return int('(')

        case string(l.input[l.pos:l.pos+1]) == ")":
            l.pos += len(")")
            return int(')')

        default:
            log.Println("can't match", "\"" + string(l.input[l.pos:]) + "\"")
            return 0
        }
    }
    return 0
}

// when l can't match
func (l *GoLex) Error(s string) {
    log.Printf("syntax error: %s\n", s)
}

// ---[ AST ]------------------------------------------------------------------

// from a string to build a AST( if s is empty then return a nil pointer )
func GetAST(s string) *TypeAst {
    if s == "" {
        return nil
    }
    yyParse(&GoLex{input: []byte(s)})
    // [WARNING] ast_result is a global variable
    return ast_result
}


