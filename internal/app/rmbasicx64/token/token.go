package token

import (
	"encoding/json"
	"log"
)

// Token defines the actual token generated by the scanner
type Token struct {
	TokenType string
	Literal   string
	Index     int
}

// Token types are defined here
const (
	// Single-character tokens
	LeftParen        = "("
	RightParen       = ")"
	Comma            = ","
	Dot              = "."
	Minus            = "-"
	Plus             = "+"
	Colon            = ":"
	Semicolon        = ";"
	ForwardSlash     = "/"
	BackSlash        = "\\"
	Star             = "*"
	Exponential      = "^"
	LessThan         = "<"
	GreaterThan      = ">"
	Exclamation      = "!"
	Hash             = "#"
	Tilde            = "~"
	LeftSquareBrace  = "["
	RightSquareBrace = "]"
	LeftCurlyBrace   = "{"
	RightCurlyBrace  = "}"

	// Two-character tokens
	Assign              = ":="
	Inequality1         = "<>"
	Inequality2         = "><"
	LessThanEqualTo1    = "<="
	LessThanEqualTo2    = "=<"
	GreaterThanEqualTo1 = ">="
	GreaterThanEqualTo2 = "=>"
	InterestinglyEqual  = "=="
	Equal               = "="

	// Literals
	StringLiteral     = "StringLiteral"
	NumericLiteral    = "NumericLiteral"
	HexLiteral        = "HexLiteral"
	IdentifierLiteral = "IdentifierLiteral"
	Comment           = "Comment"
	Illegal           = "Illegal"
	NewLine           = "\n"

	// Keywords
	ABS        = "ABS"
	AND        = "AND"
	AREA       = "AREA"
	ASC        = "ASC"
	ASK        = "ASK"
	ATN        = "ATN"
	AUTO       = "AUTO"
	BLOCK      = "BLOCK"
	COPY       = "COPY"
	READ       = "READ"
	WRITE      = "WRITE"
	BORDER     = "BORDER"
	BOUNDS     = "BOUNDS"
	BRUSH      = "BRUSH"
	BUTTONS    = "BUTTONS"
	BYE        = "BYE"
	CHAROVER   = "CHAROVER"
	CHARSET    = "CHARSET"
	CHDIR      = "CHDIR"
	CHRstr     = "CHR"
	CIRCLE     = "CIRCLE"
	CLEAR      = "CLEAR"
	CLG        = "CLG"
	CLL        = "CLL"
	CLOSE      = "CLOSE"
	CLS        = "CLS"
	COLOUR     = "COLOUR"
	CONTINUE   = "CONTINUE"
	COS        = "COS"
	CREATE     = "CREATE"
	MOVE       = "MOVE"
	CURPOS     = "CURPOS"
	CURSOR     = "CURSOR"
	DATA       = "DATA"
	DATE       = "DATE"
	DATEstr    = "DATE$"
	DEFINED    = "DEFINED"
	DEG        = "DEG"
	DELETE     = "DELETE"
	DIM        = "DIM"
	DIR        = "DIR"
	DRAWING    = "DRAWING"
	EDIT       = "EDIT"
	END        = "END"
	ENVELOPE   = "ENVELOPE"
	ERASE      = "ERASE"
	ERL        = "ERL"
	ERR        = "ERR"
	ERRstr     = "ERR$"
	EXP        = "EXP"
	FALSE      = "FALSE"
	FKEY       = "FKEY"
	FLOOD      = "FLOOD"
	EDGE       = "EDGE"
	FLUSH      = "FLUSH"
	FOR        = "FOR"
	NEXT       = "NEXT"
	FREE       = "FREE"
	FSPACE     = "FSPACE"
	FUNCTION   = "FUNCTION"
	ENDFUN     = "ENDFUN"
	GET        = "GET"
	GETstr     = "GET$"
	GLOBAL     = "GLOBAL"
	GOSUB      = "GOSUB"
	SUBROUTINE = "SUBROUTINE"
	GOTO       = "GOTO"
	HEXstr     = "HEX$"
	HOLD       = "HOLD"
	HOME       = "HOME"
	IF         = "IF"
	THEN       = "THEN"
	ELSE       = "ELSE"
	INPUT      = "INPUT"
	INSTR      = "INSTR"
	INT        = "INT"
	JOYSTICK   = "JOYSTICK"
	JOYX       = "JOYX"
	JOYY       = "JOYY"
	KEYREP     = "KEYREP"
	LEAVE      = "LEAVE"
	LEFTstr    = "LEFT$"
	LEN        = "LEN"
	LET        = "LET"
	LINE       = "LINE"
	LIST       = "LIST"
	LN         = "LN"
	LOAD       = "LOAD"
	LOADGO     = "LOADGO"
	LOG        = "LOG"
	LOOKUP     = "LOOKUP"
	LVAR       = "LVAR"
	MEM        = "MEM"
	MERGE      = "MERGE"
	MERGEGO    = "MEREGO"
	MIDstr     = "MID$"
	MIX        = "MIX"
	MKDIR      = "MKDIR"
	MOD        = "MOD"
	MODE       = "MODE"
	MOUSE      = "MOUSE"
	NEW        = "NEW"
	NOISE      = "NOISE"
	NOT        = "NOT"
	NOTE       = "NOTE"
	ON         = "ON"
	BREAK      = "BREAK"
	EOF        = "EOF"
	ERROR      = "ERROR"
	OPEN       = "OPEN"
	OR         = "OR"
	ORIGIN     = "ORIGIN"
	OVER       = "OVER"
	PAPER      = "PAPER"
	PATHstr    = "PATH$"
	PATTERN    = "PATTERN"
	PEN        = "PEN"
	PI         = "PI"
	PITCH      = "PITCH"
	PLOT       = "PLOT"
	DIRECTION  = "DIRECTION"
	FONT       = "FONT"
	CHAR       = "CHAR"
	SIZE       = "SIZE"
	POINTS     = "POINTS"
	POS        = "POS"
	POSX       = "POSX"
	POSY       = "POSY"
	PRINT      = "PRINT"
	PROCEDURE  = "PROCEDURE"
	ENDPROC    = "ENDPROC"
	PROCS      = "PROCS"
	PSAVE      = "PSAVE"
	PUT        = "PUT"
	QUEUE      = "QUEUE"
	RAD        = "RAD"
	REM        = "REM"
	RENAME     = "RENAME"
	TO         = "TO"
	RENUMBER   = "RENUMBER"
	REPEAT     = "REPEAT"
	UNTIL      = "UNTIL"
	RESTORE    = "RESTORE"
	RESULT     = "RESULT"
	RESUME     = "RESUME"
	RETURN     = "RETURN"
	RIGHTstr   = "RIGHT$"
	RMDIR      = "RMDIR"
	RND        = "RND"
	RPOINT     = "RPOINT"
	RUN        = "RUN"
	SAVE       = "SAVE"
	SGN        = "SGN"
	SIN        = "SIN"
	SLICE      = "SLICE"
	SOUND      = "SOUND"
	SPC        = "SPC"
	SQR        = "SQR"
	STOP       = "STOP"
	STRINGstr  = "STRING$"
	STRstr     = "STR$"
	STYLE      = "STYLE"
	TAB        = "TAB"
	TAN        = "TAN"
	TIME       = "TIME"
	TIMEstr    = "TIME$"
	TONE       = "TONE"
	TRACE      = "TRACE"
	TRUE       = "TRUE"
	UNDERLINE  = "UNDERLINE"
	VAL        = "VAL"
	VERSION    = "VERSION"
	VOICE      = "VOICE"
	WARN       = "WARN"
	WIDTH      = "WIDTH"
	WRITING    = "WRITING"
	XOR        = "XOR"
	SET        = "SET"
	STEP       = "STEP"
	CONFIG     = "CONFIG"
	BOOT       = "BOOT"
	FETCH      = "FETCH"
	WRITEBLOCK = "WRITEBLOCK"
	SQUASH     = "SQUASH"
	CLEARBLOCK = "CLEARBLOCK"
	DELBLOCK   = "DELBLOCK"
)

