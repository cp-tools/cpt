# cpt

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cp-tools/cpt) ![GitHub last commit](https://img.shields.io/github/last-commit/cp-tools/cpt) ![GitHub issues](https://img.shields.io/github/issues/cp-tools/cpt) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cp-tools/cpt) [![Go Report Card](https://goreportcard.com/badge/github.com/cp-tools/cpt)](https://goreportcard.com/report/github.com/cp-tools/cpt) ![GitHub](https://img.shields.io/github/license/cp-tools/cpt)

Competitive Programming Tool (cpt) is an amazingly versatile, lightweight, feature rich command line tool, to automate all boring stuff of competitive coding (we do hate repetitive stuff, don't we?)

> A little venomous advantage doesn't hurt, does it?

Written using [golang](https://golang.org/), enabling compilation to native machine binary (no dependency configurations!) while keeping the size of the executable minuscule (~ 10 megabytes).

## Table of Contents

- [About the Project]()
  - [Overview]()
  - [Built With]()
- [Getting Started]()
  - [Installation]()
  - [Shell Completions]()
- [Usage]()
- [FAQ]()
- [Contributing]()
- [Acknowledgements]()

## About the Project

### Overview

This project started with a casual fork of [xalanq/cf-tool](https://github.com/xalanq/cf-tool), soon discovering the vast scope for improvement, as no tool as comprehensive as it existed - until now.

This project incorporates the *KISS* (keep it simple, stupid) principle, achieving exactly what we wish to eliminate - boring, repetitive actions. However, effort has been made to polish the console-interface, to enhance the end user experience.

*Some notable features are:*

- Fetch sample tests from problems.
- Compile, run and validate solution.
- **Custom checkers for testing solutions.**
- Submit solutions directly to the server.
- Watch dynamic status of submissions
- Create, modify and manage templates.
- Pull public submissions of any user.
- **List and register for contests directly.**
- **Shell auto-completion support.**
- Open problem page on browser directly.
- **Self-upgrading capability.**
- and many more...



Non-exclusive feature support table for various websites is given below:

| Site                                 | Login              | Fetch Samples      | Submit Solution    | Pull Submissions   |
| ------------------------------------ | ------------------ | ------------------ | ------------------ | ------------------ |
| [codeforces](https://codeforces.com) | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| [atcoder](https://atcoder.jp)        |                    |                    |                    |                    |

### Built With

Huge thanks to the following libraries, which have been of tremendous help to this project:

- [spf13/cobra](https://github.com/spf13/cobra) - The versatility of the command line interface is only possible due to this library.
- [spf13/viper](https://github.com/spf13/viper) - Configuration file management is a breeze, all thanks to this module.
- [AlecAivazis/survey](https://github.com/AlecAivazis/survey) - The clean, interactive-prompt module; definitely brightened this project.
- [cp-tools/cpt-lib](https://github.com/cp-tools/cpt-lib) - Is it a crime to thank myself for the back-end module? The actual website parsing and processing is handled by this library!

Also, huge thanks to the many other smaller libraries (too many to list!) that make this project complete.

## Getting Started

### Installation

If you already run an older version of the tool, simply run `cpt upgrade` and select the latest version (or any version) you wish to install.

First time users can download the latest, compiled binary executable from [here](https://github.com/cp-tools/cpt/releases/latest). Move the executable to any system path folder, to access it (on the terminal) from any directory.

Alternatively, you can also compile the executable from source.

```bash
git clone https://github.com/infixint943/cf.git
cd cf/
go build -ldflags "-s -w"
```

Once installed, you can install *checkers* (for `cpt test`) with the command `cpt upgrade --checkers`.

### Shell Completions

To generate shell completions, run `cpt config` as admin, and select your shell preference from the menu after selecting **Generate tab auto completions**.

*Note:* Windows Powershell users may need to follow additional steps to source the `.ps1` profile that will be created in the current directory. Also, since Powershell is retarded (as compared to `bash`/`zsh`), some features (like dynamic tab completion) are not present. You may consider installing a suitable [WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10) and using [terminal](https://github.com/microsoft/terminal) with `bash`/`zsh` to make use of all features available.

## Usage

Enough said already! Here's a GIF to present the tool in action (made using [terminalizer]())



  

