---
title: "Setting the evolution"
linkTitle: "evolution"
weight: 100
description: >-
     This page describes the evolution placement.
---

Evolution's syntax is composed of 5 pipes `|` and a `x` which represents the position

example : `|..|..|.x.|..|`

The interval between the pipes represents the stages of evolution:

```
|.........|........|....x...|.......|
  stage1    stage2   stage3  stage4
```

you can add as many dots (`.`) (even zero) as you want. This allows fine-tuning the placement of the component on the evolution

## Inline configuration

It is possible to set the evolution "inline" for a component:

example:

`mycomponent: |..|..|..x..|..|`

## block configuration

If you have decalred a block configuration for the component, you can use the `evolution` keyword like this:

```
mycomponent: {
    evolution: |..|..|..x..|..|
}
```

# Evolution

You can add a `>` to display the evolution of the component.

example:

`mycomponent: |..|..|..x..|.>.|`

