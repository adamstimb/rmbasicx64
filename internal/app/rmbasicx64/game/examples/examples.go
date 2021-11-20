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
70 READBLOCK 0, 0, 0, 320, 250
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
			filename: "conways.BAS",
			program: `10 SET MODE 40
20 Box_Size% := 83 
30 REM Matrix_A is current, Matrix_B is next iteration
40 DIM Matrix_A%(Box_Size%, Box_Size%)
50 DIM Matrix_B%(Box_Size%, Box_Size%)
60 REM Set blinker
70 Matrix_A%(40, 40) = 1 : Matrix_A%(40, 41) = 1 : Matrix_A%(40, 42) = 1 
80 REM Set glider
90 Matrix_A%(5, 7) = 1 : Matrix_A%(6, 7) = 1 : Matrix_A%(7, 7) = 1 : Matrix_A%(7, 6) = 1 : Matrix_A%(6, 5) = 1
100 REM Main loop
110 REPEAT
120   CLS   
130   REM Draw the matrix
140   FOR X% := 0 TO Box_Size% - 2
150     FOR Y% := 0 TO Box_Size% - 2		 
160       IF Matrix_A%(X%, Y%) = 1 THEN POINTS X% * 3, Y% * 3 BRUSH 13
170     NEXT Y%
180   NEXT X%
190   REM Counts the count of the surrounding cell
200   REM Then apply the operation to the cell
210   FOR X% := 1 TO Box_Size% - 2
220     FOR Y% := 1 TO Box_Size% - 2    
230       REM Count the surrounding cells
240       Count% = 0
250       IF Matrix_A%(X% - 1, Y% + 1) = 1 THEN Count% = Count% + 1
260       IF Matrix_A%(X%, Y% + 1) = 1 THEN Count% = Count% + 1
270       IF Matrix_A%(X% + 1, Y% + 1) = 1 THEN Count% = Count% + 1
280       IF Matrix_A%(X% - 1, Y%) = 1 THEN Count% = Count% + 1
290       IF Matrix_A%(X% + 1, Y%) = 1 THEN Count% = Count% + 1
300       IF Matrix_A%(X% - 1, Y% - 1) = 1 THEN Count% = Count% + 1
310       IF Matrix_A%(X%, Y% - 1) = 1 THEN Count% = Count% + 1
320       IF Matrix_A%(X% + 1, Y% - 1) = 1 THEN Count% = Count% + 1
330       REM Apply the operations
340       REM Death
350       IF Matrix_A%(X%, Y%) = 1 THEN GOSUB Dead_Or_Alive
360       REM Birth
370       IF Matrix_A%(X%, Y%) = 0 THEN GOSUB Born_Or_Not
380     NEXT Y%
390   NEXT X%
400   REM Update the matrix with the new matrix that we have calculated.
405   REM We can optimize this with array referencing when implemented.
410   FOR X% := 0 TO Box_Size% - 1
420     FOR Y% := 0 TO Box_Size% - 1
430       Matrix_A%(X%, Y%) = Matrix_B%(X%, Y%)
440     NEXT Y%
450   NEXT X%
460 UNTIL TRUE = TRUE
470 SUBROUTINE Dead_Or_Alive
480 IF Count% = 2 OR Count% = 3 THEN Matrix_B%(X%, Y%) = 1 ELSE Matrix_B%(X%, Y%) = 0
490 RETURN
500 SUBROUTINE Born_Or_Not
510 IF Count% = 3 THEN Matrix_B%(X%, Y%) = 1 ELSE Matrix_B%(X%, Y%) = 0
520 RETURN
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
