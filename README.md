[![Go Reference](https://pkg.go.dev/badge/github.com/owulveryck/wardleyToGo.svg)](https://pkg.go.dev/github.com/owulveryck/wardleyToGo)
[![Go](https://github.com/owulveryck/wardleyToGo/actions/workflows/go.yml/badge.svg)](https://github.com/owulveryck/wardleyToGo/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/owulveryck/wardleyToGo/branch/main/graph/badge.svg?token=9BQW1KMGJS)](https://codecov.io/gh/owulveryck/wardleyToGo)

## wardleyToGo

simple tools to play with Wardley maps in Go. The syntax is a copy of [onlinewardleymaps](https://onlinewardleymaps.com/). Note all the DSL is implemented as of today.

A typical flow is to create the map with the online tool and to copy paste it to use with the tool.

## Example

```shell
❯ cat sample.txt | go run . 
```

![sample](sample.svg)


## Usage

### Wardley Map

#### To set the title (not yet implemented)

`title My Wardley Map`

#### To create a component

`component Name [Visibility (Y Axis), Maturity (X Axis)]`

#### To create a market (not yet implemented)

`market Name [Visibility (Y Axis), Maturity (X Axis)]`

#### Inertia - component likely to face resistance to change. (not yet implemented)

`component Name [Visibility (Y Axis), Maturity (X Axis)] inertia`

#### To evolve a component (not yet implemented)

`evole Name (X Axis)`

#### To link components

`Start Component->End Component`

#### To indicate flow (not yet implemented)

`Start Component+<>End Component`

#### To set component as pipeline (not yet implemented)

`pipeline Component Name [X Axis (start), X Axis (end)]`

#### To indicate flow - past components only (not yet implemented)

`Start Component+<End Component`

#### To indicate flow - future components only (not yet implemented)

`Start Component+>End Component`

#### To indicate flow - with label (not yet implemented)

`Start Component+'insert text'>End Component`

#### Pioneers, Settlers, Townplanners area (not yet implemented)

Add areas indicating which type of working approach supports component development

`pioneers [<visibility>, <maturity>, <visibility2>, <maturity2>]`

#### Build, buy, outsource components (not yet implemented)
Highlight a component with a build, buy, or outsource method of execution

* `build <component>`
* `buy <component>`
* `outsource <component>`
* `component Customer [0.9, 0.2] (buy)`
* `component Customer [0.9, 0.2] (build)`
* `component Customer [0.9, 0.2] (outsource)`
* `evolve Customer 0.9 (outsource)`
* `evolve Customer 0.9 (buy)`
* `evolve Customer 0.9 (build)`

#### Link submap to a component (not yet implemented)

Add a reference link to a submap. A component becomes a link to an other Wardley Map

* `submap Component [<visibility>, <maturity>] url(urlName)`
* `url urlName [URL]`
* `submap Website [0.83, 0.50] url(submapUrl)`
* `url submapUrl [https://onlinewardleymaps.com/#clone:qu4VDDQryoZEnuw0ZZ]`

#### Stages of Evolution

Change the stages of evolution labels on the map

* `evolution First->Second->Third->Fourth`
* `evolution Novel->Emerging->Good->Best`

#### Y-Axis Labels (not yet implemented)

Change the text of the y-axis labels

* `y-axis Label->Min->Max`
* `y-axis Value Chain->Invisible->Visible`

#### Add notes (not yet implemented)

Add text to any part of the map

* `note Note Text [0.9, 0.5]`
* `note +future development [0.9, 0.5]`

#### Available styles (not yet implemented)

Change the look and feel of a map

* `style wardley`
* `style handwritten`
* `style colour`

### Team topologies

A couple of additions has been made to add the team topologies shapes:

#### Stream Aligned Team

`streamalignedteam Team Name [<visibility>, <maturity>, <visibility2>, <maturity2>]`

#### Platform Team

`platformteam Team Name [<visibility>, <maturity>, <visibility2>, <maturity2>]`

#### Enabling Team Team

`enablingteam Team Name [<visibility>, <maturity>, <visibility2>, <maturity2>]`
