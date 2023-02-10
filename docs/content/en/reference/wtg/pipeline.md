---
title: "Pipeline"
linkTitle: "Pipeline"
type: docs
weight: 10
description: >-
     This is a description of the pipeline
---

## The pipeline component

A pipeline can be declared as a type of a component.

Example:

```
a - mycomponent
mycomponent - b

mycomponent: {
    evolution: |..|..x|..|..|
    type: pipeline
}
```

### Rendering

![](/wardleyToGo/images/pipeline1.svg)

## Declaration of the components of a pipeline

Considering `P` the main pipeline element.
Each component belonging to the `P` pipeline is attached via a semicolon `:`.

example: `P:a` means `a` is a component of the `P` pipeline

as `a` is a component, it can be configured like any other component. Its evolution is configured like any other component (with the `|.|.|.|.|` syntax)

### Value chain

A pipeline is declared on the value chain:

ex: 

- `P:a - P2`: `a` is linked to `P2`
- `P = P2:b`: `P` is linked to `b`
- `P:a - P2:b`: `a` is linked to `b`

_Note_ when a pipeline is declared, the `type` of the host is set to `Pipeline`

### Rendering

a rectangle is automatically drawn including the lower and upper components

```
anchor - P1
P1:a - P2
P1:c

P1: |...|...x...|...|...|
P2: |...|...x...|...|...|
a: |..x.|......|...|...|
c: |...|......|...|x..|

```

![](/wardleyToGo/images/pipeline2.svg)
