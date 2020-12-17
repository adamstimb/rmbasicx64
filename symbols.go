package main

// Token types/IDs defined here
const (
	// Punctuation
	PnStatementSeparator  int = 1
	PnValueSeparator      int = 2
	PnCoordinateSeperator int = 3
	PnLeftParenthesis     int = 4
	PnRightParenthesis    int = 5
	// Mathematical
	MaAssign              int = 10
	MaExponential         int = 11
	MaAddition            int = 12
	MaSubtraction         int = 13
	MaMultiplication      int = 14
	MaDivision            int = 15
	MaIntegerDivision     int = 16
	MaEquality            int = 17
	MaInequality1         int = 18
	MaInequality2         int = 19
	MaLessThan            int = 20
	MaGreaterThan         int = 21
	MaLessThanEqualTo1    int = 22
	MaLessThanEqualTo2    int = 23
	MaGreaterThanEqualTo1 int = 24
	MaGreaterThanEqualTo2 int = 25
	MaInterestinglyEqual  int = 26
	MaVariableString      int = 27
	MaVariableInteger     int = 28
	MaVariableFloat       int = 29
	// Literals
	LiString  int = 50
	LiInteger int = 51
	LiFloat   int = 52
	LiNumber  int = 53
	// Keywords
	KwABS        int = 100
	KwAND        int = 101
	KwAREA       int = 103
	KwASC        int = 104
	KwASK        int = 105
	KwATN        int = 106
	KwAUTO       int = 107
	KwBLOCK      int = 108
	KwCOPY       int = 109
	KwREAD       int = 110
	KwWRITE      int = 111
	KwBORDER     int = 112
	KwBOUNDS     int = 113
	KwBRUSH      int = 114
	KwBUTTONS    int = 115
	KwBYE        int = 116
	KwCHAROVER   int = 117
	KwCHARSET    int = 118
	KwCHDIR      int = 119
	KwCHRstr     int = 120
	KwCIRCLE     int = 121
	KwCLEAR      int = 122
	KwCLG        int = 123
	KwCLL        int = 124
	KwCLOSE      int = 125
	KwCLS        int = 126
	KwCOLOUR     int = 127
	KwCONTINUE   int = 128
	KwCOS        int = 129
	KwCREATE     int = 130
	KwMOVE       int = 131
	KwCURPOS     int = 132
	KwCURSOR     int = 133
	KwDATA       int = 134
	KwDATE       int = 135
	KwDATEstr    int = 136
	KwDEFINED    int = 137
	KwDEG        int = 138
	KwDELETE     int = 139
	KwDIM        int = 140
	KwDIR        int = 141
	KwDRAWING    int = 142
	KwEDIT       int = 143
	KwEND        int = 144
	KwENVELOPE   int = 145
	KwERASE      int = 146
	KwERL        int = 147
	KwERR        int = 148
	KwERRstr     int = 149
	KwEXP        int = 150
	KwFALSE      int = 151
	KwFKEY       int = 152
	KwFLOOD      int = 153
	KwEDGE       int = 154
	KwFLUSH      int = 155
	KwFOR        int = 156
	KwNEXT       int = 157
	KwFREE       int = 158
	KwFSPACE     int = 159
	KwFUNCTION   int = 160
	KwENDFUN     int = 161
	KwGET        int = 162
	KwGETstr     int = 163
	KwGLOBAL     int = 164
	KwGOSUB      int = 165
	KwSUBROUTINE int = 166
	KwGOTO       int = 167
	KwHEXstr     int = 168
	KwHOLD       int = 169
	KwHOME       int = 170
	KwIF         int = 171
	KwTHEN       int = 172
	KwELSE       int = 173
	KwINPUT      int = 174
	KwINSTR      int = 175
	KwINT        int = 176
	KwJOYSTICK   int = 177
	KwJOYX       int = 178
	KwJOYY       int = 179
	KwKEYREP     int = 180
	KwLEAVE      int = 181
	KwLEFTstr    int = 182
	KwLEN        int = 183
	KwLET        int = 184
	KwLINE       int = 185
	KwLIST       int = 186
	KwLN         int = 187
	KwLOAD       int = 188
	KwLOADGO     int = 189
	KwLOG        int = 190
	KwLOOKUP     int = 191
	KwLVAR       int = 192
	KwMEM        int = 193
	KwMERGE      int = 194
	KwMERGEGO    int = 195
	KwMIDstr     int = 196
	KwMIX        int = 197
	KwMKDIR      int = 198
	KwMOD        int = 199
	KwMODE       int = 200
	KwMOUSE      int = 201
	KwNEW        int = 202
	KwNOISE      int = 203
	KwNOT        int = 204
	KwNOTE       int = 205
	KwON         int = 206
	KwBREAK      int = 207
	KwEOF        int = 208
	KwERROR      int = 209
	KwOPEN       int = 210
	KwOR         int = 211
	KwORIGIN     int = 212
	KwOVER       int = 213
	KwPAPER      int = 214
	KwPATHstr    int = 215
	KwPATTERN    int = 216
	KwPEN        int = 217
	KwPI         int = 218
	KwPITCH      int = 219
	KwPLOT       int = 220
	KwDIRECTION  int = 221
	KwFONT       int = 222
	KwCHAR       int = 223
	KwSIZE       int = 224
	KwPOINTS     int = 225
	KwPOS        int = 226
	KwPOSX       int = 227
	KwPOSY       int = 228
	KwPRINT      int = 229
	KwPROCEDURE  int = 230
	KwENDPROC    int = 231
	KwPROCS      int = 232
	KwPSAVE      int = 233
	KwPUT        int = 234
	KwQUEUE      int = 235
	KwRAD        int = 236
	KwREM        int = 238
	KwRENAME     int = 239
	KwTO         int = 240
	KwRENUMBER   int = 241
	KwREPEAT     int = 242
	KwUNTIL      int = 243
	KwRESTORE    int = 244
	KwRESULT     int = 245
	KwRESUME     int = 246
	KwRETURN     int = 247
	KwRIGHTstr   int = 248
	KwRMDIR      int = 249
	KwRND        int = 250
	KwRPOINT     int = 251
	KwRUN        int = 252
	KwSAVE       int = 253
	KwSGN        int = 254
	KwSIN        int = 255
	KwSLICE      int = 256
	KwSOUND      int = 257
	KwSPC        int = 258
	KwSQR        int = 259
	KwSTOP       int = 260
	KwSTRINGstr  int = 261
	KwSTRstr     int = 262
	KwSTYLE      int = 263
	KwTAB        int = 264
	KwTAN        int = 265
	KwTIME       int = 266
	KwTIMEstr    int = 267
	KwTONE       int = 268
	KwTRACE      int = 269
	KwTRUE       int = 270
	KwUNDERLINE  int = 271
	KwVAL        int = 272
	KwVERSION    int = 273
	KwVOICE      int = 274
	KwWARN       int = 275
	KwWIDTH      int = 276
	KwWRITING    int = 277
	KwXOR        int = 278
	KwSET        int = 279
)

