# cpt

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cp-tools/cpt) ![GitHub last commit](https://img.shields.io/github/last-commit/cp-tools/cpt) ![GitHub issues](https://img.shields.io/github/issues/cp-tools/cpt) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cp-tools/cpt) [![Go Report Card](https://goreportcard.com/badge/github.com/cp-tools/cpt)](https://goreportcard.com/report/github.com/cp-tools/cpt) ![GitHub](https://img.shields.io/github/license/cp-tools/cpt)

Competitive Programming Tool (cpt) is an amazingly versatile, lightweight, feature rich command line tool, to automate all boring stuff of competitive coding (we do hate repetitive stuff, don't we?)

> A little venomous advantage doesn't hurt, does it?

Written using [golang](https://golang.org/), enabling compilation to native machine binary (no dependency configurations!) while keeping the size of the executable minuscule (~ 10 megabytes).

> Built by CP'ers, built for CP'ers!

## Table of Contents

- [About the Project](#About-the-Project)
  - [Overview](#Overview)
  - [Built With](#Built-With)
- [Getting Started](#Getting-Started)
  - [Installation](#Installation)
  - [Shell Completions](#Shell-Completions)
- [Usage](#Usage)
- [FAQ](#FAQ)
- [Contributing](#Contributing)

## About the Project

### Overview

This project initially started with a casual fork of [xalanq/cf-tool](https://github.com/xalanq/cf-tool). Discovering the vast scope for improvement, this project was brought to life.
This project incorporates the *KISS* (keep it simple, stupid) principle, achieving exactly what we wish to eliminate - boring, repetitive actions. However, effort has been made to polish the console-interface, to enhance the end user experience.

This application serves as a handy tool for competitive programmers who love using the terminal (we all do, don't we?) to accomplish as many repetitive tasks as possible.

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

What are you waiting for? Go ahead and download it. You'll love it, I promise!

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
git clone https://github.com/cp-tools/cpt.git
cd cpt/
go build -ldflags "-s -w"
```

Once installed, you can install *checkers* (for `cpt test`) with the command `cpt upgrade --checkers`.

### Shell Completions

To generate shell completions, run `cpt config` as admin, and select your shell preference from the menu after selecting **Generate tab auto completions**.

*Note:* Windows Powershell users may need to follow additional steps to source the `.ps1` profile that will be created in the current directory. Also, since Powershell is retarded (as compared to `bash`/`zsh`), some features (like dynamic tab completion) are not present. You may consider installing a suitable [WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10) and using [terminal](https://github.com/microsoft/terminal) with `bash`/`zsh` to make use of all features available.

## Usage

Enough said already! Here's a GIF to present the tool in action (made using [asciinema](https://asciinema.org)!)

![demo](assets/demo.gif)

For complete documentation, head to the [wiki]() page.

## FAQ

#### How do I install and use this tool?

For first time users, download the executable corresponding to your Operating System and add it to system PATH. This way, you will be able to access the tool through the terminal, from any directory.

If you wish to upgrade the tool to the latest version, simply run `cpt upgrade` and select the desired version to upgrade to.

#### How do I configure my login credentials? Are they securely stored?

For module (`cf`, `atc` etc) specific configurations, use the `config` sub command in the corresponding module.

Yes, none of your data is accessible by us! Your passwords are encrypted using *AES* and saved locally. However do note that, your password is **NOT** secure. The encryption is simply to prevent others from reading the password when stored as plain text (if someone got access to your device, you can consider your account compromised in any case, and the encryption won't help much.)

#### Can I enable dynamic tab completion support? How?

Yes, you can! However, some features (like custom tab completion) are not accessible unless you use `zsh` or `bash`. Refer [Shell Completions](#Shell-Completions) for complete information.

## Contributing

Here are the ways through which you can contribute to this project:

- **Create a pull request** - Yes, the best way to contribute would by helping us make this tool better and better. Simply fork the project, make additions/changes, and create a pull request! It isn't that hard, trust me!
- **Star this project** - Yep, you heard it right. I'm an insatiable beast, with a never quenching thirst for GitHub stars. Jokes aside, I'm a normal teenager, utilising the few hours of spare time that I have, to give something back to the community. So, your stars let me know that the community likes this project and it only motivates me to do better!
- **Helping in documentation** - Yes, documentation sucks (a lot)! However, the existence of a well written documentation is very essential. Help the community (and get your name etched forever in the contributors list) by improvising the [wiki pages](https://github.com/cp-tools/cpt.wiki).