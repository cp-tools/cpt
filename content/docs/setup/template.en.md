---
title: "Template"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: true
---

# Template

Templates are core components of cpt, and configuration of the same is required to use two of the most powerful features of cpt - `test` and `submit`.

Templates are stored under global settings. Shown below is a snippet of the template configuration for the language `C++` on Linux (refer [here]({{< relref "/docs/setup#overview" >}}) for details about the different configuration files).

{{< tabs "template-example" >}}
{{< tab "cpt.yaml" >}}
```yaml
template:
  cpp:
    codeFile: /home/admin/Documents/template.cpp
    preScript: g++ "{{.file}}"
    runScript: ./a.out
```
{{< /tab >}}

{{< tab "codeforces.yaml" >}}
```yaml
template:
  cpp:
    language: GNU G++17 7.3.0
```
{{< /tab >}}
{{< /tabs >}}

## Alias

The alias is an *unique* name given to every template, to differentiate it from other configured templates. In the template above, `cpp` is the alias value.

The alias should be *alpha-numeric* and not contain any whitespaces. Alias names are case sensitive, and must be different from other aliases.

## Code File

Code file/template code is the skeleton code you use in every solution. Save this template code (it can be empty too!) in a file and ensure the file stays in one location. The code from this file is duplicated when `cpt generate` is run.

Templates also support variable placeholders, parsed using the [text/template](https://golang.org/pkg/text/template/) package of golang. Presently, the following *generic* placeholders are supported:
```
{{.date}} - Current date, in dd.mm.yyyy format
{{.time}} - Current time, in hh:mm format
```
Dynamic placeholders are supported too, with values fed from the local configuration (`meta.yaml` in the current directory).

Presented below is an example showing the local configuration, template file and generated code.
{{< tabs "template-generate" >}}
{{< tab "meta.yaml" >}}
```yaml
problem:
    name: Binary Table (Hard Version)
    memoryLimit: 256 megabytes
    inputStream: standard input
    outputStream: standard output
customFlag:
    addCredits: true
```
{{< /tab >}}

{{< tab "template.cpp" >}}
```cpp
/*
Date: {{.date}} | Time: {{.time}}
{{if .customFlag.addCredits -}} Credits: cp-tools {{- end}}

{{if .problem.name -}} Problem: {{.problem.name}} {{- end}}
{{if .problem.timeLimit -}} Time limit: {{.problem.timeLimit}} {{- end}}
*/

#include<bits/stdc++.h>
using namespace std;

int main(){
    // Enter code here...
    return 0;
}
```
{{< /tab >}}

{{< tab "gen.cpp" >}}
```cpp
/*
Date: 18.11.2020 | Time: 13:48
Credits: cp-tools

Problem: Binary Table (Hard Version)

*/

#include<bits/stdc++.h>
using namespace std;

int main(){
    // Enter code here...
    return 0;
}
```
{{< /tab >}}
{{< /tabs >}}

## Test Scripts

Test scripts are required *exclusively* by `cpt test` - the testing module - and specify the compilation *(prescript)* and execution *(runscript)* commands to be executed.

Test scripts also support the following generic placeholders:
```
{{.file}}               - The absolute path to the solution code file
{{.fileBasename}}       - The name of the file, without the path
{{.fileBasenameNoExt}}  - The name of the file, without path and extension
```

### Prescript

This script is executed *before the testing*. The intended usecase is to compile your solution file. Interpreted languages like Python and Golang may leave this blank.

<!--Add folded section of commands by language-->

### Runscript

This script is executed *once per test case*. This command is meant to run the compiled binary (for compiled languages) or execute the solution file (for interpreted languages).

{{< hint warning >}}
Unless flags `input-stream` and `output-stream` are used while testing, the test case *input/output* is read from/written to the *standard input/standard output*.
{{< /hint >}}
For problems that require *file input/file output*, use the command line flags `input-stream` and `output-stream` respectively.

<!--Add folded section of commands by language-->

## Language

{{< hint info >}}
The value of this key differs from site to site, and is thus configured at the module (website) level, rather than the global level.
{{< /hint >}}

The value specified here, is the language selected while submitting the code file (corresponding to this template) to the remote judge.

# Screencast

<script id="asciicast-F3qk7qtQy5AKhAhL7Oj6xbmQP" src="https://asciinema.org/a/F3qk7qtQy5AKhAhL7Oj6xbmQP.js" async data-rows="20" data-speed="1.5" data-theme="monokai"></script>