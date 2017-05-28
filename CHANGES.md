For now, dates (DD/MM/YYYY) are used until ishell gets stable enough to warrant tags.
Attempts will be made to ensure non breaking updates as much as possible.
#### 28/05/2017
* Added `shell.Process(os.Args[1:]...)` for non-interactive execution
*


#### 07/02/2016
Added multiline support to shell mode.

#### 23/01/2016
* Added history support.
* Added tab completion support.
* Added `SetHistoryPath`, `SetMultiPrompt`
* Removed password masks.
* **Breaking Change**: changed definition of `ReadPassword` from `(string)` to `()`
* **Breaking Change**: changed name of `Shell` constructor from `NewShell` to `New`

#### 13/07/2015
* Added `ClearScreen` method.
* Added `clear` to default commands.

#### 12/07/2015:
* Added `PrintCommands`, `Commands` and `ShowPrompt` methods.
* Added default `exit` and `help` commands.
* **Breaking Change**: changed return values of `ReadLine` from `(string, error)` to `string.`
* **Breaking Change**: changed definition of `CmdFunc` from `(cmd string, args []string)` to `(args ...String)` to remove redundant command being passed.
* Added multiline input support.
* Added case insensitive command support.

#### 11/07/2015:
* Initial version.
