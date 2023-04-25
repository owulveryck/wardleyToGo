---
title: "ChatGPT Experimental prompt"
linkTitle: "GPT-4 Prompt"
weight: 100
description: >-
     This is an experiment prompt for GPT-4
---

This is an attempt to build a prompt to use GPT-4 to build wardley maps.
Do not hesitate to submit PR to improve this.

```
Here is how to express a Wardley Map with a language called wtg.
The basic element is a component. A component can be a practice, an activity, a product or a data-set.
A component is represented in wtg by a set of characters separated by one or none whitespace.

## create the value chain

To create a value chain, link components in pairs with a set of dashes (-).
One pair per line.
The component on the left requires the component on the right. 
The more dashes, the less valuable the component on the right is.
for example:

user - cup of tea
cup of tea - cup
cup of tea -- tea

means that a user needs a cup of tea, and a cup of tea needs tea and a cup, but the cup is more valuable to the user than the tea.

## place the components on the evolution axis

To add the evolution in wtg, use this syntax: 5 pipes (|), some dots (.) and an x.
You place an 'x' to set the component placement on the evolution axis. 
If the component is in stage 1, the x will be between the first and second pipe.
If the component is in stage 2, the x will be between the second and the third pipe.
If the component is in stage 3, the x will be between the third and fourth pipe.
If the component is in stage 4, the x will be between the fourth and the fifth pipe.

To fine place the component in its stage, place the component between dots. the number of characters (. and x) represents the length of a segment. 
The position of the x in the segment is the position within the stage.
for example, ...x. places the component at 80% of its evolution stage (it is at place 4 over 5 characters)

for example:
mycomponent: |.....|..x..|...|..|

means that mycomponent is in stage 2 (between pipes 2 and 3), and there are 5 characters total, it is in the middle, so 50% of the stage 2.

first, figure out where the component should be placed on the evolution by using Simon Wardley's cheat sheet.

comments can be added with // for single line comments, or /* */ for multilines comments

in résumé, to build a map:
- set the visibility when buildind a map.
- analyze the value of a component with regard to the top component and add as many dashes as required
- put some comments before each evolution to explain the component's evolution with the cheat sheet
- set the evolution for every components.
- do not set evolution of some components if they are not part of the value chain

just give a short answer if you have understood.
```

Example: 

> create the teashop example

> Here's an example of a Wardley Map for a teashop using the wtg language:
```
/* Components */
user
cup of tea
cup
tea
water
kettle
electricity

/* Value chain */
user - cup of tea
cup of tea - cup
cup of tea -- tea
cup of tea --- water
water ---- kettle
kettle ----- electricity

/* Evolution */
// user: no evolution (not part of the value chain)

// cup of tea: custom-built (stage 1)
// Simon Wardley's cheat sheet: Novel and new, few competitors, high margin, constant change
cup of tea: |x...|...|...|...|

// cup: product (stage 3)
// Simon Wardley's cheat sheet: commodity, low margin, stable, many competitors
cup: |...|...|x..|...|

// tea: product (stage 3)
// Simon Wardley's cheat sheet: commodity, low margin, stable, many competitors
tea: |...|...|x..|...|

// water: utility (stage 4)
// Simon Wardley's cheat sheet: undifferentiated, no competitors, stable, essential
water: |...|...|...|x..|

// kettle: product (stage 3)
// Simon Wardley's cheat sheet: commodity, low margin, stable, many competitors
kettle: |...|...|x..|...|

// electricity: utility (stage 4)
// Simon Wardley's cheat sheet: undifferentiated, no competitors, stable, essential
electricity: |...|...|...|x..|
```
