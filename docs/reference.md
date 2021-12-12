[Home](index.md) - [Quickstart](quickstart.md) - [History](history.md) - [**Reference**](reference.md) - [Releases](releases.md)

# Reference

This reference defines all the commands currently implemented in RM BASICx64 with details of any deviations from the original RM Basic implementation.  To learn RM Basic itself I recommend reading the original RM Basic manual, available for free download from the [Centre for Computing History](http://www.computinghistory.org.uk/det/47278/RM-Nimbus-PC-RM-Basic-PN-14351/).

# Filepaths

The original RM Nimbus shipped with MS-DOS 3.1 as standard and the differences between that are modern operating systems creates a many-bodied problem when trying to emulate the file operation behaviour of RM Basic, which implemented MS-DOS-like commands such as [DIR](#dir), [CHDIR](#chdir), etc. but also did some pre-handling before running the command on MS-DOS.  This could easily turn into a giagantic hairball in a cross-platform app so the following constraints were put in place:

- The filepath divider character is "\\", consistent with MS-DOS/Windows.
- The Workspace Directory set during the installation process is regarded as root ("\\").
- It is not possible to access folders above the root.
- Only very basic behaviour is implemented, e.g. switching between subdirectories, deleting or renaming files individually, creating subdirectories etc.
- Relative paths (e.g. "..\\mydir") are not supported.
- Files cannot be deleted or renamed recursively or with wildcards; these operations are supported only for individual files.
- As with RM Basic, RM BASICx64 only supports ASCII character encoding so file names or paths containing unicode characters cannot be accessed.
- Unlike RM Basic (and MS-DOS 3.1) filepaths are case-sensitive.

# Keywords

The format, punctuation and options are shown using the following symbols:

_..._
indicates the more items can be used

_exp, exp1_
numeric or string expression

_e, e1, ..._
numeric expressions

_e$, e1$, ..._
string expressions

_v, v1, ..._
variable names

_t_
a numeric expression that is evaluated as a boolean; "truthiness" in RM Basic was defined as -1.0 is true, any other value is false.

_lineNumber, lineNumber ..._
program line numbers

## ABS

Calculate the absolute value of a number.

### Syntax

ABS(_e_)

### Example

```
PRINT ABS(-1.2345)

   1.2345

```

## AND

Bitwise AND on two expressions.

### Syntax

_e1_ AND _e2_

## AREA

Draw a filled polygon on the screen.

### Syntax

AREA _coordinateList_ [_optionList_]

### Remarks

The AREA command draws a shape by connecting each point from a coordinate list to the next, starting and finishing at the first point in the list.  Each coordinate is seperated by a semicolon, e.g. 0, 0; 50, 0; 50, 50; 0, 50

After the coordinate list you can specify options that override the current graphics settings:

| Option | Action | Syntax |
| ------ | ------ | ------ |
| BRUSH | selects the brush colour | BRUSH _e_ |
| STYLE | the fill style | STYLE _e1_[, _e2_[, _e3_]] |
| OVER | selects the drawing style | OVER _t_ |

### Example

```
10 REM Draw a red triangle
20 SET MODE 40
30 AREA 0, 0; 100, 0; 50, 100 BRUSH 2
```

## ASK MOUSE

Get the current position and button state of the mouse.  

### Syntax

ASK MOUSE [_v1_, _v2_][, _v3_]

### Remarks

The mouse must first be initialised with the [SET MOUSE](#set-mouse) command.  _v1_ and _v2_ are the variables in which the X and Y position of the mouse cursor will be stored.  Note that these are the positions on the Nimbus's screen; if the mouse is outside RM BASICx64 application window then only its last known on-screen position will be available.  _v3_ is the variable to store the button state (0 no button pressed, 1right-hand button pressed, 2 left-hand button pressed, 3 middle or both buttons pressed depending on your mouse and operating system).  Unlike original RM Basic, it is not possible to set the starting position or sensitivity of the mouse in RM BASICx64.

### Example

```
10 REM A very, very simple drawing program
20 SET MODE 40
30 PRINT "Click any mouse button to quit"
40 SET MOUSE
50 REPEAT
60   ASK MOUSE Xpos%, Ypos%, Button%
70   CIRCLE 5, Xpos%, Ypos% BRUSH 13
80 UNTIL Button% > 0
```

## ATN

Calculate the angle with the given tangent.  The unit of the measurement for the angle can be set with [SET DEG](#set-deg) or [SET RAD](#set-rad).

### Syntax

ATN(_e_)

### Example

```
SET DEG TRUE
PRINT ATN(1.557)

   90

```

## BYE

Quit the application.

## CHDIR

Change the current working directory.

### Syntax

CHDIR _e$_

### Remarks

See [Filepaths](#filepaths) for restrictions.

## CHR$

Return an ASCII character.

### Syntax

CHR$(_e_)

### Remarks

_e_ is an integer representing the decimal code of a charcter in the [Extended ASCII table](https://www.ascii-code.com/)  Note that in charset 1 the non-alphabetic characters are very different to the standard MS-DOS/IBM characters and are unique to the RM Nimbus.

## CIRCLE

Draw one or more circles on the screen.

### Syntax

CIRCLE _e_, _coordinateList_ [_optionList_]

### Remarks

_e_ is the radius of the circle.  The coordinate list can be a single set of coordinates or a list of several, separated by a semicolon, e.g. 10, 200 or 10, 200; 20, 200; 30, 200

After the coordinate list you can specify options that override the current graphics settings:

| Option | Action | Syntax |
| ------ | ------ | ------ |
| BRUSH | selects the brush colour | BRUSH _e_ |
| STYLE | the fill style | STYLE _e1_[, _e2_[, _e3_]] |
| OVER | selects the drawing style | OVER _t_ |

### Example

```
CIRCLE 20, 0, 200
CIRCLE 15, 50, 100 BRUSH 2
CIRCLE 30, 0, 150; 50, 150; 100, 150 BRUSH 1 OVER FALSE
```

## CLOSE

Close a file channel.

### Syntax

CLOSE [#_e_]

### Remarks

`CLOSE` without an argument closes all file channels.  To close a specific channel, pass the channel number prefixed by `#`.

## CLS

Clears the screen or a selected writing area.

### Syntax

CLS [~_e_]

### Remarks

`CLS` clears the entire screen.  To clear just one writing area, pass the number of the writing area after a tilde, e.g. `CLS ~1`.

## COS

Calculate the cosine of an angle.  The unit of the measurement for the angle can be set with [SET DEG](#set-deg) or [SET RAD](#set-rad).

### Syntax

COS(_e_)

### Example

```
SET DEG TRUE
PRINT COS(90)

   -25.67272711536642

```

## CREATE

Open a file channel in writing mode.

### Syntax

CREATE #_e1_, _e2$_

### Remarks

As in the original RM Basic, channels #11 to #127 are user-defined and can be assign to files.  The filename _e2$_ must be valid file path (see [Filepaths](#filepaths) for details).

## DATA

Specify numeric and/or string constants that will be assigned to variables with the READ statement.

### Syntax

DATA _c1_[, _c2_...]

### Remarks

DATA statements can be placed anywhere in your program and unlike function or procedure statements they can be executed although no side-effects will be noticed.  A program can have many DATA statements and the data in them will be read in the order in which they appear.  All DATA statements are read into memory _before_ the program itself executes, and the values are read into variables using READ statements.

See the RM Basic manual for the details!

## DIR

Print a directory listing

### Syntax

DIR [#_e1_,] [~_e2_,] [_e3$_]

### Remarks

_e1_ is a file channel that has already been open in writing mode (see [CREATE][#create] for details).

_e2_ is a writing area that has been defined with [SET WRITING][#set-writing]; this is ignored if a file channel has been passed.

`DIR` without an argument lists all .BAS files in the current working directory.  Old-school MS-DOS wildcards are supported, so to list all JPGs use `DIR "*.JPG"` and all files `DIR "*.*"`.  To list folders in a subdirectory use `DIR "myfolder\"`.  Just like RM Basic, if the file extension is omitted as in the last example, *.BAS is automatically appended.  Unlike RM Basic (and MS-DOS 3.1) the filepaths are case-sensitive.

See [Filepaths](#filepaths) for restrictions.

## EDIT

Edit a line number in a program

### Syntax

EDIT _lineNumber_

## END

End program execution

### Syntax

END

## ERASE

Erase a file in the current working directory.

### Syntax

ERASE _e$_

### Remarks

See General - Filepaths for restrictions.

## EXP

Calculate the exponential function, e^x

### Syntax

EXP(_e_)

### Example

```
PRINT EXP(1)

   2.718281828459045

```

## FLOOD

Fill an area of the graphics screen.

## Syntax

FLOOD _coordinateList_ [_optionList_]

### Remarks

Starting from the point of points given in the coordinate list, the screen is filled with the brush colour until a boundary or drawing area edge is reached.

After the coordinate list you can specify options that override the current graphics settings:

| Option | Action | Syntax |
| ------ | ------ | ------ |
| BRUSH | selects the brush colour | BRUSH _e_ |
| STYLE | the fill style | STYLE _e1_[, _e2_[, _e3_]] |
| OVER | selects the drawing style | OVER _t_ |
| EDGE | selects the boundary colour | OVER _t_ |

## FOR ... NEXT

Repeat a series of instruction, altering a control variable on each repetition.

### Syntax

FOR _v_ [:]= _e1_ TO _e2_ [STEP _e3_]
:
Instructions
:
NEXT [_v_]

### Example

```
10 PRINT "Countdown"
20 FOR I% := 5 TO 0
30   PRINT I%
40 NEXT I%
50 PRINT "Blast off!"
RUN

  Countdown
  5
  4
  3
  2
  1
  0
  Blast off!
```

## FUNCTION / RESULT / ENDFUN

Define a function.

### Syntax

FUNCTION _v1_([_v2_ [ ,_v3_...]])

RESULT _e1_ [, _e2_ ...]

ENDFUN

### Remarks

Functions can be defined in RM Basic much like in any modern language.  The definition can be placed anywhere in your program, so even if you call a function before it's defined, the function will still be callable.  The only gotcha is that the FUNCTION command itself cannot be executed.  A good way to avoid this is to put all your function statements at the end of the program, and insert an END statement above as shown in the example below.  The result is returned to the caller whenever RESULT is called from within the function.  Note that RM Basic functions can only return one value.  To return more than one value, bizarrely enough you don't need a function at all: You need a procedure!  The ENDFUN statement marks the end of the function.  Although not strictly enforced in RM Basic, execution can be unpredictable if the ENDFUN statement is left out.

### Example

```
10 REM Simple function to generate a greeting
20 Name$ := "Slim Shady"
30 PRINT Generate_Greeting$(Name$)
40 END : REM Don't execute function definitions below
50 FUNCTION Generate_Greeting(N$)
60   RESULT "Hi! My name is " + N$
70 ENDFUN
```

## GET

Read the code of a character from the keyboard if a key was pressed.

### Syntax

GET([_e_])

## GOSUB

### Syntax

GOSUB _label_

### Remarks

Jump to a SUBROUTINE.  Execution will be returned to the place where you jumped when the subroutine ends with a RETURN statement.

## GLOBAL

Create a global variable, or set up a procedure or function to access a global variable.

### Syntax

GLOBAL _v_

### Remarks

To make a variable global, declare it as global in the main program block before assigning any values to it.  Then in any procedure or function, declare it again to make the variable accessible.

### Example

```
10 GLOBAL Cat_Name$
20 Cat_Name$ := "Bobby"
30 PRINT "First his name was "; Cat_Name$
40 Change_Name
50 PRINT "And now it's "; Cat_Name$
60 END
70 PROCEDURE Change_Name
80   GLOBAL Cat_Name$
90   Cat_Name$ = "Chubby"
100 ENDPROC

Output:
First his name was Bobby
And now it's Chubby
```

## GOTO

Interrupt the flow of the program and jump to any given line number.

### Syntax

GOTO _lineNumber_

### Example

```
10 REM This goes out to LGR
20 PRINT "Farts!"
30 GOTO 20
```

## HOME

Return the cursor to the top-left corner of the screen

### Syntax

HOME

## IF...THEN...ELSE

Conditionally execution instruction(s) on a single line.

### Syntax

IF _t_ THEN Instruction(s) [ELSE Instruction(S)]

### Example

```
IF Month = 9 THEN PRINT "September"
IF Month = 12 AND Day = 31 THEN SET PEN 2 : PRINT "Happy New Year!" ELSE SET PEN 1 : PRINT "Have a nice day!"
```

## INPUT

Receive input and assign input to a variable.

### Syntax

INPUT [#_e1_,] [~_e1_,] _e$_[;] _v_

### Remarks

_e1_ is a file channel that has already been open in reading mode (see [OPEN][#open] for details).

_e2_ is a writing area that has been defined with [SET WRITING][#set-writing]; this is ignored if a file channel has been passed.

_e$_ is a prompt that is written to the screen.  If a semicolon follows the prompt, a "?" character is printed after the prompt.  The user can then key in a response which is parsed and stored in the variable _v_.  Note that parsing input into multiple variables is not yet supported.  The prompt is ignore if a file channel has been passed.

### Example

```
10 PRINT "Mirror, mirror, on the wall,
20 INPUT "Who is the most irritating cat of them all?", Name$
30 IF Name$ = "Fluffy" OR Name$ = "fluffy" THEN PRINT "I totally agree" ELSE PRINT "I don't know a cat called", Name$

Output:
Mirror, mirror, on the wall,
Who is the most irritating cat of them all? Fluffy
I totally agree
```

## INT

Calculate the largest whole number that is less than or equal to a given value.

### Syntax

INT(_e_)

### Example 

```
PRINT INT(34.99999999)

   34

```

## LEN

Return the number of characters in a string.

### Syntax

LEN(_e$_)

```Example
PRINT LEN("Hello Mable")

   11

```

## LET

Assign the value of an expression to a variable.

### Syntax

[LET] v [:]= _e_
_or_
[LET] v$ [:]= _e$_

### Remarks

Yes this means there are several equivalent expressions to assign variables in RM Basic as shown in the example below, and it just seems to be a matter of personal preference.  `LET` is rarely mentioned in the manual, sometimes `=` is used, sometimes `:=`; pick one and stick with it!

### Example

```
Cats% := 5
Dogs% = 1
Temperature := 18.57
Country$ := "Poland"
LET Total_Animals% := Cats% + Dogs%
LET Overkill% = TRUE
```

## LINE

Draw a series of connected lines on the screen.

### Syntax

LINE _coordinateList_ [_optionList_]

### Remarks

The LINE command draws a shape by connecting each point from a coordinate list to the next.  Each coordinate is seperated by a semicolon, e.g. 0, 0; 50, 0; 50, 50; 0, 50

After the coordinate list you can specify options that override the current graphics settings:

| Option | Action | Syntax |
| ------ | ------ | ------ |
| BRUSH | selects the brush colour | BRUSH _e_ |
| STYLE | not yet implemented | n/a |
| OVER | selects the drawing style | OVER _t_ |

### Example

```
10 REM Draw a green wire triangle
20 SET MODE 40
30 LINE 0, 0; 100, 0; 50, 100; 0, 0 BRUSH 4
```

## LIST

List the stored program.

### Syntax

LIST [#_e1_,] [~_e2_,] [_e3_] [TO [_e4_]]

### Remarks

_e1_ is a file channel that has already been open in writing mode (see [CREATE][#create] for details).

_e2_ is a writing area that has been defined with [SET WRITING][#set-writing]; this is ignored if a file channel has been passed.

`LIST` by itself lists the entire program.  A single line can be listed by passing the line number, e.g. `LIST 130`.  The program from a particular line to the end can be listed by passing an unlimited range, e.g. `LIST 130 TO`.  The program can be listed between to lines by passing a limited range, e.g. `LIST 90 TO 130`.

## LN

Calculate the natural logarithm of a number.

### Syntax

LN(_e_)

### Example

```
PRINT LN(1.32)

   0.27763173659827955

```

## LOAD

Load a program from a file into memory.

### Syntax

LOAD _e$_

### Remarks

Where _e$_ must be a valid filename.  If _e$_ does not end in ".BAS" then ".BAS" will be added automatically.

## LOG

Calculate the logarithm to the base 10 of a number.

### Syntax

LOG(_e_)

### Example

```
PRINT LOG(5.2)

   0.7160033436347991

```

## LOOKUP

Check if a file exists in the current working directory and return TRUE or FALSE.

### Syntax

LOOKUP(_e$_)

### Remarks

See [Filepaths](#filepaths) for restrictions.

## MKDIR

Create a subdirectory in the current working directory.

### Syntax

MKDIR _e$_

### Remarks

See [Filepaths](#filepaths) for restrictions.

## MOD

Returns the remainder of integer division.

### Syntax

_e1_ MOD _e2_

## MOVE

Move the cursor relative to its current position.

### Syntax

MOVE _e1_, _e2_

### Remarks

_e1_ is the number of columns move, and _e2_ is the number of rows to move, relative to the current cursor position.

## NEW

Clear workspace.  Delete all variables and wipe the stored program.

### Syntax

NEW

## NOT

Bitwise NOT on an expression.

### Syntax

NOT _e_

## NOTE

Play a note.

### Syntax

NOTE _e1_ [TO _e2_] [,_e3_ [ ,_e4_]] [ENVELOPE _e5_] [VOICE _e6_]

### Remarks

Not all options are implemented.  This thing is complicated.  Please refer to the original manual!

## OPEN

Open a file channel in reading mode.

### Syntax

OPEN #_e1_, _e2$_

### Remarks

As in the original RM Basic, channels #11 to #127 are user-defined and can be assign to files.  The filename _e2$_ must be valid file path (see [Filepaths](#filepaths) for details).

## OR

Bitwise OR on two expressions.

### Syntax

_e1_ OR _e2_

## PATH$

Returns the current working directory.

### Syntax

PATH$

### Remarks

See [Filepaths](#filepaths) for restrictions.

## PLOT

Draw graphics characters on the screen.

### Syntax

PLOT _e$_, _coordinateList_ [_optionList_]

### Remarks

The PLOT command draws _e$_ at all the coordinates in the coordinateList.  Each coordinate is seperated by a semicolon, e.g. 0, 0; 50, 0; 50, 50; 0, 50

After the coordinate list you can specify options that override the current graphics settings:

| Option | Action | Syntax |
| ------ | ------ | ------ |
| BRUSH | selects the brush colour | BRUSH _e_ |
| OVER | selects the drawing style | OVER _t_ |
| DIRECTION | selects the drawing direction | DIRECTION _e_ |
| SIZE | selects the size of characters | SIZE _e1_[, _e2_] |
| FONT | selects the font of characters | FONT _e_ |

## PITCH

Return the note number for a given octave and note.

### Syntax

PITCH(_e1_, _e2_)

### Remarks

_e1_ is the octave (0 being middle), and _e2_ is the note number, 0 being C, i.e. middle C is returned by PITCH(0, 0).

## POINTS

Draw one or more points on the screen.

## Syntax

POINTS _coordinateList_ [_optionList_]

### Remarks

The coordinate list can be a single set of coordinates or a list of several, seperated by a semicolon, e.g. 10, 200 or 10, 200; 20, 200; 30, 200

After the coordinate list you can specify options that override the current graphics settings:

| Option | Action | Syntax |
| ------ | ------ | ------ |
| BRUSH | selects the brush colour | BRUSH _e_ |
| STYLE | the points style | STYLE _e_ |
| OVER | selects the drawing style | OVER _t_ |

## PRINT

Prints strings and/or numbers on the screen.

### Syntax

PRINT [#_e1_,] [~_e2_,] [_print list_]

### Remarks

_e1_ is a file channel that has already been open in writing mode (see [CREATE][#create] for details).

_e2_ is a writing area that has been defined with [SET WRITING][#set-writing]; this is ignored if a file channel has been passed.

_print list_ is a list of expressions (numeric and/or string). Each expression must be separated by a semicolon, comma, space or exclamation mark.  The expressions are evaluated and then printed on screen.  Using semicolon or space between the expressions causes the results to be printed on the same line immediately following one another.  Using a comma causes the following result to be printed in the next print zone (this isn't properly implemented yet).  Using an exclamation mark causes the next result to be printed on a new line.  Ending the _print list_ with a semicolon or comma causes the cursor to remain on the same line, so that the next PRINT statement not start on a new line.

### Example

```
PRINT "My mind is going"
    
   My mind is going

PRINT 10 * 100

   1000

First_Name$ := "Dave"
Last_Name$ := "Bowman"
PRINT First_Name$ + " " + Last_Name$

   Dave Bowman

PRINT First_Name$; Last_Name$

   DaveBowman

PRINT First_Name$ !! Last_Name$

   Dave

   Bowman
```

## PROCEDURE / RETURN / RECEIVE / LEAVE / ENDPROC

Define a procedure.

### Syntax

PROCEDURE _v1_ [_v2_ [ ,_v3_...]] [RECEIVE [_v4_ [ , _v5_ ...]]]

LEAVE

ENDPROC

### Remarks

When is a function not a function?  When it's a procedure.  When is a procedure not a procedure? When it's a procedure that can receive arguments and return a value; in fact it can return several values, making it a kind of monster function!  Confusing?  Yep.  Ahead of it's time and brilliant?  Absolutely.

As with functions, the definition can be placed anywhere in your program, so even if you call a procedure before it's defined, the procedure will still be callable.  The PROCEDURE command itself cannot be executed.  To avoid this is to put all your procedure statements at the end of the program, and insert an END statement above as shown in the example below.  The result is returned to the caller whenever LEAVE or ENDPROC is called from within the procedure.  The ENDPROC statement marks the end of the function.  Like ENDFUNC, although not strictly enforced in RM Basic, execution can be unpredictable if the ENDPROC statement is left out.

### Examples

```
10 Boom
20 END
30 PROCEDURE Boom
40   SET MODE 40 : SET BORDER 14 : SET PAPER 2 : CLS
50   Shadow_Text "BOOOMM!!!", 10, 10, 4, 15
60 ENDPROC
70 PROCEDURE Shadow_Text Msg$, X%, Y%, S%, C%
80   PLOT Msg$, X%, Y% SIZE S% BRUSH 0
90   PLOT Msg$, X% + 1, Y% + 1 SIZE S% BRUSH C%
100 ENDPROC
```

```
10 Cat_Name$ := "Fluffy"
20 Dog_Name$ := "Czeszek"
30 Cat_Food$ := "Salmon"
40 Dog_Food$ := "Steak"
50 Describe_Pets Cat_Name$, Dog_Name$, Cat_Food$, Dog_Food$ RECEIVE Names$, Foods$
60 PRINT Names$
70 PRINT Foods$
80 END
90 PROCEDURE Describe_Pets Pet_1_Name$, Pet_2_Name$, Pet_1_Food$, Pet_2_Food$ RETURN Their_Names$, Their_Foods$
100   Their_Names$ = "The names of our pets are " + Pet_1_Name$ + " and " + Pet_2_Name$
110   Their_Food$ = Pet_1_Name$ + "'s favourite food is " + Pet_1_Food$ + " but " + Pet_2_Name$ + " likes " + Pet_2_Food$
120 ENDPROC
```

## PUT

Write one or more ASCII characters to the screen.

### Syntax

PUT [~_e1_] _e2_[, _e4_ ...]

### Remarks

PUT does not add a carriage return like [PRINT](#print).

## READ

Read values from DATA and assign them to variables.

### Syntax

READ _v1_[, _v2_...]

### Remarks

See the RM Basic manual for details.

## REM

Insert a comment.

## Syntax

REM _comment_

## RENAME

Rename a file in the current working directory.

### Syntax

RENAME _e1$_ TO _e2$_

### Remarks

See [Filepaths](#filepaths) for restrictions.

## RENUMBER

Renumber the program lines.  Currently no arguments are supported, only the default option to renumber the entire program with the first line given the number 10, and all subsequent lines incremented by 10.

### Example

```
30 PRINT "Hello"
11 PRINT "Blah"
23 PRINT "Meh"
5 CLS

LIST
   5 CLS
   11 PRINT "Blah"
   23 PRINT "Meh"
   30 PRINT "Hello"

RENUMBER
LIST
   10 CLS
   20 PRINT "Blah"
   30 PRINT "Meh"
   40 PRINT "Hello"
```

## REPEAT ... UNTIL

Repeat a series of instructions until a condition is met.

### Syntax

REPEAT
:
Instructions
:
UNTIL _t_

### Example

```
10 REM Roll dice until we get 2 sixes
20 Throws% := 0
30 REPEAT
40   D1% := RND(6) : D2% := RND(6)
50   PRINT "Throw ", Throws%, ": ", D1%, "+", D2%
60   Throws% = Throws% + 1
70 UNTIL D1% = 6 AND D2% = 6
80 PRINT "We got 2 sixes after ", Throws%, " throws!"
```

## RESTORE

Prepare to reread DATA instructions.

### Syntax

RESTORE [_lineNumber_]

### Remarks

See the RM Basic manual for details.

## RMDIR

Remove a subdirectory in the current working directory.

### Syntax

RMDIR _e$_

### Remarks

See [Filepaths](#filepaths) for restrictions.

## RND

Generate a random number, or re-seed the random number generator.  

### Syntax

RND(_e_)

### Remarks

Pass any negative number to re-seed.  To return random floating-point number pass 1.  To return a random integer up to a maximum value, pass the maximum value.

## RUN

Execute the stored program.  Running from a different line number is not yet implemented.

### Syntax

RUN

## SAVE

Save a stored program to a file.

### Syntax

SAVE _e$_

### Remarks

_e$_ must be a valid filename.  Wildcard characters are not allowed.  If the file already exists the user is prompted with a warning and asked if the operation should be aborted.  If _e$_ does not end in ".BAS" then ".BAS" will be added automatically.

See [Filepaths](#filepaths) for restrictions.

## SET BORDER

Change the border colour.

### Syntax

SET BORDER _e_

### Example

```
SET BORDER 2
```

## SET COLOUR

Assign colours to the current pallete and/or set flashing colours and flash speed.  

### Syntax

SET COLOUR _e1_ TO _e2_[,_e3_,_e4_]

### Remarks

_e1_ is the number of the current pallete you want to set.  _e2_ indicates the value of the base colour to be assigned to _e1_.  The list of base colours is given below.  _e4_ indicates a second colour that will flash regularly with _e2_.  _e3_ specifies the flashing speed (0 no flash, 1 slow flash, 2 fast flash).

| Value of _e2_ or _e4_ | Colour |
| --------------------- | ------ |
| 0                     | Black |
| 1                     | Dark blue |
| 2                     | Dark red |
| 3                     | Purple |
| 4                     | Dark green |
| 5                     | Dark cyan |
| 6                     | Brown |
| 7                     | Light grey |
| 8                     | Dark grey |
| 9                     | Light blue |
| 10                    | Light red |
| 11                    | Magenta |
| 12                    | Light green |
| 13                    | Cyan |
| 14                    | Yellow |
| 15                    | White |

### Example

```
10 REM Write a flashing yellow warning message
20 REM on a dark grey background in hi-res mode
30 SET MODE 80
40 SET COLOUR 0 TO 8
50 SET COLOUR 1 TO 14, 2, 8
60 SET COLOUR 3 TO 0
70 SET PEN 1 : PRINT "WARNING - Imminent meltdown!"
80 SET PEN 2 : PRINT "Evacuate to at least 100 km distance immediately."
90 SET PEN 3 : PRINT "Good luck and have a nice day."
```

## SET CONFIG BOOT

This a new command only implemented in RM BASICx64.  It is used to enable or disable the RM Nimbus "Welcome" boot sequence when RM BASICx64 starts.

### Syntax

SET CONFIG BOOT _t_

### Example

```
SET CONFIG BOOT TRUE
```

## SET CURPOS

Move the cursor to a specific position.

### Syntax

SET CURPOS _e1, _e2_

### Remarks

_e1_ is the column number and _e2_ is the row number to move the cursor to.

## SET DEG

Set the angle measurement unit to degrees.  Note that this is equivalent to [SET RAD](#set-rad) which sets the angle measurement to radians.

### Syntax

SET DEG _t_

### Example

```
SET DEG TRUE
```

## SET DRAWING

Select a drawing area or define the boundaries of a drawing area.

### Syntax

SET DRAWING _e1_ [TO _e2_, _e3_; _e4_, _e5_]

### Remarks

Select a drawing area by passing one argument, _e1_.  The Nimbus had 10 drawing areas, 0 being the entire screen and cannot be user-defined.  Drawing areas 1 to 10 are user definable.  To define a drawing area, pass in the x, y coordinates of the bottom-left and top-right corners that will form the boundaries of the area.

When a user-defined writing area is selected, all graphics commands  will intepret x, y coordinates relative to the drawing area's boundaries.

### Example

```
SET DRAWING 0 : REM Select entire screen for drawing
SET DRAWING 5 TO 10, 10; 100, 100 : REM Define a drawing box near the bottom-left corner of the screen
```

## SET ENVELOPE

Select a sound envelope, or define a sound envelop.

### Syntax

SET ENVELOPE _e1_ [TO _e2_, _e3_; _e4_, _e5_; _e6_, _e7_; _e8_]

### Remarks

Select an envelope by passing one argument, _e1_.  The Nimbus had 9 definable sound envelopes (0 - 9).  To define a drawing area, pass all the arguments _e2_ to _e8_ corresponding to the envelop properties as listed here:

| Argument | Property |
| -------- | -------- |
| _e2_ | Attack time |
| _e3_ | Attack level |
| _e4_ | Decay time |
| _e5_ | Decay level |
| _e6_ | Sustain time |
| _e7_ | Sustain level |
| _e8_ | Release time |

Times are in centiseconds, levels are 0 - 15 (0 being inaudible)

## SET MODE

Change the screen mode between high-resolution, 4-colour mode (80) and low-resolution, 16-colour mode (40)

### Syntax

SET MODE _e_

### Example 

```
SET MODE 40
SET MODE 80
```

## SET PAPER

Change the paper colour.

### Syntax

SET PAPER _e_

### Example

```
SET PAPER 1
```

## SET PATTERN

Define a pattern that can be used as a BRUSH colour when drawing.

### Syntax

SET PATTERN _e1_, _e2_ TO _e3_, _e4_, _e5_, _e6_

### Remarks

See the RM Basic manual for how this works and the default pattern settings!

## SET PEN

Change the pen colour.

### Syntax

SET PEN _e_

### Example

```
SET PEN 2
```

## SET RAD

Set the angle measurement unit to radians.  Note that this is equivalent to [SET DEG](#set-rad) which sets the angle measurement to degrees.

### Syntax

SET RAD _t_

### Example

```
SET RAD TRUE
```

## SET SOUND

Turn the Nimbus sound engine on or off.

### Syntax

SET SOUND _t_

### Remarks

`SET SOUND TRUE` must be executed before using any of the sound-related commands.

## SET TONE

Switch the current voice between square-wave and white noise.

### Syntax

SET TONE _t_

### Remarks

`SET TONE TRUE` sets the current voice to square-wave, `SET TONE FALSE` sets it to white noise.  The Nimbus used pink noise but this has not yet been implemented.

## SET VOICE

Select a voice to play sounds.

### Syntax

SET VOICE _e1_

### Remarks

There are 3 voices, 1 - 3.

## SET WRITING

Select a writing area (textbox) or define the boundaries of a writing area.

### Syntax

SET WRITING _e1_ [TO _e2_, _e3_; _e4_, _e5_]

### Remarks

Select a writing area by passing one argument, _e1_.  The Nimbus had 10 writing areas, 0 being the entire screen and cannot be user-defined.  Writing areas 1 to 10 are user definable.  To define a writing area, pass in the column, row cursor positions of the bottom-left and top-right corners that will form the boundaries of the area.

When a user-defined writing area is selected, [SET CURPOS](#set-curpos) will reposition the cursor relative to the writing area's boundaries.

### Examples

```
SET WRITING 0 : REM Select the entire screen for writing
SET WRITING 3 TO 1, 1; 10, 10 : REM Define a small writing area in the top-left of the screen
```

## SIN

Calculate the sine of an angle. The unit of the measurement for the angle can be set with [SET DEG](#set-deg) or [SET RAD](#set-rad).

### Syntax

SIN(_e_)

### Example

```
SET DEG TRUE
PRINT SIN(90)

   1

```

## SQR

Calculate the square root of a number.

### Syntax

SQR(_e_)

### Example

```
PRINT SQR(23)

   4.795831523312719

```

## STR$

Convert a number into a string representation.

### Syntax

STR$(_e_)

## SUBROUTINE ... RETURN

Label a section of code as a subroutine.

### Syntax

SUBROUTINE _label_

...

RETURN

### Remarks

See RM Basic manual for details.

## TAN

Calculate the tangent of an angle. The unit of the measurement for the angle can be set with [SET DEG](#set-deg) or [SET RAD](#set-rad).

### Syntax

TAN(_e_)

### Example

```
```

## XOR

Bitwise XOR on two expressions.

### Syntax

_e1_ XOR _e2_

# ANIMATE Extension

This extension was introduced in RM Basic 2.00C released in 1987.  It was used to store and retrieve blocks of video memory and to load and save images in PaintSPA's file format.  ANIMATE has been re-implemented in RM BASICx64, supporting instead full-colour JPG and BMP files which, upon loading, are downsampled to which ever colour pallete is being used at the time.  Yes - this means that your 16 million colour JPG will be rendered with only 4 colours if you load it in MODE 80!  You can also save images in JPG or BMP format, making it possible to share screenshots and even generate memes with RM Basic (see `meme.BAS` in the example programs).  Keep in mind that the resolution of the Nimbus is tiny by today's standards (320x250 in MODE 40) so it is recommended to scale down images to a comparable size beforehand.  Results are often further improved by boosting the contrast and brightness as well.

The syntax and original documentation of the ANIMATE extensions's command aren't _quite_ consistent with the core RM Basic commands.  For authenticity these inconsistencies have been left in instead of being "fixed".

## READBLOCK

Read the data displayed in a specified area of the screen into a numbered block of memory numbered 0 - 99.

### Syntax

READBLOCK _block-number_, _x-min_, _y-min_; _x-max_, _y-max_

### Example

```
10 SET MODE 80
20 REM Read the whole display into memory block 0
30 READBLOCK 0, 0, 0; 639, 249
```

## WRITEBLOCK 

Display the contents of a numbered block of memory at a specified position on the screen.

### Syntax

WRITEBLOCK _block-number_, _x-pos_, _y-pos_ [, _plot-mode_]

### Example

```
10 REM Display the contents of block 0
20 WRITEBLOCK 0, 0, 0, -1, 1
```

### Remarks

The specified memory block must have been previously allocated by READBLOCK or FETCH.

Set _plot-mode_ to 0 for XOR plotting, or -1 for OVERWRITE plotting (-1 is default).

Selecting a transparency colour has not yet been implemented.

## SQUASH

Same syntax and similar behaviour to WRITEBLOCK expect the image is scaled to 1/16 size before writing.

## ASK BLOCKSIZE

Returns the width and height (in pixels) of the specified memory block and the screen width (in characters) that was in use when the block was created.

### Syntax

ASK BLOCKSIZE _block-size_, _x-pixels_ [ , _y-pixels_ [ ,_screen-chars_ ]]

### Example

```
10 SET MODE 80
20 READBLOCK 0, 0, 0; 639, 249
30 ASK BLOCKSIZE 0, X%, Y%, M%
40 PRINT X%, Y%, M%
60 REM Prints 640    250    80
```

## COPYBLOCK 

Copy the data displayed in one area of the screen into another area.

### Syntax

COPYBLOCK _x-min_, _y-min_; _x-max_, _y-max_; _x-dest_, _y-dest_ [ , _plot-mode_ ]

## DELBLOCK 

Delete a numbered block of memory.

### Syntax

DELBLOCK _block-number_

## CLEARBLOCK 

Delete all blocks of memory.

### Syntax

CLEARBLOCK

## FETCH

Loads the contents of the specified image file into the specified memory block.  The image will be downsamples to the colour pallette in use at the time.  Supported formats are JPG and BMP.  The format is inferred from the file extension, which must be either `.jpg` or `.bmp`.

### Syntax

FETCH _block-number_, _filename_

### Example

```
10 SET MODE 40
20 REM Load the picture of an astronaut
30 FETCH 0, "astronaut.jpg"
```

## KEEP

Save the image held in the specified memory block to an image file.  The image format is inferred from the file extension, which must be either `.jpg` or `.bmp`.

### Syntax

KEEP _block-number_, _filename_

### Example

```
10 REM Save the image in block 99 as a bitmap
20 KEEP 99, "mypic.bmp"
```

[< Home](index.md)