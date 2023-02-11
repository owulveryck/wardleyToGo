---
title: "Demo howtos"
linkTitle: "Demo howtos"
type: docs
description: >-
     How to do various things with the demo
---

The demo runs entirely in the browser and does not store any information on a remote server.

## Sharing elements with the demo

### The sharing

The sharing elements works by encoding the wtg elements in gzipped' base64. Then it is added to the url with the `wtg=` parameter.

example: [https://owulveryck.github.io/wardleyToGo/demo/?wtg=H4sIAAAAAAAAE8tIzcnJV9BVSEnNzQcAdsIHTgwAAAA%3D](https://owulveryck.github.io/wardleyToGo/demo/?wtg=H4sIAAAAAAAAE8tIzcnJV9BVSEnNzQcAdsIHTgwAAAA%3D)

### The gist/GitHub integration

It is possible to reference and render wtg files hosted on github or gist.github.com.

To do so, get the url of the raw file, for example [https://raw.githubusercontent.com/owulveryck/wardleyToGo/main/docs/content/en/illustration.wtg](https://raw.githubusercontent.com/owulveryck/wardleyToGo/main/docs/content/en/illustration.wtg)

Then append it to the demo url with the `url=` param:

exemple:

[https://owulveryck.github.io/wardleyToGo/demo/?url=https://raw.githubusercontent.com/owulveryck/wardleyToGo/main/docs/content/en/illustration.wtg](https://owulveryck.github.io/wardleyToGo/demo/?url=https://raw.githubusercontent.com/owulveryck/wardleyToGo/main/docs/content/en/illustration.wtg)

