package main

import "encoding/json"

// Token defines the actual token generated by the scanner
type Token struct {
	TokenType int
	Literal   string
}

// Token types are defined here
const (
	// Single-character tokens
	LeftParen = iota
	RightParen
	Comma
	Dot
	Minus
	Plus
	Colon
	Semicolon
	ForwardSlash
	BackSlash
	Star
	Exponential
	LessThan
	GreaterThan

	// Two-character tokens
	Assign
	IntegerDivision
	Inequality1
	Inequality2
	LessThanEqualTo1
	LessThanEqualTo2
	GreaterThanEqualTo1
	GreaterThanEqualTo2
	InterestinglyEqual
	Equal

	// Literals
	StringLiteral
	NumericalLiteral
	IdentifierLiteral
	Comment
	Illegal
	EndOfLine

	// Keywords
	ABS
	AND
	AREA
	ASC
	ASK
	ATN
	AUTO
	BLOCK
	COPY
	READ
	WRITE
	BORDER
	BOUNDS
	BRUSH
	BUTTONS
	BYE
	CHAROVER
	CHARSET
	CHDIR
	CHRstr
	CIRCLE
	CLEAR
	CLG
	CLL
	CLOSE
	CLS
	COLOUR
	CONTINUE
	COS
	CREATE
	MOVE
	CURPOS
	CURSOR
	DATA
	DATE
	DATEstr
	DEFINED
	DEG
	DELETE
	DIM
	DIR
	DRAWING
	EDIT
	END
	ENVELOPE
	ERASE
	ERL
	ERR
	ERRstr
	EXP
	FALSE
	FKEY
	FLOOD
	EDGE
	FLUSH
	FOR
	NEXT
	FREE
	FSPACE
	FUNCTION
	ENDFUN
	GET
	GETstr
	GLOBAL
	GOSUB
	SUBROUTINE
	GOTO
	HEXstr
	HOLD
	HOME
	IF
	THEN
	ELSE
	INPUT
	INSTR
	INT
	JOYSTICK
	JOYX
	JOYY
	KEYREP
	LEAVE
	LEFTstr
	LEN
	LET
	LINE
	LIST
	LN
	LOAD
	LOADGO
	LOG
	LOOKUP
	LVAR
	MEM
	MERGE
	MERGEGO
	MIDstr
	MIX
	MKDIR
	MOD
	MODE
	MOUSE
	NEW
	NOISE
	NOT
	NOTE
	ON
	BREAK
	EOF
	ERROR
	OPEN
	OR
	ORIGIN
	OVER
	PAPER
	PATHstr
	PATTERN
	PEN
	PI
	PITCH
	PLOT
	DIRECTION
	FONT
	CHAR
	SIZE
	POINTS
	POS
	POSX
	POSY
	PRINT
	PROCEDURE
	ENDPROC
	PROCS
	PSAVE
	PUT
	QUEUE
	RAD
	REM
	RENAME
	TO
	RENUMBER
	REPEAT
	UNTIL
	RESTORE
	RESULT
	RESUME
	RETURN
	RIGHTstr
	RMDIR
	RND
	RPOINT
	RUN
	SAVE
	SGN
	SIN
	SLICE
	SOUND
	SPC
	SQR
	STOP
	STRINGstr
	STRstr
	STYLE
	TAB
	TAN
	TIME
	TIMEstr
	TONE
	TRACE
	TRUE
	UNDERLINE
	VAL
	VERSION
	VOICE
	WARN
	WIDTH
	WRITING
	XOR
	SET
)

