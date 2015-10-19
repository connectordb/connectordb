//line pipeline_generator.y:6
package transforms

import __yyfmt__ "fmt"

//line pipeline_generator.y:6
import (
	//"fmt"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

//line pipeline_generator.y:20
type TransformSymType struct {
	yys        int
	val        TransformFunc
	strVal     string
	stringList []string
	funcList   []TransformFunc
}

const NUMBER = 57346
const BOOL = 57347
const STRING = 57348
const COMPOP = 57349
const THIS = 57350
const OR = 57351
const AND = 57352
const NOT = 57353
const RB = 57354
const LB = 57355
const EOF = 57356
const PIPE = 57357
const RSQUARE = 57358
const LSQUARE = 57359
const COMMA = 57360
const IDENTIFIER = 57361
const HAS = 57362
const IF = 57363
const SET = 57364
const PLUS = 57365
const MINUS = 57366
const MULTIPLY = 57367
const DIVIDE = 57368
const UMINUS = 57369

var TransformToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"NUMBER",
	"BOOL",
	"STRING",
	"COMPOP",
	"THIS",
	"OR",
	"AND",
	"NOT",
	"RB",
	"LB",
	"EOF",
	"PIPE",
	"RSQUARE",
	"LSQUARE",
	"COMMA",
	"IDENTIFIER",
	"HAS",
	"IF",
	"SET",
	"PLUS",
	"MINUS",
	"MULTIPLY",
	"DIVIDE",
	"UMINUS",
}
var TransformStatenames = [...]string{}

const TransformEofCode = 1
const TransformErrCode = 2
const TransformMaxDepth = 200

//line pipeline_generator.y:235

/* Start of lexer, hopefully go will let us do this automatically in the future */

const (
	eof         = 0
	errorString = "<ERROR>"
	eofString   = "<EOF>"
	builtins    = `has|if|set`
	logicals    = `true|false|and|or|not`
	numbers     = `(-)?[0-9]+(\.[0-9]+)?`
	compops     = `<=|>=|<|>|==|!=`
	stringr     = `\"(\\["nrt\\]|.)*?\"|'(\\['nrt\\]|.)*?'`
	pipes       = `:|\||,`
	syms        = `\$|\[|\]|\(|\)`
	idents      = `([a-zA-Z_][a-zA-Z_0-9]*)`
	maths       = `\-|\*|/|\+`
	allregex    = builtins + "|" + logicals + "|" + numbers + "|" + compops + "|" + stringr + "|" + pipes + "|" + syms + "|" + idents + "|" + maths
)

var (
	tokenizer   = regexp.MustCompile(`^(` + allregex + `)`)
	numberRegex = regexp.MustCompile("^" + numbers + "$")
	stringRegex = regexp.MustCompile("^" + stringr + "$")
	identRegex  = regexp.MustCompile("^" + idents + "$")
)

// ParseTransform takes a transform input and returns a function to do the
// transforms.
func ParseTransform(input string) (TransformFunc, error) {
	tl := TransformLex{input: input}

	TransformParse(&tl)

	if tl.errorString == "" {
		return tl.output, nil
	}

	return tl.output, errors.New(tl.errorString)
}

type TransformLex struct {
	input    string
	position int

	errorString string
	output      TransformFunc
}

// Are we at the end of file?
func (t *TransformLex) AtEOF() bool {
	return t.position >= len(t.input)
}

// Return the next string for the lexer
func (l *TransformLex) Next() string {
	var c rune = ' '

	// skip whitespace
	for c == ' ' || c == '\t' {
		if l.AtEOF() {
			return eofString
		}
		c = rune(l.input[l.position])
		l.position += 1
	}

	l.position -= 1

	rest := l.input[l.position:]

	token := tokenizer.FindString(rest)
	l.position += len(token)

	if token == "" {
		return errorString
	}

	return token
}

