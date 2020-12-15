---
title: "Installation"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: true
---

# Installation

{{< hint info >}}
The below instructions are for first time users only.
Existing users can upgrade using:
```bash
cpt upgrade -m s # Upgrade cli tool
cpt upgrade -m c # Upgrade checkers
```
{{< /hint >}}

## GitHub Release

Download the tarball corresponding to your system from the [latest release](https://github.com/cp-tools/cpt/releases/latest) page. Extract the executable from the archive and follow the instructions given below, according to your OS.

{{< hint warning >}}
**Linux/MacOS**

Zsh users may need to use `.zshrc` instead of `.bashrc`.
{{< /hint >}}
{{< hint danger >}}
**Windows**

[WSL/WSL2](https://docs.microsoft.com/en-us/windows/wsl/about) is currently not supported.
{{< /hint >}}

{{< tabs "github-releases" >}}
{{< tab "Linux" >}}
Move the executable to the following directory,
```bash
/home/<username>/.local/bin
```
where `<username>` is your system username.

To add the above directory to the system PATH, follow [this](https://askubuntu.com/q/60218/994766). 

{{< /tab >}}

{{< tab "MacOS" >}}
Move the executable to the following directory,
```bash
/usr/local/bin
```
You may require `sudo` access to do this.
{{< /tab >}}

{{< tab "Windows" >}}
Move the executable to the following directory,
```powershell
C:\Program Files\
```

To add the above directory to the system PATH, follow [this](https://www.architectryan.com/2018/03/17/add-to-the-path-on-windows-10/). 
{{< /tab >}}

{{< /tabs >}}

## Source Build

To build from source, execute the following commands.
```bash
git clone https://github.com/cp-tools/cpt.git
cd cpt/
go install -ldflags "-s -w"
```

Ensure that the `$GOBIN` directory is present in the system PATH.

# Verification

{{< hint warning >}}
You may need to re-login to your desktop session, for the command to be picked up from the system PATH.
{{< /hint >}}

A successful installation would output text similar to
```bash
cpt version vX.Y.Z
```
on running the command `cpt --version`.