// keywordsMap returns a map of keyword symbols to token ids
func keywordMap() map[string]int {
	return map[string]int{
		"ABS":        ABS,
		"AND":        AND,
		"AREA":       AREA,
		"ASC":        ASC,
		"ASK":        ASK,
		"ATN":        ATN,
		"AUTO":       AUTO,
		"BLOCK":      BLOCK,
		"COPY":       COPY,
		"READ":       READ,
		"WRITE":      WRITE,
		"BORDER":     BORDER,
		"BOUNDS":     BOUNDS,
		"BRUSH":      BRUSH,
		"BUTTONS":    BUTTONS,
		"BYE":        BYE,
		"CHAROVER":   CHAROVER,
		"CHARSET":    CHARSET,
		"CHDIR":      CHDIR,
		"CHR$":       CHRstr,
		"CIRCLE":     CIRCLE,
		"CLEAR":      CLEAR,
		"CLG":        CLG,
		"CLL":        CLL,
		"CLOSE":      CLOSE,
		"CLS":        CLS,
		"COLOUR":     COLOUR,
		"CONTINUE":   CONTINUE,
		"COS":        COS,
		"CREATE":     CREATE,
		"MOVE":       MOVE,
		"CURPOS":     CURPOS,
		"CURSOR":     CURSOR,
		"DATA":       DATA,
		"DATE":       DATE,
		"DATE$":      DATEstr,
		"DEFINED":    DEFINED,
		"DEG":        DEG,
		"DELETE":     DELETE,
		"DIM":        DIM,
		"DIR":        DIR,
		"DRAWING":    DRAWING,
		"EDIT":       EDIT,
		"END":        END,
		"ENVELOPE":   ENVELOPE,
		"ERASE":      ERASE,
		"ERL":        ERL,
		"ERR":        ERR,
		"ERR$":       ERRstr,
		"EXP":        EXP,
		"FALSE":      FALSE,
		"FKEY":       FKEY,
		"FLOOD":      FLOOD,
		"EDGE":       EDGE,
		"FLUSH":      FLUSH,
		"FOR":        FOR,
		"NEXT":       NEXT,
		"FREE":       FREE,
		"FSPACE":     FSPACE,
		"FUNCTION":   FUNCTION,
		"ENDFUN":     ENDFUN,
		"GET":        GET,
		"GET$":       GETstr,
		"GLOBAL":     GLOBAL,
		"GOSUB":      GOSUB,
		"SUBROUTINE": SUBROUTINE,
		"GOTO":       GOTO,
		"HEX$":       HEXstr,
		"HOLD":       HOLD,
		"HOME":       HOME,
		"IF":         IF,
		"THEN":       THEN,
		"ELSE":       ELSE,
		"INPUT":      INPUT,
		"INSTR":      INSTR,
		"INT":        INT,
		"JOYSTICK":   JOYSTICK,
		"JOYX":       JOYX,
		"JOYY":       JOYY,
		"KEYREP":     KEYREP,
		"LEAVE":      LEAVE,
		"LEFT$":      LEFTstr,
		"LEN":        LEN,
		"LET":        LET,
		"LINE":       LINE,
		"LIST":       LIST,
		"LN":         LN,
		"LOAD":       LOAD,
		"LOADGO":     LOADGO,
		"LOG":        LOG,
		"LOOKUP":     LOOKUP,
		"LVAR":       LVAR,
		"MEM":        MEM,
		"MERGE":      MERGE,
		"MERGEGO":    MERGEGO,
		"MID$":       MIDstr,
		"MIX":        MIX,
		"MKDIR":      MKDIR,
		"MOD":        MOD,
		"MODE":       MODE,
		"MOUSE":      MOUSE,
		"NEW":        NEW,
		"NOISE":      NOISE,
		"NOT":        NOT,
		"NOTE":       NOTE,
		"ON":         ON,
		"BREAK":      BREAK,
		"EOF":        EOF,
		"ERROR":      ERROR,
		"OPEN":       OPEN,
		"OR":         OR,
		"ORIGIN":     ORIGIN,
		"OVER":       OVER,
		"PAPER":      PAPER,
		"PATH$":      PATHstr,
		"PATTERN":    PATTERN,
		"PEN":        PEN,
		"PI":         PI,
		"PITCH":      PITCH,
		"PLOT":       PLOT,
		"DIRECTION":  DIRECTION,
		"FONT":       FONT,
		"CHAR":       CHAR,
		"SIZE":       SIZE,
		"POINTS":     POINTS,
		"POS":        POS,
		"POSX":       POSX,
		"POSY":       POSY,
		"PRINT":      PRINT,
		"PROCEDURE":  PROCEDURE,
		"ENDPROC":    ENDPROC,
		"PROCS":      PROCS,
		"PSAVE":      PSAVE,
		"PUT":        PUT,
		"QUEUE":      QUEUE,
		"RAD":        RAD,
		"REM":        REM,
		"RENAME":     RENAME,
		"TO":         TO,
		"RENUMBER":   RENUMBER,
		"REPEAT":     REPEAT,
		"UNTIL":      UNTIL,
		"RESTORE":    RESTORE,
		"RESULT":     RESULT,
		"RESUME":     RESUME,
		"RETURN":     RETURN,
		"RIGHT$":     RIGHTstr,
		"RMDIR":      RMDIR,
		"RND":        RND,
		"RPOINT":     RPOINT,
		"RUN":        RUN,
		"SAVE":       SAVE,
		"SGN":        SGN,
		"SIN":        SIN,
		"SLICE":      SLICE,
		"SOUND":      SOUND,
		"SPC":        SPC,
		"SQR":        SQR,
		"STOP":       STOP,
		"STRING$":    STRINGstr,
		"STR$":       STRstr,
		"STYLE":      STYLE,
		"TAB":        TAB,
		"TAN":        TAN,
		"TIME":       TIME,
		"TIME$":      TIMEstr,
		"TONE":       TONE,
		"TRACE":      TRACE,
		"TRUE":       TRUE,
		"UNDERLINE":  UNDERLINE,
		"VAL":        VAL,
		"VERSION":    VERSION,
		"VOICE":      VOICE,
		"WARN":       WARN,
		"WIDTH":      WIDTH,
		"WRITING":    WRITING,
		"XOR":        XOR,
		"SET":        SET,
	}
}

// PrintToken prints a token in the console
func PrintToken(thisToken Token) {
	out, err := json.Marshal(thisToken)
	if err != nil {
		panic(err)
	}
	logMsg(string(out))
}
