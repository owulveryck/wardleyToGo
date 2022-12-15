Given sample.txt with this content

```text
business - cup of tea
public - cup of tea
cup of tea - cup
cup of tea -- tea
cup of tea --- hot water
hot water - water
hot water -- kettle
kettle - power

cup of tea: {
    type: buy
    evolution: |....|....|...x..|.........|
}
water: {
    type: build
    evolution: |....|....|....|....x....|
}
kettle: {
    type: build
    evolution: |....|..x.|....|.........|
}
power: {
    type: outsource
    evolution: |....|....|....x|.........|
}
business: {
    evolution: |....|....|..x.|.......|
}
public: {
    evolution: |....|....|....|.x....|
}
cup: {
    evolution: |....|....|....|.x.......|
}
tea: {
    evolution: |....|....|....|..x......|
}
hot water: {
    evolution: |....|....|....|...x.....|
}
```

When we execute

```shell
cat sample.txt | go run cmd/wtg2svg/main.go > sample.svg
```

The result is

![](sample.svg)
