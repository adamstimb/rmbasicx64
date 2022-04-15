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
620 PROCEDURE Directory
630   SET WRITING 0
640   SET PAPER 10
650   SET BORDER 10
660   SET BRUSH 0
670   CLG
680   AREA 30, 0; 310, 0; 310, 230; 30, 240; 30, 0 BRUSH 2
690   AREA 30, 0; 310, 0; 310, 230; 30, 240; 30, 0 BRUSH 0 STYLE 0
700   AREA 10, 20; 290, 20; 290, 249; 10, 249; 10, 20 BRUSH 15
710   AREA 10, 20; 290, 20; 290, 249; 10, 249; 10, 20 BRUSH 0 STYLE 0
720   PLOT "Shapes Available.", 40, 0 BRUSH 0 SIZE 2, 2
730   PLOT "Shapes Available.", 41, 1 BRUSH 13 SIZE 2, 2
740   SET WRITING 2 TO 3, 2; 36, 22
750   SET WRITING 2
760   SET PAPER 15 : SET PEN 0
770   DIR "*.3dt"
780   ASK MOUSE X%, Y%, B%
790   IF B% = 0 THEN 780
800 ENDPROC
810 PROCEDURE Quit
820   SET WRITING 0
830   SET PAPER 8 : SET BORDER 8
840   CLG
850   PLOT "Have a nice day!", 35, 230 SIZE 2, 2 BRUSH 0
860   PLOT "Have a nice day!", 36, 229 SIZE 2, 2 BRUSH 14
870   PLOT "Thank you for using this program.", 30, 200 SIZE 1, 2 BRUSH 0
880   PLOT "Thank you for using this program.", 31, 199 SIZE 1, 2 BRUSH 12
890   PLOT "** GOODBYE **", 0, 100 SIZE 3, 3 BRUSH 15
900   Of := 0
910   FOR N := 0 TO 15
920     LINE 60, 116 + N; 247, 166 + N BRUSH (N + Of) AND 15 OVER FALSE
930     LINE 60, 115 - N; 247, 115 - N BRUSH (N + Of) AND 15 OVER FALSE
940   NEXT
950   Of := Of + 1 : IF Of = 16 THEN Of := 0
960   ASK MOUSE X%, Y%, B%
970   IF B% = 0 THEN 910
980 ENDPROC
990 PROCEDURE Ed
1000   GLOBAL Ro(), Yo(), T%
1010   SET PAPER 0 : SET BORDER 0 : CLG : SET WRITING 0 : HOME
1020   SET POINTS STYLE 3
1030   AREA 56, 38; 56, 238; 282, 238; 282, 38; 56, 38 BRUSH 8
1040   AREA 56, 38; 56, 238; 282, 238; 282, 38; 56, 38 BRUSH 0 STYLE 0
1050   FOR N := 46 TO 240 STEP 16
1060     LINE 50, N; 280, N BRUSH 4
1070   NEXT
1080   FOR N := 50 TO 280 STEP 16
1090	 LINE N, 46; N, 240 BRUSH 4
1100   NEXT
1110   Ro(1) := 0 : Yo(1) := 50
1120   Cx := 160 : Cv := 125
1130   Px := Cx + Po(1) : Py := Cy - Yo(1)
1140   Xpos% := Px : Ypos% := Py : Xposn% := Px : Yposn% := Py
1150   T% := 1
1160   Mx := Px : My := Py : Pxm := Px : Pym := Py
1170   LINE Pxm, Pym; Mx, My BRUSH 15 OVER FALSE
1180   LINE Px, Py; Xpos%, Ypos% BRUSH 15 OVER FALSE
1190   Mx := Cx + (-1 * (Xposn% - Cx))
1200   My := Cy - ((Yposn% - Cy))
1210   POINTS Xpos%, Ypos% BRUSH 2 OVER FALSE
1220   ASK MOUSE Xposn%, Yposn%, Button%
1230   IF Button% = 2 AND Xposn% > 160 THEN 1340
1240   IF Button% = 3 THEN 1430
1250   IF Xposn% = Xpos% AND Yposn% = Ypos% THEN 1220
1260   IF Xposn% < 160 THEN 1220
1270   POINTS Xpos%, Ypos% BRUSH 2 OVER FALSE
1280   LINE Px, Py; Xpos%, Ypos% BRUSH 15 OVER FALSE
1290   LINE Pxm, Pym; Mx, My BRUSH 15 OVER FALSE
1300   Xpos% := Xposn% : Ypos% := Yposn%
1310   Mx := Cx + (-1 * (Xposn% - Cx))
1320   My := Cy - (Yposn% - Cy)
1330   GOTO 1170
1340   Px := Xposn% : Py := Yposn%
1350   Pxm := Cx + (-1 * (Xposn% - Cx))
1360   Pym := Cy - (Yposn% - Cy)
1370   Yo(T%) := Yposn% - Cv
1380   Ro(T%) := Xposn% - Cx
1390   T% := T% + 1
1400   IF T% = 50 THEN 1430
1410   FOR P := 0 TO 500 : NEXT
1420   GOTO 1170
1430 ENDPROC
1440 PROCEDURE Savdat
1450   GLOBAL Ro(), Yo(), T%, F$
1460   SET PAPER 12 : SET BORDER 12 : CLG : SET WRITING 0 : CLS
1470   AREA 50, 100; 270, 230; 50, 230; 50, 100 BRUSH 4
1480   AREA 50, 100; 270, 230; 50, 230; 50, 100 BRUSH 0 STYLE 0
1490   AREA 30, 120; 250, 120; 250, 250; 30, 250; 30, 120 BRUSH 15
1500   AREA 30, 120; 250, 120; 250, 250; 30, 250; 30, 120 BRUSH 0 STYLE 0
1510   SET WRITING 2 TO 5, 2; 31, 12
1520   SET WRITING 2
1530   SET PAPER 15
1540   SET PEN 0
1550   CLS
1560   PLOT "Enter filename to be save. (8 Chars)", 10, 75 SIZE 1, 2 BRUSH 0
1570   PLOT "Enter filename to be save. (8 Chars)", 11, 74 SIZE 1, 2 BRUSH 9
1580   HOME : INPUT F%
1590   IF LEN F$  > 8 OR LEN F$ = 0 THEN 1730 : REM Uh oh!
1600   F$ := F$ + ".3dt"
1610   PLOT "Please Wait...", 91, 199 SIZE 1, 2 BRUSH 0
1620   PLOT "Please Wait...", 90, 200 SIZE 1, 2 BRUSH 8
1630   P$ := "Saving:- " + F$
1640   PLOT P$, 70 + ((8 - LEN F$) * 4), 170 SIZE 1, 2 BRUSH 0
1650   PLOT P$, 71 + ((8 - LEN F$) * 4), 169 SIZE 1, 2 BRUSH 2
1660   CREATE #11, F$
1670   PRINT #11, T%
1680   FOR D% := 1 TO T%
1690     PRINT #11, Yo(D%), Ro(D%)
1700   NEXT
1710   CLOSE #11
1720 ENDPROC
1730 CLS : GOTO 1580 : REM Yikes...
1740 PROCEDURE Loaddat
1750   GLOBAL Ro(), Yo(), T%, F$
1760   SET PAPER 12 : SET BORDER 12 : CLG : SET WRITING 0 : CLS
1770   AREA 50, 100; 270, 100; 270, 230; 50, 230; 50, 100 BRUSH 4
1780   AREA 50, 100; 270, 100; 270, 230; 50, 230; 50, 100 BRUSH 0 STYLE 0
1790   AREA 30, 120; 250, 120; 250, 250; 30, 250; 30, 120 BRUSH 15
1800   AREA 30, 120; 250, 120; 250, 250; 30, 250; 30, 120 BRUSH 0 STYLE 0
1810   SET WRITING 2 TO 5, 2; 31, 12
1820   SET WRITING 2
1830   SET PAPER 15
1840   SET PEN 0
1850   CLS
1860   PLOT "Enter filename to be loaded. (8 Chars)", 10, 75 SIZE 1, 2 BRUSH 0
1870   PLOT "Enter filename to be loaded. (8 Chars)", 11, 74 SIZE 1, 2 BRUSH 9
1880   HOME : INPUT F$
1890   IF LEN F$ > 8 OR LEN F$ = 0 THEN 1880
1900   F$ := F$ + ".3dt"
1910   PLOT "Please Wait...", 91, 199 SIZE 1, 2 BRUSH 0
1920   PLOT "Please Wait...", 90, 200 SIZE 1, 2 BRUSH 8
1930   P$ := "Loading:- " + F$
1940   PLOT P$, 70 + ((8 - LEN F$) * 4), 170 SIZE 1, 2 BRUSH 0
1950   PLOT P$, 71 + ((8 - LEN F$) * 4), 169 SIZE 1, 2 BRUSH 2
1960   OPEN #11, F$
1970   INPUT #11, T%
1980   FOR D% := 1 TO T%
1990     INPUT #11, Yo(D%), Ro(D%)
2000   NEXT
2010   CLOSE #11
2020 ENDPROC
2030 PROCEDURE View
2040   GLOBAL Ro(), Yo(), T%, F$
2050   SET MODE 40
2060   SET DEG TRUE
2070   SET PAPER 2 : CLG
2080   SET WRITING 3 TO 2, 2; 38, 3
2090   SET WRITING 3
2100   AREA 10, 239; 319, 239; 319, 200; 100, 200 BRUSH 8
2110   AREA 10, 239; 319, 239; 319, 200; 100, 200 BRUSH 0 STYLE 0
2120   SET PAPER 15 : SET PEN 0 : SET BRUSH 136
2130   AREA 0, 249; 309, 249; 309, 210; 0, 210 BRUSH 15
2140   AREA 0, 249; 309, 249; 309, 210; 0, 210 BRUSH 0 STYLE 0
2150   CLS
2160   SET DRAWING 1 TO 0, 0; 319, 199
2170   SET DRAWING 1
2180   Angle RECEIVE Ang
2190   A := Ang
2200   SET PAPER 2 : CLG : SET PAPER 15
2210   S := SIN A : C := COS 30
2220   B := 1 : Cx := 160 : Cy := 75 : Dist := 100 : D := 200
2230   Mousey RECEIVE Cx, Cy
2240   SET PATTERN 136, 1 TO 0, 8, 0, 8
2250   SET PATTERN 136, 2 TO 8, 0, 8, 0
2260   SET PATTERN 136, 3 TO 0, 8, 0, 8
2270   SET PATTERN 136, 4 TO 8, 0, 8, 0
2280   SET PATTERN 137, 1 TO 8, 7, 8, 7
2290   SET PATTERN 137, 2 TO 7, 8, 7, 8
2300   SET PATTERN 137, 3 TO 8, 7, 8, 7
2310   SET PATTERN 137, 4 TO 7, 8, 7, 8
2320   SET PATTERN 138, 1 TO 15, 7, 15, 7
2330   SET PATTERN 138, 2 TO 7, 15, 7, 15
2340   SET PATTERN 138, 3 TO 15, 7, 15, 7
2350   SET PATTERN 138, 4 TO 7, 15, 7, 15
2360   FOR N := 1 TO T% - 2
2370     B := 0
2380     R1 := Ro(N) : R2 := Ro(N + 1) : Y1 := Yo(N) : Y2 := Yo(N + 1) GOSUB 2440
2390   NEXT
2400   ASK MOUSE X%, Y%, B%
2410   IF B% = 0 THEN 2400
2420   SET DRAWING 0
2430 ENDPROC
2440 FOR A := 135 - 185 - 180 TO 105 STEP 30
2450   X := R1 * SIN A
2460   Z := R1 * COS A
2470   Px1 := ((X - Z) * C) + Cx
2480   Py1 := (Y1 - (X + Z) * S) + Cy
2490   X := R2 * SIN A
2500   Z := R2 * COS A
2510   Px2 := ((X - Z) * C) + Cx
2520   Py2 := (Y2 - (X + Z) * S) + Cy
2530   X := R2 * SIN(A + 30)
2540   Z := R2 * COS(A + 30)
2550   Px3 := ((X - Z) * C) + Cx
2560   Py3 := (Y2 - (X + Z) * S) + Cy
2570   X := R1 * SIN(A + 30)
2580   Z := R1 * COS(A + 30)
2590   Px4 := ((X - Z) * C) + Cx
2600   Py4 := (Y1 - (X + Z) * S) + Cy
2610   AREA Px1, Py1; Px2, Py2; Px3, Py3; Px4, Py4
2620   AREA Px1, Py1; Px2, Py2; Px3, Py3; Px4, Py4 STYLE 3 BRUSH 15
2630   B := B + 1
2640   IF B = 1 THEN SET BRUSH 0
2650   IF B = 2 OR B = 12 THEN SET BRUSH 136
2660   IF B = 3 OR B = 11 THEN SET BRUSH 8
2670   IF B = 4 OR B = 10 THEN SET BRUSH 137
2680   IF B = 5 OR B = 9 THEN SET BRUSH 7
2690   IF B = 7 THEN SET BRUSH 15
2700   IF B = 8 OR B = 6 THEN SET BRUSH 138
2710 NEXT
2720 RETURN
2730 PROCEDURE Mousey RETURN Cx, Cy
2740   PRINT "Select centre of shape with mouser."
2750   SET MOUSE
2760   SET POINTS STYLE 2
2770   ASK MOUSE Cx%, Cy%, B%
2780   POINTS Cx, Cy OVER FALSE
2790   POINTS Cx, Cy OVER FALSE
2800 IF B% = 2 THEN ENDPROC
2810 GOTO 2770
2820 PROCEDURE Plott
2830   CREATE #15, "ltp1"
2840   GLOBAL T%, F$, Yo(), Ro()
2850   SET DEG TRUE
2860   SET PAPER 1 : SET BORDER 1 : CLG
2870   PLOT "Plot Drawing", 55, 220 SIZE 2, 2 BRUSH 0
2880   PLOT "Plot Drawing", 56, 219 SIZE 2, 2 BRUSH 2
2890   Angle RECEIVE Ang
2900   IF Ang = 0 THEN Sta := - 50
2910   IF Ang = 90 THEN Sta := - 230
2920   FOR N := 1 TO T% - 2
2930     Ox := 700 : Ov := 700 : Scx := 6 : Scv = 5 : R1 := Ro(N) : R2 := Ro(N + 1) : Y1 := Yo(N) : Y2 := Yo(N + 1) : GOSUB 3010
2940     X := Ro(N + 1) * SIN A
2950     Z := Ro(N + 1) * COS A
2960     PRINT #15, "M "; Ox - LEN F$ * 4; " "; Scv * Yo(0) - Ov
2970     PRINT #15, "P Drawing := ": F$
3000 CLOSE #15 : ENDPROC
3010 Cv := 400 : Cx := 400 : S := SIN Ang : C := COS 30
3020 FOR A := Sta TO 105 STEP 30
3030   X := R1 * SIN A
3040   Z := R1 * COS A
3050   Px1 := ((X - Z) * C) + Cx
3060   Py1 := (Y1 - (X + Z) * S) + Cy
3070   X := R2 * SIN A
3080   Z := R2 * COS A
3090   Px2 := ((X - Z) * C) + Cx
3100   Py2 := (Y2 - (X + Z) * S) + Cy
3110   X := R2 * SIN(A + 30)
3120   Z := R2 * COS(A + 30)
3130   Px3 := ((X - Z) * C) + Cx
3140   Py3 := (Y2 - (X + Z) * S) + Cy
3150   X := R1 * SIN(A + 30)
3160   Z := R1 * COS(A + 30)
3170   Px4 := ((X - Z) * C) + Cx
3180   Py4 := (Y1 - (X + Z) * S) + Cy
3190   PRINT #15, "M "; (Px1 * Scx) - Ox; " "; (Py1 * Scv) - Ov;
3200   PRINT #15, "D "; (Px1 * Scx) - Ox; " "; (Py1 * Scv) - Ov; " "; (Px2 * Scx) - Ox; " "; (Py2 * Scv) - Ov;
3210   PRINT #15, "D "; (Px2 * Scx) - Ox; " "; (Py2 * Scv) - Ov; " "; (Px3 * Scx) - Ox; " "; (Py3 * Scv) - Ov;
3220   PRINT #15, "D "; (Px3 * Scx) - Ox; " "; (Py3 * Scv) - Ov; " "; (Px4 * Scx) - Ox; " "; (Py4 * Scv) - Ov;
3230   PRINT #15, "D "; (Px4 * Scx) - Ox; " "; (Py4 * Scv) - Ov; " "; (Px1 * Scx) - Ox; " "; (Py1 * Scv) - Ov;
3240 NEXT
3250 RETURN
3260 PROCEDURE Angle RETURN Ang
3270   SET MOUSE
3280   PLOT "Please select an angle.", 65, 170 SIZE 1, 2 BRUSH 0
3290   PLOT "Please select an angle.", 66, 169 SIZE 1, 2 BRUSH 4
3300   PLOT "0", 50, 100 SIZE 5, 5 BRUSH 0
3310   PLOT "0", 51, 99 SIZE 5, 5 BRUSH 15
3320   PLOT "30", 190, 100 SIZE 5, 5 BRUSH 0
3330   PLOT "30", 191, 99 SIZE 5, 5 BRUSH 15
3340   PLOT "o", 85, 140 BRUSH 0 SIZE 1, 2
3350   PLOT "o", 86, 139 BRUSH 15 SIZE 1, 2
3360   PLOT "o", 265, 140 BRUSH 0 SIZE 1, 2
3370   PLOT "o", 266, 139 BRUSH 15 SIZE 1, 2
3380   ASK MOUSE X%, Y%, B%
3390   POINTS X%, Y% OVER FALSE BRUSH 15 STYLE 2
3400   POINTS X%, Y% OVER FALSE BRUSH 15 STYLE 2
3410   IF B% = 0 THEN 3380
3420   IF X% > 50 AND X% < 98 THEN Ang := 0 : GOTO 3450
3430   IF X% > 190 AND X% < 270 THEN Ang := 30 : GOTO 3450
3440   GOTO 3380
3450 ENDPROC
3460 SET MODE 40 : SET MOUSE
3470 PROCEDURE Cut
3480   GLOBAL T%, F$, Yo(), Ro()
3490   SET PAPER 1 : SET BORDER 1 : CLG
3500   Ymax% := 0 : Rmax% := 0 : Ymin% := 0
3510   FOR N := 1 TO T%
3520     IF Yo(N) > Ymax% THEN Ymax% := Yo(N)
3530     IF Yo(N) < Ymin% THEN Ymin% := Yo(N)
3540     IF ABS Ro(N) > Rmax% THEN Rmax% := ABS Ro(N)
3550   NEXT
3560   Cx := 160 : Cy := 125
3570   AREA Ymax% + Cx, Cy + Rmax%; Ymax% + Cx, Cy - Rmax%; Ymin% + Cx, Cy + Rmax%; Ymin% + Cx, Cy - Rmax%; Ymin% + Cx, Cy + Rmax% BRUSH 7
3580   AREA Ymax% + Cx, Cy + Rmax%; Ymax% + Cx, Cy - Rmax%; Ymin% + Cx, Cy + Rmax%; Ymin% + Cx, Cy - Rmax%; Ymax% + Cx, Cy + Rmax% BRUSH 0 STYLE 0
3590   FOR N := 1 TO T% - 2
3600     Cx := 160 : Cy := 125
3610     LINE Yo(N) + Cx, Cy + Ro(N); Yo(N + 1) + Cx, Cy + Ro(N + 1) BRUSH 15 STYLE 2
3620     LINE Yo(N) + Cx, Cy - Ro(N); Yo(N + 1) + Cx, Cy - Ro(N + 1) BRUSH 15 STYLE 2
3630   NEXT
3640   Maxin% := Ymax% + 160 : Tn% := 1
3650   FOR Ty := Rmax% TO 5 STEP - 5
3660     FOR Tx := Ymax% + 160 TO Ymin% + 160 STEP - 1
3670       Calc Tx RECEIVE Ra
3680       Tooldraw Tx, Ty, Tn%
3690       Erasetool Tx, Ty, Tn%
3700     IF Tv <> INT(Ra) THEN NEXT Tx
3705     IF Tx < Maxin% THEN Maxin% := Tx
3710   Tx := Ymax% + 160 : NEXT Ty
3711   FOR Tx := Ymax% + 160 TO Maxin% STEP - 1
3712     Calc Tx RECEIVE Ra
3713     Tooldraw Tx, Ty, Tn%
3714     Erasetool Tx, Ty, Tn%
3715   Ty := Ra : NEXT Tx
3720   ASK MOUSE X%, Y%, B%
3730   IF B% = 0 THEN 3720
3740 ENDPROC
3750 PROCEDURE Tooldraw Tx, Ty, Tn%
3760   GLOBAL Tlx%(), Tly%(), Pt%
3770   IF Tn% = Pt% THEN 3860
3780   IF Tn% = 1 THEN RESTORE 3880
3790   IF Tn% = 2 THEN RESTORE 3890
3800   READ Ts%
3810   FOR D% := 1 TO Ts%
3820     READ Tlx%(D%), Tly%(D%)
3830   NEXT
3840   FOR D% := 1 TO Ts% - 1
3850     LINE Tx + Tlx%(D%), (125 - Ty) + Tly%(D%); Tx + Tlx%(D% + 1), (125 - Ty) + Tly%(D% + 1) BRUSH 15
3860   NEXT
3870 ENDPROC
3880 DATA 5,0,0,10,-10,10,-50,0,-50,0,0
3890 DATA 5,0,-10,10,0,10,-50,0,-50,0,-10
3900 PROCEDURE Erasetool Tx, Ty, Tn%
3910   GLOBAL Tlx%(), Tly%(), Pt%
3920   IF Tn% = 1 THEN 3940
3930   IF Tn% = 2 THEN 3960
3940   AREA Tx + 0, (125 - Ty) - 0; Tx + 10, (125 - Ty) - 10; Tx + 10, (125 - Ty) - 50; Tx + 0, (125 - Ty) - 50; Tx + 0, (125 - Ty) - 0; Tx + 0, (125 - Ty) - 0 BRUSH 1
3950 AREA Tx + 0, (125 + Ty) + 0; Tx + 10, (125 + Ty) + 10; Tx + 10, (125 + Ty) + 50; Tx + 0, (125 + Ty) + 50; Tx + 0, (125 + Ty) + 0; Tx + 0, (125 + Ty) + 0 BRUSH 1 : ENDPROC
3960 AREA Tx + 0, Ty - 10; Tx + 10, Ty + 0; Tx + 10, Ty - 50; Tx + 0, Ty - 50; Tx + 0, Ty - 10 BRUSH 1
3970 AREA Tx + 0, Ty + 125 + 10; Tx + 10, Ty + 125 - 0; Tx + 10, Ty + 125 + 50; Tx + 0, Ty + 125 + 50; Tx + 0, Ty + 125 + 10 BRUSH 1 : ENDPROC
3980 PROCEDURE Calc Tx RECEIVE Ra
3990   GLOBAL Ro(), Yo(), T%, F$, Tlx%(), Tly%(), Pt%
4000   Calx := Tx - 160
4010   FOR N := 1 TO T% - 1
4020     IF Calx > Yo(N) AND Calx < Yo(N + 1) THEN 4080
4030     If Calx = Yo(N) THEN 4060
4040   NEXT
4050 Ra := 90 : ENDPROC
4060 Ra := Ro(N)
4070 ENDPROC
4080 L := Calx - Yo(N) : M := Yo(N + 1) - Calx
4090 Le := Yo(N + 1) - Yo(N) : H := Ro(N + 1) - Ro(N)
4100 H1 := (M * H) / (L + M)
4110 Ra := Ro(N) - H1 + H
4120 ENDPROC 



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
