[Home](index.md) - [**Quickstart**](quickstart.md) - [History](history.md) - [Reference](reference.md) - [Releases](releases.md)

# Quickstart

Load the RM BASICx64 application.  The prompt with a flashing cursor means it's ready to receive instructions.

```
:_
```

Enter a program to change the screen mode, change the colours, and print a message.  Putting a line number at the beginning of an instruction tells the interpret to store this as a line of code in memory instead of executing it immediately.

```
:10 SET MODE 40 : REM High-colour mode
:20 SET BORDER 5 : SET PAPER 9 : CLS : REM Dark cyan border and light blue paper
:30 SET PEN 0 : REM Black pen
:40 PRINT "Hello!  This is my first RM Basic program."
:_
```

Use the `LIST` command to review your program.

```
:LIST
10 SET MODE 40 : REM High-colour mode
20 SET BORDER 5 : SET PAPER 9 : CLS : REM Dark cyan border and light blue paper
30 SET PEN 0 : REM Black pen
40 PRINT "Hello!  This is my first RM Basic program in a while."
:_
```

Execute the program with the `RUN` command.

```
:RUN
```

Now save the program so we can load it next time with the `LOAD` command.

```
:SAVE "firstprog"
:_
```

To start again from scratch just wipe the workspace with the `NEW` command.  `LIST` shows there's no longer any program stored.

```
:NEW
:LIST
:_
```

Or try loading some of the example programs.  Use the `DIR` command to list all the programs in your workspace including the examples.  

[< Home](index.md)