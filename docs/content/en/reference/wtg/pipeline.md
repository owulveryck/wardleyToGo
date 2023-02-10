---
title: "Pipeline"
linkTitle: "Pipeline"
type: docs
weight: 10
description: >-
     This is a description of the pipeline
---

## Declaration

Considering `P` the main pipeline element.
Each component belonging to the `P` pipeline is attached via a semicolon `:`.

example: `P:a` means `a` is a component of the `P` pipeline

as `a` is a component, it can be configured like any other component. Its evolution is configured like any other component (with the `|.|.|.|.|` syntax)

## Value chain

A pipeline is declared on the value chain:

ex: 

- `P:a - P2`: `a` is linked to `P2``
- `P = P2:b`: `P` is linked to `b`
- `P:a - P2:b`: `a` is linked to `b`

_Note_ when a pipeline is declared, the `type` of the host is set to `Pipeline`
