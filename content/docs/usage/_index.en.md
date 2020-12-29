---
title: "Usage"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: true
---

# Usage

{{< hint info >}}
Refer the corresponding subsections (in the side menu) for usage examples and explanations. This section deals with terminologies and usage of command line tools in general.
{{< /hint >}}

Running the command `cpt --help` produces the following.
```text
Lightweight cli tool for competitive programming!

Usage:
  cpt [command]

Available Commands:
  codeforces  Functions exclusive to codeforces
  config      Configure global settings
  generate    Create file using template
  help        Help about any command
  test        Run code file against sample tests
  upgrade     Upgrade cli to latest version

Flags:
  -h, --help      help for cpt
  -v, --version   version for cpt

Use "cpt [command] --help" for more information about a command.
```

# Anatomy

A typical command line interface has the following usage pattern:
```text
APPNAME COMMAND ARGUMENT --FLAG
```

**Commands** represent actions, **Arguments** are things and **Flags** are modifiers of those actions.

In the above help message, 'config' is a command, and 'version' is a flag.

## Command

A command is the central point of the application. Each interaction that the application supports will be encapsulated within a command. A command can have children commands and can optionally run an action.

In the following example,
```text
cpt codeforces config
```
'codeforces' is a command, and 'config' is a subcommand of 'codeforces'.

## Flag

A command line flag (also known as an option or switch) modifies the operation of a command. They are preceeded by two hyphens, and in the case of short flags (containing only a single character) they are preceeded by a single hyphen.

Flags are non-positional, that is, they can be specified anywhere **after the command** is called. Flags may take an optional value (seperated by a equal sign), or act as a switch (boolean flag). 

A valid flag usage example is given below.
```text
cpt test --mode=i -f=test.cpp
```

## Argument

A command line argument is an item of information provided to a program when it is started. A program can have many command line arguments that identify sources or destinations of information, or that alter the operation of the program.

```text
cpt codeforces fetch 1123 c
```
Here, '1123' and 'c' are positional arguments.