// invertStringIntMap receives a map[string][int], swaps the keys for values and returns a map[int][string]
func invertStringIntMap(mapToInvert map[string]int) map[int]string {
	var newMap map[int]string
	for k, v := range mapToInvert {
		newMap[v] = k
	}
	return newMap
}

// punctuationToTokens returns a map of punctuation symbols to token ids
func punctuationToTokens() map[string]int {
	return map[string]int{
		":": PnStatementSeparator,
		",": PnValueSeparator,
		";": PnCoordinateSeperator,
		"(": PnLeftParenthesis,
		")": PnRightParenthesis,
	}
}

// tokensToPunctuation returns a map of tokens to punctuation symbols
func tokensToPunctuation() map[int]string {
	return invertStringIntMap(punctuationToTokens())
}

// mathematicalToTokens returns a map of mathematical symbols to token ids
func mathematicalToTokens() map[string]int {
	return map[string]int{
		":=": MaAssign,
		"^":  MaExponential,
		"+":  MaAddition,
		"-":  MaSubtraction,
		"*":  MaMultiplication,
		"/":  MaDivision,
		"\\": MaIntegerDivision,
		"=":  MaEquality,
		"<>": MaInequality1,
		"><": MaInequality2,
		"<":  MaLessThan,
		">":  MaGreaterThan,
		"<=": MaLessThanEqualTo1,
		"=<": MaLessThanEqualTo2,
		">=": MaGreaterThanEqualTo1,
		"=>": MaGreaterThanEqualTo2,
		"==": MaInterestinglyEqual,
	}
}

