---
title: "Tea Shop in wtg"
linkTitle: "Tea Shop in wtg"
type: docs
weight: 10
description: >-
     Your first wardley map with WTG
---

This is the "tea shop" example described step by step in `wtg`

## Value Chain

Imagine you are a tea shop.
The `business` and the `public` need a `cup of tea`; this is described like this:

```
business - cup of tea
public - cup of tea
```

![](/wardleyToGo/images/cupoftea1.svg)

The the `cup of tea` required `tea`, `cup`, and `hot water`.
`tea` is more visible than `cup` which is more visible than `hot water` 

```
cup of tea - tea
cup of tea -- cup
cup of tea --- hot water
```

![](/wardleyToGo/images/cupoftea2.svg)

Adding all the dependencies will lead to this rended:

```
business - cup of tea
public - cup of tea
cup of tea - tea
cup of tea -- cup
cup of tea --- hot water
hot water - water
hot water -- kettle
kettle - power
```

![](/wardleyToGo/images/cupoftea3.svg)

## Evolution

Now that we have built a value chain, we can place the components according to their evolution.
Let's start with the anchors. Business is at the end of stage 3. public is in stage 4.

```
business: |..|..|..x|..|
public: |..|..|..|.x.|
```

![](/wardleyToGo/images/cupoftea4.svg)

configuring all the components leads to this map:

```
business - cup of tea
public - cup of tea
cup of tea - cup
cup of tea -- tea
cup of tea --- hot water
hot water - water
hot water -- kettle
kettle - power

business:   |...|.....|...x.|..........|
public:     |...|.....|.....|..x.......|
cup of tea: |...|.....|..x..|..........|
cup:        |...|.....|.....|.....x....|
tea:        |...|.....|.....|.....x....|
hot water:  |...|.....|.....|....x.....|
water:      |...|.....|.....|.....x....|
kettle:     |...|...x.|..>..|..........|
power:      |...|.....|....x|.....>....|
```

![](/wardleyToGo/images/cupoftea5.svg)
