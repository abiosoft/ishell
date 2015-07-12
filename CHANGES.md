For now, dates (DD/MM/YYYY) are used until ishell gets stable enough to warrant tags.
Attempts will be made to ensure non breaking updates as much as possible.

#### 12/07/2015:
* Added `PrintCommands`, `Commands` and `ShowPrompt` methods.
* Added default "exit" and "help" commands to new shell returned by NewShell().
* **Breaking Change**: changed return values of `ReadLine` from `(string, error)` to `string.`
* **Breaking Change**: changed definition of `CmdFunc` from `(cmd string, args []string)` to `(args ...String)` to cater for redundant command that is passed.
* Added multiline input support.
* Added case insensitive command support.

#### 11/07/2015:
* Initial version.