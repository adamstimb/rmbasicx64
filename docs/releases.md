[Home](index.md) - [Quickstart](quickstart.md) - [History](history.md) - [Reference](reference.md) - [**Releases**](releases.md)

# Releases

## 0.23

12th December 2021

### Bug fixes

- Fixed automatic indenting of stored programs

### New stuff

- OPEN implemented
- CREATE implemented
- CLOSE implemented
- PRINT, DIR, LIST, etc. accept optional file channel

## 0.22

7th December 2021

### Bug fixes

- READ works with arrays now
- NEXT does not require a variable name
- DIM can dimension multiple arrays in one command
- FOR loops no longer break on repeated runs
- Fixed "Array already dimensioned" on repeated runs
- Corrected READBLOCK and COPYBLOCK syntax
- Fixed minor font rendering glitches
- Cosmetic changes to boot screen and editor screen to be more authentic(ish)

### New stuff

- SET WRITING implemented
- CLG implemented
- PRINT, DIR, LIST, CLS can target a writing area with tidle (~)
- SET DRAWING implemented
- CHR$() implemented (smileys!!)
- PUT implemented
- INPUT implemented
- GLOBAL implemented
- Functions and procedures now accept array references
- SET SOUND implemented
- SET TONE implemented
- SET VOICE implemented
- SET ENVELOPE implemented
- PITCH() implemented
- NOTE implemented

## 0.21

23rd November 2021

### New stuff

- DIR now receives a path argument, e.g. DIR "*.JPG" or DIR "myprogs\"
- CHDIR implemented
- MKDIR implemeneted
- RMDIR implemented
- ERASE implemented
- RENAME implemented
- LOOKUP implemented
- STR$() implemented
- MOD implemented

### Bug fixes

- Fixed broken array indexing for 2 dimensions
- Fixed broken IF statement parsing when condition expression contains an array
- Fixed broken parsing of statements between THEN and ELSE
- Fixed incorrect "End of instruction expected" error after DIM statement

## 0.20

14th November 2021

### New stuff

- ANIMATE Extension for loading, displaying and saving image files:
    - READBLOCK implemented
    - WRITEBLOCK implemented
    - SQUASH implemented
    - ASK BLOCKSIZE implemented
    - COPYBLOCK implemented
    - DELBLOCK implemented
    - CLEARBLOCK implemented
    - FETCH implemented
    - KEEP implemented
- Functions implemented
- Procedures implemented
- Subroutines implemented
- POINTS implemented
- FLOOD implemented
- FILL/FLOOD STYLE implemented
- SET PATTERN implemented
- END implemented
- DATA implemented
- RESTORE implemented
- Arrays implemented (although they can't yet be referenced in procedure or function calls)
- CTRL-B also sends sends BREAK interrupt signal (for keyboards without a Scroll-Lock key)
- Location of workspace folder can now be set by RM_BASICX64_WORKSPACE_DIR env var
- Windows installer asks user for location of workspace folder and sets RM_BASICX64_WORKSPACE_DIR accordingly

### Bugfixes

- Fixed borked drawing near top of screen

See the [Reference](https://adamstimb.github.io/rmbasicx64site/docs/reference.html) for an up-to-date list of implemented commands.

## 0.10A

12th September 2021

### New stuff

- Loads and saves BASIC programs to a workspace folder in the installation directory
- Some example programs are included in the workspace folder 
- Application now loads with RM Nimbus "Welcome" boot sequence; disable this feature with the command `SET CONFIG BOOT FALSE`
- Plumbing for 3-channel sound synthesizer is done but not yet accessable (watch this space)
- Error messages highlight the position in the code where an error was detected
- DIR implemented but it can only list BASIC programs in the workspace folder
- All equality operators implemented
- GOTO implemented (woop-woop!)
- REPEAT ... UNTIL loop implemented
- FOR ... NEXT loop implemented
- EDIT implemented
- RENUMBER implemented
- LOAD/SAVE improved error handling
- PRINT extra features added
- HOME implemented
- MOVE implemented
- GET implemented
- SET CURPOS implemented
- PLOT implemented
- AREA implemented
- LINE implemented
- CIRCLE implemented
- SET COLOUR implemented
- SET/ASK MOUSE implemented

### Bugfixes

- Graphics near top or bottom of the screen no longer get sliced in half
- Fixed major slowdown after many graphics operations
- CPU usage is a little more restrained but still runs quite hot
- "Unknown/command procedure" error is now produced if a statement contains only a solitary variable name

See the [Reference](https://adamstimb.github.io/rmbasicx64site/docs/reference.html) for an up-to-date list of implemented commands.

## 0.01B

21st July 2021

LOAD, SAVE, RUN commands implemented with minimal functionality.  EDIT and AUTO commands not yet implemented but program lines can be added or edited by typing the line number followed by the instruction(s), e.g. `10 SET MODE 40 : PRINT "Hi there!"`

IF...THEN...ELSE implemented.

Logical operators are all bitwise, consistent with RM Basic.

See the [Reference](https://adamstimb.github.io/rmbasicx64site/docs/reference.html) for an up-to-date list of implemented commands.

## 0.01A

19th July 2021

This is a preliminary release just so you can get a snifter of things to come.  The reason for
doing this is simply due to the amount of effort that has gone in to getting the project this
far and not really having anything to show for it.

Development work has been in three slightly overlapping phases: 1. Building an emulation of the Nimbus text and graphics drivers, 2. Building an extendable parser to support RM Basic, and 3. Implementing RM Basic in the parser according to the RM Basic manual.  The past 18 months have been spent on phases 1 and 2.  Despite that being a fair chunk of work, there's actually almost nothing to share beyond screen shots because it's all "under-the-bonnet stuff".  

So in this release, you can try some simple commands and evaluate some not-so-simple expressions
and even experience some unspecific syntax errors!

See the [Reference](https://adamstimb.github.io/rmbasicx64site/docs/reference.html) for an up-to-date list of implemented commands.

[< Home](index.md)