---
title: "Setup"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: true
---

# Setup

CPT is pre-configured, to match the requirements of a majority of users. Two configurations however, need to be manually set by the user - [browser]({{< relref "browser" >}}) and [template]({{< relref "template" >}}) - to enable many core functionalities.

Refer the subsections (from the side menu) for respective configuration instructions.

# Overview

All configuration data is stored in yaml files, at `CONFIG_DIR/cp-tools/cpt/`, where value of `CONFIG_DIR` is determined by [`os.UserConfigDir()`](https://golang.org/pkg/os/#UserConfigDir)

The configurations are broken down by submodule, and are parsed in the following order:

{{< mermaid class="text-center" >}}
graph LR
    id1(Global)-->
    id2(Checker)-->
    id3(Submodule)-->
    id3-->
    id4(Local)
{{< /mermaid >}}

- Checker configurations are present at `CONFIG_DIR/cp-tools/cpt-checker/`.
- Submodules are recursively loaded, in the order they are called.
- Local configurations are loaded from the file `meta.yaml` in the **current directory**.

If the same key is present in multiple configurations, the value from the *last* parsed configuration file containing the key is considered.
