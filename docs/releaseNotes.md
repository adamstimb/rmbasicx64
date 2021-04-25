# Release Notes

## 0.01 (alpha)

This is an alpha release so it is far from complete and there are many bugs.  At this stage it is only intended as a proof-of-concept that RM Basic and the RM Nimbus SUBBIOS can be emulated in a standalone Go application using ebiten game engine as the underlying platform.

### Known issues

- FOR and REPEAT statements can break if they are not the first instruction in a line.
- FOR ... NEXT and REPEAT ... UNTIL statements must be on separate lines.
- NEXT does not yet accept a variable name as a parameter.
- There is no scoping of variables at the moment so everything is global.
- Keypad controls are a bit glitchy, particularly in delete mode.
- RESTORE does not yet accept a line number as a parameter.

### New features

- RM Basic user interface
- String-type variables
- Integer-type variables
- Float-type variables
- Variable assignment
- Expression evaluation
- BYE
- RUN 
- GOTO 
- PRINT 
- LIST 
- LOAD 
- SAVE 
- EDIT 
- AUTO 
- INPUT 
- NEW
- GOTO
- FOR/NEXT
- REPEAT/UNTIL
- DATA/READ/RESTORE
- SET BORDER
- SET PAPER
- SET PEN
- SET CURSOR
- SET CURPOS
- SET MODE
- Windows installer

