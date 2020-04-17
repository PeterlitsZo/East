%{
// ----------------------------------------------------------------------------
// ---[ begin of source file head ]--------------------------------------------

// ---[ need github.com/pingcap/parser ]---------------------------------------

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

// the part of parser's type that it want:
type typeAtom struct {
    not bool
    value interface{}
}
type typeExpr []*typeAtom
type typeAst  []*typeExpr

var ast_result *typeAst

// ---[ end of source file head ]----------------------------------------------
// ----------------------------------------------------------------------------
%}

// the part of parser's type that it want:
%union{
    Atom *typeAtom
    Expr *typeExpr
    Ast  *typeAst
    Str  string
}

// <<< %token AND OR NOT
%token AND OR NOT '(' ')'
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
            | NOT '(' ast ')' {
                $$ = &typeAtom{true, $3}
            }
            | '(' ast ')' {
                $$ = &typeAtom{false, $2}
            }
            ;

%%
// ----------------------------------------------------------------------------
// ---[ begin of the tail source file ]----------------------------------------

// ---[ Parser ]---------------------------------------------------------------

// the regex machine
var re = map[int]*regexp.Regexp{
    // e.g. and, AND, And, &&
    AND: regexp.MustCompile(`^([aA][nN][dD]|&&)`),
    // e.g. or, OR, Or, ||
    OR:  regexp.MustCompile(`^([oO][rR]|\|\|)`),
    // e.g. not NOT, Not, !
    NOT: regexp.MustCompile(`^([nN][oO][tT]|!)`),
    // e.g. "PETER", "\"", '\'', '"'
    STR: regexp.MustCompile(`^("(\\"|[^"])*"|'(\\'|[^'])*')`),
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
        // e.g. AND, And, and, &&
        case re[AND].Match(l.input[l.pos:]):
            l.pos += len(re[AND].Find(l.input[l.pos:]))
            return AND
        // e.g. OR, Or, or, ||
        case re[OR].Match(l.input[l.pos:]):
            l.pos += len(re[OR].Find(l.input[l.pos:]))
            return OR
        // e.g. NOT, Not, not, !
        case re[NOT].Match(l.input[l.pos:]):
            l.pos += len(re[NOT].Find(l.input[l.pos:]))
            return NOT
        // e.g. "", '', "\"", "abc"
        case re[STR].Match(l.input[l.pos:]):
            str_result := re[STR].Find(l.input[l.pos:])
            l.pos += len(str_result)
            // let itself has the value that it want.
            lval.Str = getstring(string(str_result))
            return STR
        // match "(": "("'s length is 1
        case string(l.input[l.pos:l.pos+1]) == "(":
            l.pos += len("(")
            return int('(')
        // match ")": ")"'s length is 1
        case string(l.input[l.pos:l.pos+1]) == ")":
            l.pos += len(")")
            return int(')')
        // match error
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

// return the typeExpr object's format string
func (expr *typeExpr) String() string {
    result := ""
    result += "---[ EXPR ]---\n"
    for _, v := range *expr {
        if v.not { result += "<NOT> "; }
        // TODO: show the atom's string
        // result += v.value
        result += "\n"
    }
    return result
}

// return the typeAst object's format string
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

// from a string to build a AST( if s is empty then return a nil pointer )
func getAST(s string) *typeAst {
    if s == "" {
        return nil
    }
    yyParse(&GoLex{input: []byte(s)})
    // [WARNING] ast_result is a global variable
    return ast_result
}

// ---[ test ]-----------------------------------------------------------------

/*
    // test main function
    func main(){
        log.SetFlags(log.Ldate|log.Lshortfile)
        log.Println("parsing command ...")
        getAST([]byte("'me' AND 'bala' OR 'hey' AND 'hellp'"))
    }
*/

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------


