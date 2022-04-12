package examples

import (
	"log"
	"os"
	"path/filepath"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/game/examples/resources/images"
)

// WriteExamples writes the example programs and supporting files to the workspace path
func WriteExamples(workspacePath string) {
	examples := []struct {
		filename string
		program  string
	}{
		{
			filename: "hello.BAS",
			program: `10 REM Pretty much the simplest RM Basic program
20 PRINT "Hello from RM BASICx64`,
		},
		{
			filename: "datatest.BAS",
			program: `10 FOR I% := 0 to 5
20   READ A%
30   PRINT A%
40   RESTORE 60
50 NEXT I%
60 DATA 1,2,3,4`,
		},
		{
			filename: "subroutine.BAS",
			program: `10 PRINT "This is how subroutines work in RM Basic."
20 GOSUB Second_Bus
30 GOSUB First_Bus
40 GOTO 110
50 SUBROUTINE First_Bus
60 PRINT "This is the first bus"
70 RETURN
80 SUBROUTINE Second_Bus
90 PRINT "This is the second bus"
100 RETURN
110 PRINT "Naturally, the second bus goes first lol."`,
		},
		{
			filename: "function.BAS",
			program: `10 PRINT "This is how functions work in RM Basic."
20 PRINT Add_Ten(110)
30 END : REM Function definitions cannot be executed
40 FUNCTION Add_Ten(Number%)
50    RESULT Number% + 10
60 ENDFUN`,
		},
		{
			filename: "procedure.BAS",
			program: `10 PRINT "This is how procedures work in RM Basic."
20 Say_Hello
30 Say_Goodbye
40 Shout_Message "Hellooo!!!", 4, 2
50 Multiply 10, 50 RECEIVE Answer
60 PRINT "The result is: "; Answer
70 END : REM Procedure definitions cannot be executed
80 PROCEDURE Say_Hello
90    PRINT "Hello"
100 ENDPROC
110 PROCEDURE Say_Goodbye
120    PRINT "Goodbye"
130 ENDPROC
140 PROCEDURE Shout_Message Msg$, Font_Size%, Font_Colour%
150   PLOT Msg$, 0, 0 SIZE Font_Size% BRUSH Font_Colour%
160 ENDPROC
170 PROCEDURE Multiply A, B RETURN C
180   PRINT "A: "; A
190   PRINT "B: "; B
200   C = A + B
210 ENDPROC`,
		},
		{
			filename: "hello2.BAS",
			program: `10 REM A slightly more intereting way to say hello
20 SET MODE 40
30 SET BORDER 1 : SET PAPER 5 : CLS
40 PLOT "Greetings!", 45, 150 SIZE 3 BRUSH 0
50 PLOT "Greetings!", 46, 151 SIZE 3 BRUSH 13
60 PLOT "Welcome to", 120, 120 BRUSH 14
70 PLOT "RM BASICx64", 30, 50 SIZE 3, 4 BRUSH 0
80 PLOT "RM BASICx64", 31, 51 SIZE 3, 4`,
		},
		{
			filename: "mouse.BAS",
			program: `10 REM A very, very simple drawing program
20 SET MODE 40
30 PRINT "Click any mouse button to quit"
40 SET MOUSE
50 REPEAT
60   ASK MOUSE Xpos%, Ypos%, Button%
70   POINTS Xpos%, Ypos% BRUSH 13 STYLE 2
80 UNTIL Button% > 0`,
		},
		{
			filename: "meltdown.BAS",
			program: `10 REM Write a flashing yellow warning message
20 REM on a dark grey background in hi-res mode
30 SET MODE 80
40 SET COLOUR 0 TO 8
50 SET COLOUR 3 TO 0
60 SET COLOUR 1 TO 14, 2, 8
70 SET PEN 1 : PRINT "WARNING - Imminent meltdown!"
80 SET PEN 2 : PRINT "Evacuate to at least 100 km distance immediately."
90 SET PEN 3 : PRINT "Good luck and have a nice day."`,
		},
		{
			filename: "mandelbrot.BAS",
			program: `10 REM Render the Mandelbrot set
20 REM Adapted from https://rosettacode.org/wiki/Mandelbrot_set#BASIC
30 SET MODE 40 : SET BORDER 1
40 Maxiteration% := 150
50 FOR X0 := -2 TO 2 STEP 0.01
60   FOR Y0 := -1.5 TO 1.5 STEP 0.01
70     X := 0
80     Y := 0
90     Iteration% := 0
100     REPEAT
110       Xtemp := X * X - Y * Y + X0
120       Y := 2 * X * Y + Y0
130       X := Xtemp
140       Iteration% := Iteration% + 1
150    UNTIL (X * X + Y * Y > (2 * 2)) OR (Iteration% >= Maxiteration%)
160    IF Iteration% <> Maxiteration% THEN C% := Iteration% ELSE C% := 0
170    Xpos% := 50 + ((X0 + 2) * 80)
180    Ypos% := (Y0 + 1.5) * 80
190    Col% := C% / (Maxiteration% / 15)
200    POINTS Xpos%, Ypos% BRUSH Col% STYLE 1
210   NEXT Y0
220 NEXT X0
230 PLOT "The Mandelbrot Set", 90, 2 BRUSH 1
240 PLOT "The Mandelbrot Set", 91, 3 BRUSH 13`,
		},
		{
			filename: "meme.BAS",
			program: `10 REM RM Basic Meme Generator
20 SET MODE 40 : SET BORDER 15
30 FETCH 0, "meme.jpg"
40 WRITEBLOCK 0, 0, 0, -1
50 Meme_Text "LEARN TO PROGRAM", 30, 200
60 Meme_text "LIKE THE ANCIENTS", 25, 0
70 READBLOCK 0, 0, 0; 320, 250
80 KEEP 0, "learnrmbasic.jpg"
90 END
100 PROCEDURE Meme_Text Text$, X%, Y%
110   PLOT Text$, X% - 1, Y% - 1 SIZE 2, 4 BRUSH 0 FONT 1
120   PLOT Text$, X%, Y% - 1 SIZE 2, 4 BRUSH 0 FONT 1
130   PLOT Text$, X% + 1, Y% - 1 SIZE 2, 4 BRUSH 0 FONT 1
140   PLOT Text$, X% + 1, Y% SIZE 2, 4 BRUSH 0 FONT 1
150   PLOT Text$, X% + 1, Y% + 1 SIZE 2, 4 BRUSH 0 FONT 1
160   PLOT Text$, X%, Y% + 1 SIZE 2, 4 BRUSH 0 FONT 1
170   PLOT Text$, X% - 1, Y% + 1 SIZE 2, 4 BRUSH 0 FONT 1
180   PLOT Text$, X% - 1, Y% SIZE 2, 4 BRUSH 0 FONT 1
190   PLOT Text$, X%, Y% SIZE 2, 4 BRUSH 15 FONT 1
200 ENDPROC
`,
		},
		{
			filename: "globals.BAS",
			program: `10 GLOBAL Is_Global_Var%
20 Is_Global_Var% := 10
30 Is_Local_Var% := 20
40 PRINT "Main: Is_Global_Var% = "; Is_Global_Var%
50 PRINT "Main: Is_Local_Var% = "; Is_Local_Var%
60 Test_Globals
70 PRINT "Main: Is_Global_Var% = "; Is_Global_Var%
80 PRINT "Main: Is_Local_Var% = "; Is_Local_Var%
90 END
100 PROCEDURE Test_Globals
120   GLOBAL Is_Global_Var%
130   Is_Local_Var% := 30
140   PRINT "  Test_Globals: Is_Global_Var% = "; Is_Global_Var%
150   Is_Global_Var% := 1000
160   PRINT "  Test_Globals: Is_Global_Var% = "; Is_Global_Var%
170   PRINT "  Test_Globals: Is_Local_Var% = "; Is_Local_Var%
180 ENDPROC
`,
		},
		{
			filename: "procrefs.BAS",
			program: `10 DIM Test%(10)
11 Test%(1) := 1000
15 Do_Stuff Test%()
16 PRINT "Main: Test%(2) = "; Test%(2)
17 PRINT "Get_Stuff(Test%()) = "; Get_Stuff(Test%())
20 END
30 PROCEDURE Do_Stuff Stuff%()
40   PRINT "Do_Stuff: Stuff%(1) = "; Stuff%(1)
45   Stuff%(2) = 2000
50 ENDPROC
60 FUNCTION Get_Stuff(Stuff%())
70   RESULT Stuff%(0) + 1000
80 ENDFUN
`,
		},
		{
			filename: "funcrefs.BAS",
			program: `10 DIM Test%(10)
15 PRINT "Get_Stuff(Test%()) = "; Get_Stuff(Test%())
20 END
60 FUNCTION Get_Stuff(Stuff%())
70   RESULT Stuff%(5) + 1000
80 ENDFUN
`,
		},
		{
			filename: "music1.BAS",
			program: `10 SET SOUND TRUE
20 DIM First%(53), Second%(53)
30 FOR I% = 0 TO 53
40   READ First%(I%)
50   DATA 0,100,15,4,50,15,2,100,15,7,50,15
60   DATA 4,50,15,2,50,15,0,50,15,2,100,15,7,50,15 
70   DATA 4,50,15,2,50,15,0,50,15,2,50,15
80   DATA 4,50,15,5,50,15,4,100,15,2,50,15
90   DATA 0,150,15
100 NEXT I%
110 FOR J% = 0 TO 53
120   READ Second%(J%)
130   DATA 4,100,15,0,50,15,7,150,15,0,100,15
140   DATA 4,50,15,0,2,0,7,150,15,0,100,15,4,50,15 
150   DATA 5,100,15,2,50,15,7,100,15,0,2,0,5,50,15 
160   DATA 0,2,0,4,150,15,0,0,0,0,0,0
170 NEXT J%
180 FOR K% = 0 TO 51 STEP 3
190   SET VOICE 1
200   NOTE PITCH(1, First%(K%)), First%(K% + 1), First%(K% + 2)
210   SET VOICE 2
220   NOTE PITCH (0, Second%(K%)), Second%(K% + 1), Second%(K% + 2)
230 NEXT K%
`,
		},
		{
			filename: "music2.BAS",
			program: `10 SET SOUND TRUE
20 SET ENVELOPE 5 TO 10, 15; 0, 15; 0, 15; 10 
30 SET ENVELOPE 5
40 FOR I% = 1 TO 30
50   NOTE PITCH (RND(3), RND(11))
60 NEXT I%
`,
		},
		{
			filename: "music3.BAS",
			program: `10 SET SOUND TRUE
20 SET SOUND TRUE
30 DIM First%(35), Second%(35)
40 FOR I% = 0 TO 35
50   READ First%(I%)
60   DATA 0,3,4,2,2,3,7,2,4,2,2,2,0,2,2,3,7,2,4,2,2,2,0,2,2,2
70   DATA 4,2,5,2,4,3,2,2,0,4
80 NEXT I%
90 FOR J% = 0 TO 35
100   READ Second%(J%)
110   DATA 4,3,0,2,7,4,0,3,4,1,1,7,4,1,3,4,2,5,3,2,2,7,3
120   DATA 1,1,5,2,1,1,4,4,1,0,1,0,1
130 NEXT J%
140 SET ENVELOPE 1 TO 0, 0; 2, 0; 0, 0; 0
150 SET ENVELOPE 2 TO 5, 15; 10, 15; 0, 15; 70 
160 SET ENVELOPE 3 TO 10, 15; 20, 15; 0, 15; 70 
170 SET ENVELOPE 4 TO 15, 15; 30, 15; 0, 15; 105 
180 FOR K% = 0 TO 34 STEP 2
190   SET VOICE 1
200   NOTE PITCH(1, First%(K%)) ENVELOPE First%(K% + 1)
210   SET VOICE 2
220   NOTE PITCH(0, Second%(K%)) ENVELOPE Second%(K% + 1)
230 NEXT K%
`,
		},
		{
			filename: "lathe.BAS",
			program: `10 REM *****************************
20 REM *                           *
30 REM * Lathe Simulation Program. *
40 REM * 1986-1987 (C) Rob Baines. *
50 REM *          A Level          *
60 REM *      Technical Graphics   *
70 REM *                           *
80 REM *****************************
90 SET MODE 40
100 GLOBAL Ro(), Yo(), T%, F%, T1x%(), T1y%()
110 SET MOUSE
120 DIM Ro(50), Yo(50), T1x%(10), T1%(10)
130 SET PAPER 5
140 SET BORDER 5
150 CLG
160 AREA 80, 40; 270, 40; 270, 200; 80, 200; 80, 40 BRUSH 8
170 AREA 80, 40; 270, 40; 270, 200; 80, 200; 80, 40 BRUSH 0 STYLE 0
180 AREA 70, 50; 260, 50; 260, 210; 70, 210; 70, 50 BRUSH 7
190 AREA 70, 50; 260, 50; 260, 210; 70, 210; 70, 50 BRUSH 0 STYLE 0
200 SET WRITING 4 TO 11, 10; 31, 18
210 SET WRITING 4
220 HOME
230 SET PAPER 7
240 AREA 75, 205; 255, 205; 255, 180; 75, 180; 75, 180 BRUSH 15
250 AREA 75, 205; 255, 205; 255, 180; 75, 180; 75, 180 BRUSH 0 STYLE 0
260 PLOT "Lathe Simulation", 82, 193 BRUSH 0
270 PLOT "Lathe Simulation", 83, 192 BRUSH 2
280 PLOT "by Rob Baines", 120, 182 BRUSH 0
290 PLOT "by Rob Baines", 121, 181 BRUSH 4
300 AREA 75, 60; 255, 60; 255, 170; 75, 170; 75, 60 BRUSH 15
310 AREA 75, 60; 255, 60; 255, 170; 75, 170; 75, 60 BRUSH 0 STYLE 0
320 SET PAPER 15
330 SET PEN 3
340 PRINT "1.) Edit a shape"
350 PRINT "2.) Save data"
360 PRINT "3.) Load data"
370 PRINT "4.) 3D view"
380 PRINT "5.) Plot shape"
390 PRINT "6.) Directory"
400 PRINT "7.) Cut shape"
410 PRINT "8.) Quit"
420 SET PEN 5
430 PRINT "Please select option";
440 SET POINTS STYLE 2
450 ASK MOUSE X%, Y%, B%
460 AREA X%, Y%; X%, Y% - 10; X% + %, Y% - 10; X% + 7, Y% - 20; 10 + X%, Y% - 20; X% + 15, Y% - 20; X% + 9, Y% - 10; X% + 15, Y% - 10; X%, Y% BRUSH 5 OVER FALSE STYLE 0
470 AREA X%, Y%; X%, Y% - 10; X% + %, Y% - 10; X% + 7, Y% - 20; 10 + X%, Y% - 20; X% + 15, Y% - 20; X% + 9, Y% - 10; X% + 15, Y% - 10; X%, Y% BRUSH 5 OVER FALSE STYLE 0
480 IF B% = 1 OR B% = 2 THEN 500
490 GOTO 450
500 IF X% < 75 OR X% > 255 THEN 450
510 Op% := Y% / 10 - 8
520 Op% := - Op% + 8
530 IF Op% = 1 THEN Ed : GOTO 130
540 IF Op% = 2 THEN Savdat : GOTO 130
550 IF Op% = 3 THEN Loadat : GOTO 130
560 IF Op% = 4 THEN View : GOTO 130
570 IF Op% = 5 THEN Plott : GOTO 130
580 IF Op% = 6 THEN Directory : GOTO 130
590 IF Op% = 7 THEN Cut : GOTO 130
600 IF Op% = 8 THEN Quit : END
610 GOTO 450
`,
		},
	}

	for _, example := range examples {
		fullpath := filepath.Join(workspacePath, example.filename)
		file, err := os.Create(fullpath)
		if err != nil {
			log.Printf("Error creating example program %q - %e", example.filename, err)
			continue
		}
		defer file.Close()
		file.WriteString(example.program)
	}

	// Images
	fullpath := filepath.Join(workspacePath, "astronaut.jpg")
	file, err := os.Create(fullpath)
	if err != nil {
		log.Printf("Error creating example program %q - %e", fullpath, err)
	}
	defer file.Close()
	file.Write(images.Astronaut_jpg)

	fullpath = filepath.Join(workspacePath, "meme.jpg")
	file, err = os.Create(fullpath)
	if err != nil {
		log.Printf("Error creating example program %q - %e", fullpath, err)
	}
	defer file.Close()
	file.Write(images.Meme_jpg)
}
