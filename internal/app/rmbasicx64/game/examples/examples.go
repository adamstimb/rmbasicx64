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
}