func (lexer *TransformLex) Lex(lval *TransformSymType) int {

	token := lexer.Next()
	//fmt.Println("token: " + token)
	lval.strVal = token

	switch token {
	case eofString:
		return 0
	case errorString:
		lexer.Error("Error, unknown token")
		return 0
	case "true", "false":
		return BOOL
	case ")":
		return RB
	case "(":
		return LB
	case "[":
		return LSQUARE
	case "]":
		return RSQUARE
	case "$":
		return THIS
	case "has":
		return HAS
	case "and":
		return AND
	case "or":
		return OR
	case "not":
		return NOT
	case ">=", "<=", ">", "<", "==", "!=":
		return COMPOP
	case "if":
		return IF
	case "|", ":":
		return PIPE
	case ",":
		return COMMA
	case "set":
		return SET
	case "-":
		return MINUS
	case "+":
		return PLUS
	case "/":
		return DIVIDE
	case "*":
		return MULTIPLY
	default:
		switch {
		case numberRegex.MatchString(token):
			return NUMBER
		case stringRegex.MatchString(token):
			// unquote token
			strval := token[1 : len(token)-1]

			// replace escape characters
			strval = strings.Replace(strval, "\\n", "\n", -1)
			strval = strings.Replace(strval, "\\r", "\r", -1)
			strval = strings.Replace(strval, "\\t", "\t", -1)
			strval = strings.Replace(strval, "\\\\", "\\", -1)
			strval = strings.Replace(strval, "\\\"", "\"", -1)
			strval = strings.Replace(strval, "\\'", "'", -1)

			lval.strVal = strval
			return STRING
		default:
			return IDENTIFIER
		}
	}
}

func (l *TransformLex) Error(s string) {
	l.errorString = s
}

//line yacctab:1
var TransformExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const TransformNprod = 40
const TransformPrivate = 57344

var TransformTokenNames []string
var TransformStates []string

const TransformLast = 114

var TransformAct = [...]int{

	4, 1, 13, 54, 35, 36, 27, 12, 71, 19,
	20, 21, 39, 22, 11, 33, 34, 37, 18, 69,
	38, 61, 7, 25, 28, 24, 3, 23, 19, 20,
	21, 46, 22, 29, 32, 10, 44, 18, 51, 52,
	70, 49, 50, 6, 24, 5, 23, 48, 14, 43,
	33, 34, 42, 19, 20, 21, 8, 22, 62, 63,
	10, 65, 18, 41, 40, 68, 67, 31, 6, 24,
	5, 23, 72, 14, 19, 20, 21, 60, 22, 61,
	64, 10, 30, 18, 19, 20, 21, 47, 22, 28,
	24, 58, 23, 18, 14, 26, 53, 59, 73, 28,
	24, 26, 23, 56, 14, 55, 66, 57, 2, 45,
	17, 16, 15, 9,
}
var TransformPact = [...]int{

	49, -1000, 8, -1000, 92, 70, 20, 72, -1000, -1000,
	70, 27, -21, -1000, 5, -1000, -1000, -1000, 49, -1000,
	-1000, -1000, -5, 51, 50, 49, 70, 92, 20, 24,
	70, -1000, 80, 80, 80, 5, 5, -1000, 84, 99,
	95, 101, -1000, 72, -1000, 79, -1000, -1000, -8, -21,
	-21, -1000, -1000, -1000, 61, -1000, 41, 68, -1000, 49,
	-1000, 100, 99, 49, -1000, -1000, -1000, 3, 28, -10,
	-1000, 70, 86, -1000,
}
var TransformPgo = [...]int{

	0, 0, 22, 56, 113, 2, 26, 1, 112, 111,
	110, 7, 14, 109, 108, 3,
}
var TransformR1 = [...]int{

	0, 7, 14, 14, 6, 6, 6, 1, 1, 2,
	2, 3, 3, 4, 4, 12, 12, 12, 11, 11,
	11, 11, 5, 5, 5, 5, 8, 8, 8, 9,
	9, 10, 10, 10, 10, 10, 15, 15, 13, 13,
}
var TransformR2 = [...]int{

	0, 1, 1, 3, 1, 2, 1, 1, 3, 1,
	3, 1, 2, 1, 3, 1, 3, 3, 1, 3,
	3, 2, 1, 1, 1, 3, 1, 1, 1, 4,
	1, 9, 6, 4, 3, 4, 1, 3, 1, 3,
}
var TransformChk = [...]int{

	-1000, -7, -14, -6, -1, 21, 19, -2, -3, -4,
	11, -12, -11, -5, 24, -8, -9, -10, 13, 4,
	5, 6, 8, 22, 20, 15, 9, -1, 19, 13,
	10, -3, 7, 23, 24, 25, 26, -5, -7, 17,
	13, 13, -6, -2, 12, -13, -7, -3, -12, -11,
	-11, -5, -5, 12, -15, 6, 8, 6, 12, 18,
	16, 18, 17, 18, 12, -7, 6, -15, -7, 16,
	12, 18, -1, 12,
}
var TransformDef = [...]int{

	0, -2, 1, 2, 4, 0, 6, 7, 9, 11,
	0, 13, 15, 18, 0, 22, 23, 24, 0, 26,
	27, 28, 30, 0, 0, 0, 0, 5, 0, 0,
	0, 12, 0, 0, 0, 0, 0, 21, 0, 0,
	0, 0, 3, 8, 34, 0, 38, 10, 14, 16,
	17, 19, 20, 25, 0, 36, 0, 0, 35, 0,
	29, 0, 0, 0, 33, 39, 37, 0, 0, 0,
	32, 0, 0, 31,
}
var TransformTok1 = [...]int{

	1,
}
var TransformTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27,
}
var TransformTok3 = [...]int{
	0,
}

var TransformErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	TransformDebug        = 0
	TransformErrorVerbose = false
)

type TransformLexer interface {
	Lex(lval *TransformSymType) int
	Error(s string)
}

type TransformParser interface {
	Parse(TransformLexer) int
	Lookahead() int
}

type TransformParserImpl struct {
	lookahead func() int
}

func (p *TransformParserImpl) Lookahead() int {
	return p.lookahead()
}

func TransformNewParser() TransformParser {
	p := &TransformParserImpl{
		lookahead: func() int { return -1 },
	}
	return p
}

const TransformFlag = -1000

func TransformTokname(c int) string {
	if c >= 1 && c-1 < len(TransformToknames) {
		if TransformToknames[c-1] != "" {
			return TransformToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func TransformStatname(s int) string {
	if s >= 0 && s < len(TransformStatenames) {
		if TransformStatenames[s] != "" {
			return TransformStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func TransformErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !TransformErrorVerbose {
		return "syntax error"
	}

	for _, e := range TransformErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + TransformTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := TransformPact[state]
	for tok := TOKSTART; tok-1 < len(TransformToknames); tok++ {
		if n := base + tok; n >= 0 && n < TransformLast && TransformChk[TransformAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if TransformDef[state] == -2 {
		i := 0
		for TransformExca[i] != -1 || TransformExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; TransformExca[i] >= 0; i += 2 {
			tok := TransformExca[i]
			if tok < TOKSTART || TransformExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if TransformExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += TransformTokname(tok)
	}
	return res
}

func Transformlex1(lex TransformLexer, lval *TransformSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = TransformTok1[0]
		goto out
	}
	if char < len(TransformTok1) {
		token = TransformTok1[char]
		goto out
	}
	if char >= TransformPrivate {
		if char < TransformPrivate+len(TransformTok2) {
			token = TransformTok2[char-TransformPrivate]
			goto out
		}
	}
	for i := 0; i < len(TransformTok3); i += 2 {
		token = TransformTok3[i+0]
		if token == char {
			token = TransformTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = TransformTok2[1] /* unknown char */
	}
	if TransformDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", TransformTokname(token), uint(char))
	}
	return char, token
}

func TransformParse(Transformlex TransformLexer) int {
	return TransformNewParser().Parse(Transformlex)
}

func (Transformrcvr *TransformParserImpl) Parse(Transformlex TransformLexer) int {
	var Transformn int
	var Transformlval TransformSymType
	var TransformVAL TransformSymType
	var TransformDollar []TransformSymType
	_ = TransformDollar // silence set and not used
	TransformS := make([]TransformSymType, TransformMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	Transformstate := 0
	Transformchar := -1
	Transformtoken := -1 // Transformchar translated into internal numbering
	Transformrcvr.lookahead = func() int { return Transformchar }
	defer func() {
		// Make sure we report no lookahead when not parsing.
		Transformstate = -1
		Transformchar = -1
		Transformtoken = -1
	}()
	Transformp := -1
	goto Transformstack

ret0:
	return 0

ret1:
	return 1

Transformstack:
	/* put a state and value onto the stack */
	if TransformDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", TransformTokname(Transformtoken), TransformStatname(Transformstate))
	}

	Transformp++
	if Transformp >= len(TransformS) {
		nyys := make([]TransformSymType, len(TransformS)*2)
		copy(nyys, TransformS)
		TransformS = nyys
	}
	TransformS[Transformp] = TransformVAL
	TransformS[Transformp].yys = Transformstate

Transformnewstate:
	Transformn = TransformPact[Transformstate]
	if Transformn <= TransformFlag {
		goto Transformdefault /* simple state */
	}
	if Transformchar < 0 {
		Transformchar, Transformtoken = Transformlex1(Transformlex, &Transformlval)
	}
	Transformn += Transformtoken
	if Transformn < 0 || Transformn >= TransformLast {
		goto Transformdefault
	}
	Transformn = TransformAct[Transformn]
	if TransformChk[Transformn] == Transformtoken { /* valid shift */
		Transformchar = -1
		Transformtoken = -1
		TransformVAL = Transformlval
		Transformstate = Transformn
		if Errflag > 0 {
			Errflag--
		}
		goto Transformstack
	}

Transformdefault:
	/* default state action */
	Transformn = TransformDef[Transformstate]
	if Transformn == -2 {
		if Transformchar < 0 {
			Transformchar, Transformtoken = Transformlex1(Transformlex, &Transformlval)
		}

		/* look through exception table */
		xi := 0
		for {
			if TransformExca[xi+0] == -1 && TransformExca[xi+1] == Transformstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			Transformn = TransformExca[xi+0]
			if Transformn < 0 || Transformn == Transformtoken {
				break
			}
		}
		Transformn = TransformExca[xi+1]
		if Transformn < 0 {
			goto ret0
		}
	}
	if Transformn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			Transformlex.Error(TransformErrorMessage(Transformstate, Transformtoken))
			Nerrs++
			if TransformDebug >= 1 {
				__yyfmt__.Printf("%s", TransformStatname(Transformstate))
				__yyfmt__.Printf(" saw %s\n", TransformTokname(Transformtoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for Transformp >= 0 {
				Transformn = TransformPact[TransformS[Transformp].yys] + TransformErrCode
				if Transformn >= 0 && Transformn < TransformLast {
					Transformstate = TransformAct[Transformn] /* simulate a shift of "error" */
					if TransformChk[Transformstate] == TransformErrCode {
						goto Transformstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if TransformDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", TransformS[Transformp].yys)
				}
				Transformp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if TransformDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", TransformTokname(Transformtoken))
			}
			if Transformtoken == TransformEofCode {
				goto ret1
			}
			Transformchar = -1
			Transformtoken = -1
			goto Transformnewstate /* try again in the same state */
		}
	}

	/* reduction by production Transformn */
	if TransformDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", Transformn, TransformStatname(Transformstate))
	}

	Transformnt := Transformn
	Transformpt := Transformp
	_ = Transformpt // guard against "declared and not used"

	Transformp -= TransformR2[Transformn]
	// Transformp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if Transformp+1 >= len(TransformS) {
		nyys := make([]TransformSymType, len(TransformS)*2)
		copy(nyys, TransformS)
		TransformS = nyys
	}
	TransformVAL = TransformS[Transformp+1]

	/* consult goto table to find next state */
	Transformn = TransformR1[Transformn]
	Transformg := TransformPgo[Transformn]
	Transformj := Transformg + TransformS[Transformp].yys + 1

	if Transformj >= TransformLast {
		Transformstate = TransformAct[Transformg]
	} else {
		Transformstate = TransformAct[Transformj]
		if TransformChk[Transformstate] != -Transformn {
			Transformstate = TransformAct[Transformg]
		}
	}
	// dummy call; replaced with literal code
	switch Transformnt {

	case 1:
		TransformDollar = TransformS[Transformpt-1 : Transformpt+1]
		//line pipeline_generator.y:42
		{
			TransformVAL.val = pipeline(TransformDollar[1].funcList)
			Transformlex.(*TransformLex).output = TransformVAL.val
		}
	case 2:
		TransformDollar = TransformS[Transformpt-1 : Transformpt+1]
		//line pipeline_generator.y:50
		{
			TransformVAL.funcList = []TransformFunc{TransformDollar[1].val}
		}
	case 3:
		TransformDollar = TransformS[Transformpt-3 : Transformpt+1]
		//line pipeline_generator.y:54
		{
			//$$ = append([]TransformFunc{$3}, $1...)
			TransformVAL.funcList = append(TransformDollar[1].funcList, TransformDollar[3].val)
		}
	case 5:
		TransformDollar = TransformS[Transformpt-2 : Transformpt+1]
		//line pipeline_generator.y:64
		{
			TransformVAL.val = pipelineGeneratorIf(TransformDollar[2].val)
		}
	case 6:
		TransformDollar = TransformS[Transformpt-1 : Transformpt+1]
		//line pipeline_generator.y:68
		{
			fun, err := InstantiateRegisteredFunction(TransformDollar[1].strVal)

			if err != nil {
				Transformlex.Error(err.Error())
			}

			TransformVAL.val = fun
		}
	case 8:
		TransformDollar = TransformS[Transformpt-3 : Transformpt+1]
		//line pipeline_generator.y:83
		{
			TransformVAL.val = pipelineGeneratorOr(TransformDollar[1].val, TransformDollar[3].val)
		}
	case 10:
		TransformDollar = TransformS[Transformpt-3 : Transformpt+1]
		//line pipeline_generator.y:91
		{
			TransformVAL.val = pipelineGeneratorAnd(TransformDollar[1].val, TransformDollar[3].val)
		}
	case 12:
		TransformDollar = TransformS[Transformpt-2 : Transformpt+1]
		//line pipeline_generator.y:99
		{
			TransformVAL.val = pipelineGeneratorNot(TransformDollar[2].val)
		}
	case 14:
		TransformDollar = TransformS[Transformpt-3 : Transformpt+1]
		//line pipeline_generator.y:107
		{
			TransformVAL.val = pipelineGeneratorCompare(TransformDollar[1].val, TransformDollar[3].val, TransformDollar[2].strVal)
		}
	case 16:
		TransformDollar = TransformS[Transformpt-3 : Transformpt+1]
		//line pipeline_generator.y:115
		{
			TransformVAL.val = addTransformGenerator(TransformDollar[1].val, TransformDollar[3].val)
		}
	case 17:
		TransformDollar = TransformS[Transformpt-3 : Transformpt+1]
		//line pipeline_generator.y:119
		{
			TransformVAL.val = subtractTransformGenerator(TransformDollar[1].val, TransformDollar[3].val)
		}
	case 19:
		TransformDollar = TransformS[Transformpt-3 : Transformpt+1]
		//line pipeline_generator.y:127
		{
			TransformVAL.val = multiplyTransformGenerator(TransformDollar[1].val, TransformDollar[3].val)
		}
	case 20:
		TransformDollar = TransformS[Transformpt-3 : Transformpt+1]
		//line pipeline_generator.y:131
		{
			TransformVAL.val = divideTransformGenerator(TransformDollar[1].val, TransformDollar[3].val)
		}
	case 21:
		TransformDollar = TransformS[Transformpt-2 : Transformpt+1]
		//line pipeline_generator.y:135
		{
			TransformVAL.val = inverseTransformGenerator(TransformDollar[2].val)
		}
	case 25:
		TransformDollar = TransformS[Transformpt-3 : Transformpt+1]
		//line pipeline_generator.y:145
		{
			TransformVAL.val = TransformDollar[2].val
		}
	case 26:
		TransformDollar = TransformS[Transformpt-1 : Transformpt+1]
		//line pipeline_generator.y:152
		{
			num, err := strconv.ParseFloat(TransformDollar[1].strVal, 64)
			TransformVAL.val = ConstantValueGenerator(num, err)
		}
	case 27:
		TransformDollar = TransformS[Transformpt-1 : Transformpt+1]
		//line pipeline_generator.y:157
		{
			val, err := strconv.ParseBool(TransformDollar[1].strVal)
			TransformVAL.val = ConstantValueGenerator(val, err)
		}
	case 28:
		TransformDollar = TransformS[Transformpt-1 : Transformpt+1]
		//line pipeline_generator.y:162
		{
			TransformVAL.val = ConstantValueGenerator(TransformDollar[1].strVal, nil)
		}
	case 29:
		TransformDollar = TransformS[Transformpt-4 : Transformpt+1]
		//line pipeline_generator.y:169
		{
			TransformVAL.val = pipelineGeneratorGet(TransformDollar[3].stringList)
		}
	case 30:
		TransformDollar = TransformS[Transformpt-1 : Transformpt+1]
		//line pipeline_generator.y:173
		{
			TransformVAL.val = PipelineGeneratorIdentity()
		}
	case 31:
		TransformDollar = TransformS[Transformpt-9 : Transformpt+1]
		//line pipeline_generator.y:180
		{
			TransformVAL.val = pipelineGeneratorSet(TransformDollar[5].stringList, TransformDollar[8].val)
		}
	case 32:
		TransformDollar = TransformS[Transformpt-6 : Transformpt+1]
		//line pipeline_generator.y:184
		{
			TransformVAL.val = pipelineGeneratorSet([]string{}, TransformDollar[5].val)
		}
	case 33:
		TransformDollar = TransformS[Transformpt-4 : Transformpt+1]
		//line pipeline_generator.y:188
		{
			TransformVAL.val = pipelineGeneratorHas(TransformDollar[3].strVal)
		}
	case 34:
		TransformDollar = TransformS[Transformpt-3 : Transformpt+1]
		//line pipeline_generator.y:192
		{
			fun, err := InstantiateRegisteredFunction(TransformDollar[1].strVal)

			if err != nil {
				Transformlex.Error(err.Error())
			}

			TransformVAL.val = fun
		}
	case 35:
		TransformDollar = TransformS[Transformpt-4 : Transformpt+1]
		//line pipeline_generator.y:202
		{
			fun, err := InstantiateRegisteredFunction(TransformDollar[1].strVal, TransformDollar[3].funcList...)

			if err != nil {
				Transformlex.Error(err.Error())
			}

			TransformVAL.val = fun
		}
	case 36:
		TransformDollar = TransformS[Transformpt-1 : Transformpt+1]
		//line pipeline_generator.y:215
		{
			TransformVAL.stringList = []string{TransformDollar[1].strVal}
		}
	case 37:
		TransformDollar = TransformS[Transformpt-3 : Transformpt+1]
		//line pipeline_generator.y:219
		{
			TransformVAL.stringList = append(TransformDollar[1].stringList, TransformDollar[3].strVal)
		}
	case 38:
		TransformDollar = TransformS[Transformpt-1 : Transformpt+1]
		//line pipeline_generator.y:226
		{
			TransformVAL.funcList = []TransformFunc{TransformDollar[1].val}
		}
	case 39:
		TransformDollar = TransformS[Transformpt-3 : Transformpt+1]
		//line pipeline_generator.y:230
		{
			TransformVAL.funcList = append([]TransformFunc{TransformDollar[3].val}, TransformDollar[1].funcList...)
		}
	}
	goto Transformstack /* stack new state and value */
}
