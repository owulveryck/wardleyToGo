---
title: "component types"
linkTitle: "types"
weight: 100
type: docs
description: >-
     The component types
---

wtg supports four types of components:

- build
- buy
- outsource
- pipeline

those components are set in a block configuration of the component thanks to the `type` keyword.

Example:

```
build - buy
buy - outsource
outsource - pipeline

build: {
    type: build
}
buy: {
    type: buy
}
outsource: {
    type: outsource
}
pipeline: {
    type: pipeline
}
```

![](/wardleyToGo/images/types.svg)

