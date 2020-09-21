# Design conventions

General conventions followed when developing in eksctl:

- general
  - never use `panic` in your code, if a library call panics - think twice
  - never prompt interactively, assume automation as prime use-case
  - never animate output
  - avoid foreign output message

- log message style
  - keep it neutral, don't address the user as a human (i.e. instead of "waiting for X, please be patient" just say "waiting for X")
  - don't suggest how long something is going to take, it's always going to be wrong
  - must always begin a message with lower-case letter, it's not a sentence
      - if the first word is a name, it's fine to use capital letter e.g. '[i]  Kubernetes welcomes you'
  - do not add spaces or tabs before the first character
  - use '%q' for names of resources
  - only reference flag names or API fields when absolutely certain of the use-case
  - use single quotes when suggesting use to run a command
    - only log these from top-level command handlers, not from a library
    - only suggest simple commands, no shell one-liners
    - use long form of the flags in suggestions
    - if referring to eksctl, try to make the command as full as possible with cluster name and region
  - do not use any punctuation other than column, dash or comma, and only when absolutely essential
    - i.e. no dots, ellipsis or semicolons
  - do not write long log messages, try to break up into two if extra information is needed

- use `errors.Wrapf` to wrap errors
  - avoid wrapping errors that are already meaningful

- we use Kris Nova's logger, this may change in the future and we should probably abstract it

- we use Cobra (for better or worse)
    - we put all shared things in `pkg/ctl/cmdutils`, familiarise yourself with it
      - `cmdutils.Cmd` wraps most of the things that are shared between various commands
      - we have customisation on top of Cobra that allow us to group flags in usage summary
      - this package is all about CLI use-cases, it shouldn't be consumed by libraries
    - don't use `cobra.MarkFlagRequired`
    - `fs.MarkHidden` is evil, but sometimes we have to use it

- CLI
  - commands follows a similar convention as `kubectl`: `eksctl <verb> <resource> <flags...>` e.g.: `eksctl create nodegroup`, `eksctl enable profile`
  - eksctl does not use arguments except when used for the name of the resource (e.g.: `eksctl create nodegroup ng-1`). This is what we call `nameArg` in the code. The reason for this is that users get very easily confused when arguments are used and what their meaning is, and relative position with respect to flags. The name of the resource will also have the option to pass it as part of the flag called `--name` (e.g.: `eksctl create nodegroup --name ng-1`).
  - when the main resource is not the cluster, the cluster will be set in the `--cluster` flag (e.g.: `eksctl create nodegroup --cluster clus1 ng-1`)
  - for boolean flags, use `*bool` so that eksctl can know if a flag was explicitly enabled or disabled by the user
  - add unit tests for the flags and config file parsing. These are tedious to manual test. See [profile_test.go](pkg/ctl/enable/profile_test.go) as an example
  - keep the number of CLI flags small. Not all options should have flags, but rather have a way to use them through the config file. In general, basic use cases can be covered by flags but more advanced options should only live in the config file.
  
