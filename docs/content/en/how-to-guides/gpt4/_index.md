---
title: "ChatGPT Experimental prompt"
linkTitle: "GPT-4 Prompt"
weight: 100
type: docs
description: >-
     This is an experiment prompt for GPT-4
---

This is an attempt to build a prompt to use GPT-4 to build wardley maps.
Do not hesitate to submit PR to improve this.

```
I want you to be my assistant to create and to understand Wardley Maps.
You will use the WTG language as explained here as well as you knowledge to challenge the existing maps and to create new ones.

Explanation of the WTG language:

The basic element is a component. A component can be a practice, an activity, a product or a data-set.

- a practice is something that is done repeatedly in order to acquire or polish a skill.
- an activity is An act that requires use of resources or time.
- a product is an asset that solves a problem for a user
- a data-set is a a digital representation of some knowledge

A component is represented in WTG by a set of characters separated by one or none whitespace (the whitespaces are used to seperate words if needed, for example a product can be "cup of tea").

comment can be added with // for single line comments, or /* */ for multilines comments

How to create a map:

1. create the value chain

To create a value chain, link components in pairs with a set of dashes (-).
One pair per line.
The component on the left requires the component on the right. 
The more dashes, the less valuable the component on the right is.
for example:
user - cup of tea
cup of tea - cup
cup of tea -- tea

means that a user needs a cup of tea, and a cup of tea needs tea and a cup, but the cup is more valuable to the user than the tea.

another example is:
cup of tea - hot water
hot water -- kettle
hot water - water
kettle ----- power

meaning that kettle and water are required to do hot water, but water is more visible (valuable) than the kettle.
and from the point of view of the cup of tea, the power is not very visible (one dash from hot water, two dashed from kettle and five dashes from power)
two components 

For one component, avoid setting the same visibility for two dependencies.
For example do not do this:

cup of tea - cup
cup of tea - water

but decide which is more visible, like this:

cup of tea - cup
cup of tea -- water

2. place the components on the evolution axis

First, figure out where the component should be placed on the evolution by using Simon Wardley's cheat sheet.
See which characteristics and properties are matching the components, and do a an average to evalute the stage of evolution of the component.

To add the evolution in WTG, use this syntax: 5 pipes (|), some dots (.) and an x.
You place an 'x' to set the component placement on the evolution axis. 

If the component is in stage 1, the x will be between the first and second pipe (the make it work phase).
If the component is in stage 2, the x will be between the second and the third pipe (the make it deliver phase).
If the component is in stage 3, the x will be between the third and fourth pipe (the make it efficient phase).
If the component is in stage 4, the x will be between the fourth and the fifth pipe (the make it ambiant phase).

Stages depends of the type of the component.
- for a practice stage 1 is novel, stage 2 is emerging, stage 3 is good and stage 4 is best
- for an activity or a product: stage 1 is genesis, stage 2 is custom, stage 3 is product and stage 4 is commodity
- for data-set stage 1 is unmodelled, stage 2 is divergent, stage 3 is convergent and stage 4 is modelled.

To fine place the component in its stage, place the component between dots. the number of characters (. and x) represents the length of a segment. 
The position of the x in the segment is the position within the stage.

for example, ...x. places the component at 80% of its evolution stage (it is at place 4 over 5 characters)

example:
// mycomponent is a custom system because...
mycomponent: |.....|..x..|...|..|

means that mycomponent is in stage 2 (between pipes 2 and 3), and there are 5 characters total, it is in the middle, so 50% of the stage 2.

// agility is a best practice nowadays. As seen on the cheat sheet, most of the publication focus on use.
agility: |...|...|...|x...|

means agility, a practice, is a best practice (stage 4)

// power is a commodity. Its ubiquitous	and stable as seen in the cheat sheet
power: |...|...|...|....x..| 

is a product in stage 4, a commodity.

// chatGPT is a product in genesis. OpenAI made it work, and is not yet efficient.
chatGPT: |...x|...|...|...| 

means that chatGPT is at the end of the genesis.

Don't forget to add some comments in front of each component definition to explain why it is in this stage of evolution

Questions: 

Can you provide an example of a complete Wardley Map using the WTG language?
How do you determine the visibility of a component within the value chain?
Can you provide more information on Simon Wardley's cheat sheet and how to use it effectively when creating a map?
```

Example: 

> create the teashop example

> Here's an example of a Wardley Map for a teashop using the WTG language:
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
