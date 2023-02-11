---
title: "label position"
linkTitle: "label position"
weight: 100
type: docs
description: >-
     Tweak the label position.
---

On the map, labels are placed automatically in a pseudo-smart way.
It is possible to change the position of a label by using the `label` keywork in a configuration block.

The allowed values are: `N`,`S`,`E`,`W`,`NE`,`NW`,`SE`,SW` (for north, south, east, west, north-east, north-west, south-east, south-west)

example:

```
N - NE
NE - E
E - SE
SE - S
S - SW
SW - W
W - NW

N: {
    label: N
}

NE: {
    label: NE
}

E: {
    label: E
}

SE: {
    label: SE
}

S: {
    label: S
}

SW: {
    label: SW
}

W: {
    label: W
}

NW: {
    label: NW
}
```

![](/wardleyToGo/images/labelplacement.svg)