// tokensToMathematical returns a map of tokens to mathematical symbols
func tokensToMathematical() map[int]string {
	return invertStringIntMap(mathematicalToTokens())
}

// keywordsToTokens returns a map of keyword symbols to token ids
func keywordsToTokens() map[string]int {
	return map[string]int{
		"ABS":        KwABS,
		"AND":        KwAND,
		"AREA":       KwAREA,
		"ASC":        KwASC,
		"ASK":        KwASK,
		"ATN":        KwATN,
		"AUTO":       KwAUTO,
		"BLOCK":      KwBLOCK,
		"COPY":       KwCOPY,
		"READ":       KwREAD,
		"WRITE":      KwWRITE,
		"BORDER":     KwBORDER,
		"BOUNDS":     KwBOUNDS,
		"BRUSH":      KwBRUSH,
		"BUTTONS":    KwBUTTONS,
		"BYE":        KwBYE,
		"CHAROVER":   KwCHAROVER,
		"CHARSET":    KwCHARSET,
		"CHDIR":      KwCHDIR,
		"CHR$":       KwCHRstr,
		"CIRCLE":     KwCIRCLE,
		"CLEAR":      KwCLEAR,
		"CLG":        KwCLG,
		"CLL":        KwCLL,
		"CLOSE":      KwCLOSE,
		"CLS":        KwCLS,
		"COLOUR":     KwCOLOUR,
		"CONTINUE":   KwCONTINUE,
		"COS":        KwCOS,
		"CREATE":     KwCREATE,
		"MOVE":       KwMOVE,
		"CURPOS":     KwCURPOS,
		"CURSOR":     KwCURSOR,
		"DATA":       KwDATA,
		"DATE":       KwDATE,
		"DATE$":      KwDATEstr,
		"DEFINED":    KwDEFINED,
		"DEG":        KwDEG,
		"DELETE":     KwDELETE,
		"DIM":        KwDIM,
		"DIR":        KwDIR,
		"DRAWING":    KwDRAWING,
		"EDIT":       KwEDIT,
		"END":        KwEND,
		"ENVELOPE":   KwENVELOPE,
		"ERASE":      KwERASE,
		"ERL":        KwERL,
		"ERR":        KwERR,
		"ERR$":       KwERRstr,
		"EXP":        KwEXP,
		"FALSE":      KwFALSE,
		"FKEY":       KwFKEY,
		"FLOOD":      KwFLOOD,
		"EDGE":       KwEDGE,
		"FLUSH":      KwFLUSH,
		"FOR":        KwFOR,
		"NEXT":       KwNEXT,
		"FREE":       KwFREE,
		"FSPACE":     KwFSPACE,
		"FUNCTION":   KwFUNCTION,
		"ENDFUN":     KwENDFUN,
		"GET":        KwGET,
		"GET$":       KwGETstr,
		"GLOBAL":     KwGLOBAL,
		"GOSUB":      KwGOSUB,
		"SUBROUTINE": KwSUBROUTINE,
		"GOTO":       KwGOTO,
		"HEX$":       KwHEXstr,
		"HOLD":       KwHOLD,
		"HOME":       KwHOME,
		"IF":         KwIF,
		"THEN":       KwTHEN,
		"ELSE":       KwELSE,
		"INPUT":      KwINPUT,
		"INSTR":      KwINSTR,
		"INT":        KwINT,
		"JOYSTICK":   KwJOYSTICK,
		"JOYX":       KwJOYX,
		"JOYY":       KwJOYY,
		"KEYREP":     KwKEYREP,
		"LEAVE":      KwLEAVE,
		"LEFT$":      KwLEFTstr,
		"LEN":        KwLEN,
		"LET":        KwLET,
		"LINE":       KwLINE,
		"LIST":       KwLIST,
		"LN":         KwLN,
		"LOAD":       KwLOAD,
		"LOADGO":     KwLOADGO,
		"LOG":        KwLOG,
		"LOOKUP":     KwLOOKUP,
		"LVAR":       KwLVAR,
		"MEM":        KwMEM,
		"MERGE":      KwMERGE,
		"MERGEGO":    KwMERGEGO,
		"MID$":       KwMIDstr,
		"MIX":        KwMIX,
		"MKDIR":      KwMKDIR,
		"MOD":        KwMOD,
		"MODE":       KwMODE,
		"MOUSE":      KwMOUSE,
		"NEW":        KwNEW,
		"NOISE":      KwNOISE,
		"NOT":        KwNOT,
		"NOTE":       KwNOTE,
		"ON":         KwON,
		"BREAK":      KwBREAK,
		"EOF":        KwEOF,
		"ERROR":      KwERROR,
		"OPEN":       KwOPEN,
		"OR":         KwOR,
		"ORIGIN":     KwORIGIN,
		"OVER":       KwOVER,
		"PAPER":      KwPAPER,
		"PATH$":      KwPATHstr,
		"PATTERN":    KwPATTERN,
		"PEN":        KwPEN,
		"PI":         KwPI,
		"PITCH":      KwPITCH,
		"PLOT":       KwPLOT,
		"DIRECTION":  KwDIRECTION,
		"FONT":       KwFONT,
		"CHAR":       KwCHAR,
		"SIZE":       KwSIZE,
		"POINTS":     KwPOINTS,
		"POS":        KwPOS,
		"POSX":       KwPOSX,
		"POSY":       KwPOSY,
		"PRINT":      KwPRINT,
		"PROCEDURE":  KwPROCEDURE,
		"ENDPROC":    KwENDPROC,
		"PROCS":      KwPROCS,
		"PSAVE":      KwPSAVE,
		"PUT":        KwPUT,
		"QUEUE":      KwQUEUE,
		"RAD":        KwRAD,
		"REM":        KwREM,
		"RENAME":     KwRENAME,
		"TO":         KwTO,
		"RENUMBER":   KwRENUMBER,
		"REPEAT":     KwREPEAT,
		"UNTIL":      KwUNTIL,
		"RESTORE":    KwRESTORE,
		"RESULT":     KwRESULT,
		"RESUME":     KwRESUME,
		"RETURN":     KwRETURN,
		"RIGHT$":     KwRIGHTstr,
		"RMDIR":      KwRMDIR,
		"RND":        KwRND,
		"RPOINT":     KwRPOINT,
		"RUN":        KwRUN,
		"SAVE":       KwSAVE,
		"SGN":        KwSGN,
		"SIN":        KwSIN,
		"SLICE":      KwSLICE,
		"SOUND":      KwSOUND,
		"SPC":        KwSPC,
		"SQR":        KwSQR,
		"STOP":       KwSTOP,
		"STRING$":    KwSTRINGstr,
		"STR$":       KwSTRstr,
		"STYLE":      KwSTYLE,
		"TAB":        KwTAB,
		"TAN":        KwTAN,
		"TIME":       KwTIME,
		"TIME$":      KwTIMEstr,
		"TONE":       KwTONE,
		"TRACE":      KwTRACE,
		"TRUE":       KwTRUE,
		"UNDERLINE":  KwUNDERLINE,
		"VAL":        KwVAL,
		"VERSION":    KwVERSION,
		"VOICE":      KwVOICE,
		"WARN":       KwWARN,
		"WIDTH":      KwWIDTH,
		"WRITING":    KwWRITING,
		"XOR":        KwXOR,
		"SET":        KwSET,
	}
}

// tokensToKeywords returns a map of tokens to keywords
func tokensToKeywords() map[int]string {
	return invertStringIntMap(keywordsToTokens())
}
