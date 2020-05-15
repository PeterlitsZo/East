// Code generated by goyacc -o ./src/parse.go ./src/parse.y. DO NOT EDIT.

//line ./src/parse.y:2
// ----------------------------------------------------------------------------
// ---[ begin of source file head ]--------------------------------------------

// ---[ need github.com/pingcap/parser ]---------------------------------------

package main

import __yyfmt__ "fmt"

//line ./src/parse.y:7

import (
	"log"
	"regexp"
	"unicode"
	"unicode/utf8"
)

const (
	Debug        = 4
	ErrorVerbose = true
)

// the part of parser's type that it want:
type typeAtom struct {
	not   bool
	value interface{}
}
type typeExpr []*typeAtom
type typeList []*typeExpr
type typeAst struct {
	command string
	value   interface{}
}

var ast_result *typeAst

// ---[ end of source file head ]----------------------------------------------
// ----------------------------------------------------------------------------

//line ./src/parse.y:40
type yySymType struct {
	yys  int
	Atom *typeAtom
	Expr *typeExpr
	List *typeList
	Ast  *typeAst
	Str  string
}

const SREACH = 57346
const LIST = 57347
const AND = 57348
const OR = 57349
const NOT = 57350
const STR = 57351

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"SREACH",
	"LIST",
	"AND",
	"OR",
	"NOT",
	"'('",
	"')'",
	"STR",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line ./src/parse.y:107

// ----------------------------------------------------------------------------
// ---[ begin of the tail source file ]----------------------------------------

// ---[ Parser ]---------------------------------------------------------------

// the regex machine
var re = map[int]*regexp.Regexp{
	// e.g. sreach, Sreach, SREACH
	SREACH: regexp.MustCompile(`^[sS][rR][eE][aA][cC][hH]`),
	// e.g. and, AND, And, &&
	AND: regexp.MustCompile(`^([aA][nN][dD]|&&)`),
	// e.g. or, OR, Or, ||
	OR: regexp.MustCompile(`^([oO][rR]|\|\|)`),
	// e.g. not NOT, Not, !
	NOT: regexp.MustCompile(`^([nN][oO][tT]|!)`),
	// e.g. "PETER", "\"", '\'', '"'
	STR: regexp.MustCompile(`^("(\\"|[^"])*"|'(\\'|[^'])*')`),
}

// the struct of input (member input is the string of its input)
// poos hold its postion
type GoLex struct {
	pos   int
	input []byte
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
	for _, v := range org[1 : len(org)-1] {
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
	for l.pos < len(l.input) {
		// get the next Rune(part of string) and its length
		r, n := utf8.DecodeRune(l.input[l.pos:])
		// if it is a space, then skip it.
		if unicode.IsSpace(r) {
			l.pos += n
			continue
		}
		switch {
		case re[SREACH].Match(l.input[l.pos:]):
			l.pos += len(re[SREACH].Find(l.input[l.pos:]))
			return SREACH

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
			log.Println("can't match", "\""+string(l.input[l.pos:])+"\"")
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
func getAST(s string) *typeAst {
	if s == "" {
		return nil
	}
	yyParse(&GoLex{input: []byte(s)})
	// [WARNING] ast_result is a global variable
	return ast_result
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 20

var yyAct = [...]int{

	5, 8, 10, 6, 9, 14, 12, 13, 20, 19,
	11, 15, 16, 3, 4, 18, 17, 1, 7, 2,
}
var yyPact = [...]int{

	9, -1000, -1000, -7, -1000, -1000, 3, 0, -4, -1000,
	-7, -7, -7, -1000, -7, -1, -1000, -1000, -2, -1000,
	-1000,
}
var yyPgo = [...]int{

	0, 19, 0, 3, 18, 17,
}
var yyR1 = [...]int{

	0, 5, 1, 1, 2, 2, 3, 3, 4, 4,
	4, 4,
}
var yyR2 = [...]int{

	0, 1, 2, 1, 3, 1, 3, 1, 2, 1,
	4, 3,
}
var yyChk = [...]int{

	-1000, -5, -1, 4, 5, -2, -3, -4, 8, 11,
	9, 7, 6, 11, 9, -2, -2, -3, -2, 10,
	10,
}
var yyDef = [...]int{

	0, -2, 1, 0, 3, 2, 5, 7, 0, 9,
	0, 0, 0, 8, 0, 0, 4, 6, 0, 11,
	10,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	9, 10,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 11,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
//line ./src/parse.y:60
		{
			ast_result = yyDollar[1].Ast
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
//line ./src/parse.y:64
		{
			yyVAL.Ast = &typeAst{"sreach", yyDollar[2].List}
		}
	case 3:
		yyDollar = yyS[yypt-1 : yypt+1]
//line ./src/parse.y:67
		{
			yyVAL.Ast = &typeAst{"list", nil}
		}
	case 4:
		yyDollar = yyS[yypt-3 : yypt+1]
//line ./src/parse.y:71
		{
			temp := append(*yyDollar[3].List, yyDollar[1].Expr)
			yyVAL.List = &temp
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
//line ./src/parse.y:75
		{
			temp := make(typeList, 0)
			temp = append(temp, yyDollar[1].Expr)
			yyVAL.List = &temp
		}
	case 6:
		yyDollar = yyS[yypt-3 : yypt+1]
//line ./src/parse.y:82
		{
			temp := append(*yyDollar[3].Expr, yyDollar[1].Atom)
			yyVAL.Expr = &temp
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
//line ./src/parse.y:86
		{
			temp := make(typeExpr, 0)
			temp = append(temp, yyDollar[1].Atom)
			yyVAL.Expr = &temp
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
//line ./src/parse.y:93
		{
			yyVAL.Atom = &typeAtom{true, yyDollar[2].Str}
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
//line ./src/parse.y:96
		{
			yyVAL.Atom = &typeAtom{false, yyDollar[1].Str}
		}
	case 10:
		yyDollar = yyS[yypt-4 : yypt+1]
//line ./src/parse.y:99
		{
			yyVAL.Atom = &typeAtom{true, yyDollar[3].List}
		}
	case 11:
		yyDollar = yyS[yypt-3 : yypt+1]
//line ./src/parse.y:102
		{
			yyVAL.Atom = &typeAtom{false, yyDollar[2].List}
		}
	}
	goto yystack /* stack new state and value */
}
