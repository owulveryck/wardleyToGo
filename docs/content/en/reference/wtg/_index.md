---
title: "The wtg language"
linkTitle: "The wtg language"
type: docs
description: >-
     This is a description of the wtg language
---

`wtg` (wardleyToGo) is a descriptive language used to build Wardley Maps.

This language seperate the construction of the value chain and the position of the components regarding their evolution.

It allows to build a map in two steps:

- you build the value chain without taking care of the placement of the components. 
- each component can then be configured independently. You can set their evolutions, type and even colors

The parser has a placement algorithm and it takes care of the position of the component on the vertical axis.

Components that are not configured are spread on the map on the horizontal axis and their color indicates that they are not configured yet
