# cpt 
[![Go Report Card](https://goreportcard.com/badge/github.com/cp-tools/cpt-lib)](https://goreportcard.com/report/github.com/cp-tools/cpt-lib) ![GitHub](https://img.shields.io/github/license/cp-tools/cpt)

Short for competitive programming tool, `cpt` is an extensively configurable, feature rich, yet lightweight command line tool, to automate the mundane stuff in competitive coding (don't we hate repetitive tasks?)

> Built by CP'ers, built for CP'ers!

Written in [GO](https://golang.org), the compiled executable is cross platform, standalone (no external dependencies) and minisculine in size.

Don't forget to star :star:the project if you found it useful. :smile:

# Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Usage](#usage)
- [Contributing](#contributing)
- [FAQ](#faq)

# Overview

Started as a casual fork of [xalanq/cf-tool](https://github.com/xalanq/cf-tool), and with heavy inspiration from [online-judge-tools/oj](https://github.com/online-judge-tools/oj), this project was brought to life. The project follows the *[KISS](https://en.wikipedia.org/wiki/KISS_principle)* principle, making the user experience as smooth as possible.

Uses [cp-tools/cpt-lib](https://github.com/cp-tools/cpt-lib) as the backend to handle and process website related tasks. **The backend uses the currently active browser session, doing away with any account setup related hassles.**

A non-exhaustive list of notable features of this project are:

- Fetch sample test cases of a problem from the website.
- Manage code templates and create new code files with ease.
- Compile, run and validate solution code against sample tests.
- **Multitude of default checkers, with the ability to use your own**.
- Submit solution to server, and view real time verdict status of the submission.
- List contests, submissions and other tabulated data directly on the terminal.
- **Provides the ability to dynamically auto complete shell commands.**
- Flexible configurations, at the global, sub-module and local directory level.
- **Ability to automatically upgrade to the latest version, once installed.**
- and much, much more...

This list is severely incomplete, and the best way to acquaint yourself with all the available features is by downloading and using the tool (refer instructions below). **You'll fall in love with it, I promise! :heart_eyes:**

# Getting Started

*:warning: Comprehensive documentation is underway (coming soon). :warning:*
**:mega: Please contribute in making the documentation better!**

## Installation

*The following instructions are only for first time users. Existing users can use `cpt upgrade` to upgrade the tool to the latest GitHub released version.*

Download the tarball corresponding to your system from the [latest release](https://github.com/cp-tools/cpt/releases/latest) page. Extract the binary file (executable) to a suitable folder registered as a system path, to run the executable from any directory.

Alternatively, you can compile the executable from source.

```bash
git clone https://github.com/cp-tools/cpt.git && cd cpt
go build -ldflags "-s -w"
```

---

You will also require *checkers* to test (with the command `cpt test`) your solution code against. A set of default *testlib* checkers (maintained at [cp-tools/cpt-checker](https://github.com/cp-tools/cpt-checker)) can be downloaded and installed using the command `cpt upgrade -mc`. 

## Usage

Run the following command to bring up the main help menu:

```bash
cpt --help
```

```
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

Each sub command has a dedicated help menu, that can be accessed in a similar fashion.

```bash
cpt codeforces --help
```
```
Functions exclusive to codeforces

Usage:
  cpt codeforces [command]

Aliases:
  codeforces, cf

Available Commands:
  config      Configure codeforces settings
  fetch       Fetch and save problem tests
  list        Lists specified data in tabular form
  open        Open required page in default browser
  pull        Pulls submissions to local storage
  submit      Submit problem solution to judge

Flags:
  -h, --help   help for codeforces

Use "cpt codeforces [command] --help" for more information about a command.
```

Many commands also have command flags, which can also be viewed in the help menu.

```bash
cpt test --help
```
```
Run code file against sample tests

Usage:
  cpt test [flags]

Flags:
  -c, --checker string       testlib checker to use (default "lcmp")
  -f, --file string          code file to run tests on
  -h, --help                 help for test
  -m, --mode string          mode to run tests on (default "j")
  -t, --timelimit duration   timelimit per test (default 2s)
```

# FAQ

#### What is the minimal setup needed to access all core features?

There are 2 configurations that need to be done, in order to use all functionalities.

- **Headless browser** - Run `cpt config` - select `browser - set headless browser` - follow the instructions to configure the browser data (and cookies) to use.
- **Template** - Run `cpt config` - select `template - add new` - follow the template generation wizard to complete setting up the template. Also, to map the template to a language on a website, run `cpt <website> config` - select `template - set language` - select the corresponding template and set the language name (on the website) that it corresponds to.

With these configurations done, all provided functionalities are usable.