// IsKeyword returns true if a TokenType represents a keyword
func IsKeyword(testString string) bool {
	keywords := []string{
		ABS,
		AND,
		AREA,
		ASC,
		ASK,
		ATN,
		AUTO,
		BLOCK,
		COPY,
		READ,
		WRITE,
		BORDER,
		BOUNDS,
		BRUSH,
		BUTTONS,
		BYE,
		CHAROVER,
		CHARSET,
		CHDIR,
		CHRstr,
		CIRCLE,
		CLEAR,
		CLG,
		CLL,
		CLOSE,
		CLS,
		COLOUR,
		CONTINUE,
		COS,
		CREATE,
		MOVE,
		CURPOS,
		CURSOR,
		DATA,
		DATE,
		DATEstr,
		DEFINED,
		DEG,
		DELETE,
		DIM,
		DIR,
		DRAWING,
		EDIT,
		END,
		ENVELOPE,
		ERASE,
		ERL,
		ERR,
		ERRstr,
		EXP,
		FALSE,
		FKEY,
		FLOOD,
		EDGE,
		FLUSH,
		FOR,
		NEXT,
		FREE,
		FSPACE,
		FUNCTION,
		ENDFUN,
		GET,
		GETstr,
		GLOBAL,
		GOSUB,
		SUBROUTINE,
		GOTO,
		HEXstr,
		HOLD,
		HOME,
		IF,
		THEN,
		ELSE,
		INPUT,
		INSTR,
		INT,
		JOYSTICK,
		JOYX,
		JOYY,
		KEYREP,
		LEAVE,
		LEFTstr,
		LEN,
		LET,
		LINE,
		LIST,
		LN,
		LOAD,
		LOADGO,
		LOG,
		LOOKUP,
		LVAR,
		MEM,
		MERGE,
		MERGEGO,
		MIDstr,
		MIX,
		MKDIR,
		MOD,
		MODE,
		MOUSE,
		NEW,
		NOISE,
		NOT,
		NOTE,
		ON,
		BREAK,
		EOF,
		ERROR,
		OPEN,
		OR,
		ORIGIN,
		OVER,
		PAPER,
		PATHstr,
		PATTERN,
		PEN,
		PI,
		PITCH,
		PLOT,
		DIRECTION,
		FONT,
		CHAR,
		SIZE,
		POINTS,
		POS,
		POSX,
		POSY,
		PRINT,
		PROCEDURE,
		ENDPROC,
		PROCS,
		PSAVE,
		PUT,
		QUEUE,
		RAD,
		REM,
		RENAME,
		TO,
		RENUMBER,
		REPEAT,
		UNTIL,
		RESTORE,
		RESULT,
		RESUME,
		RETURN,
		RIGHTstr,
		RMDIR,
		RND,
		RPOINT,
		RUN,
		SAVE,
		SGN,
		SIN,
		SLICE,
		SOUND,
		SPC,
		SQR,
		STOP,
		STRINGstr,
		STRstr,
		STYLE,
		TAB,
		TAN,
		TIME,
		TIMEstr,
		TONE,
		TRACE,
		TRUE,
		UNDERLINE,
		VAL,
		VERSION,
		VOICE,
		WARN,
		WIDTH,
		WRITING,
		XOR,
		SET,
		STEP,
		CONFIG,
		BOOT,
		FETCH,
		WRITEBLOCK,
		SQUASH,
		CLEARBLOCK,
		DELBLOCK,
	}
	for _, keyword := range keywords {
		if testString == keyword {
			return true
		}
	}
	return false
}

// IsOperator receives a token and returns true if the token represents an operator
// otherwise false
func IsOperator(t Token) bool {
	operators := []string{
		Minus,
		Plus,
		ForwardSlash,
		Star,
		Exponential,
		BackSlash,
		Equal,
		InterestinglyEqual,
		LessThan,
		GreaterThan,
		LessThanEqualTo1,
		LessThanEqualTo2,
		GreaterThanEqualTo1,
		GreaterThanEqualTo2,
		Inequality1,
		Inequality2,
		AND,
		OR,
		XOR,
		NOT,
	}
	for _, op := range operators {
		if op == t.TokenType {
			return true
		}
	}
	return false
}

// PrintToken prints a token in the console
func PrintToken(thisToken Token) {
	out, err := json.Marshal(thisToken)
	if err != nil {
		panic(err)
	}
	log.Println(string(out))
}
