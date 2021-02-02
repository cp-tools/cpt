---
title: "Browser"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: true
---

# Browser

All website related functions (fetch, submit, list etc) are executed using an automated, headless browser instance (using the DevTools protocol; controlled using [go-rod/rod](https://github.com/go-rod/rod)). **This setup is important if you wish to use website related features.**

{{< hint info >}}
No data of your browser is modified in this method. See [this](https://github.com/cp-tools/cpt-lib#is-sensitive-data-of-my-browser-at-risk) for more details on the process behind the browser automation.
{{< /hint >}}

Browser settings can be configured, both at the global as well as the module level. Shown below is a snippet of the browser configuration for `Google Chrome` on Linux. (refer [here]({{< relref "/docs/setup#overview" >}}) for details about the different configuration files)

{{< tabs "browser-example" >}}
{{< tab "cpt.yaml" >}}
```yaml
browser:
  binary: /opt/google/chrome/google-chrome
  profile: /home/infinitepro/.config/google-chrome
```
{{< /tab >}}
{{< /tabs >}}


{{< hint warning >}}
Not all browsers are supported. Refer [here](https://github.com/go-rod/rod#q-does-it-support-other-browsers-like-firefox-or-edge) for the list of supported browsers.
{{< /hint >}}

---

1. Navigate to `chrome://version/` or `edge://version/`, depending on your browser.
2. Copy the value of key <u>Executable Path</u> as the *binary*.
3. Copy the value of key <u>Profile Path</u> (strip the suffix `/Default`) as the *profile*.

# Screencast

<script id="asciicast-eNMwl4oJtkjBfz3bJLzuUnS2G" src="https://asciinema.org/a/eNMwl4oJtkjBfz3bJLzuUnS2G.js" async data-rows="10" data-speed="1.5" data-theme="monokai"></script>