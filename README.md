# ppi

Generate `draw.Image` from PSD (PhotoShop Document) and PFV (PSDToolFavorites)

## Usage
```go
package main

import (
	"fmt"
	"image/png"
	"os"
	"strings"

	"github.com/oov/psd"
	"github.com/raa0121/pfv"
	"github.com/raa0121/ppi"
)

func main() {
	s := `
[PSDToolFavorites-v1]
root-name/お気に入り
faview-mode/1

//Heart
マーク/ハート
アクセサリ/オタマン帽普
眉/*普通
目/*普通
口/*笑い
本体/*普段着
手/*手
ツインテール/*通常

//hoge/Heart2
マーク/ハート
アクセサリ/オタマンベレー帽
眉/*普通
目/*普通
口/*笑い
本体/*普段着
手/*手
ツインテール/*通常

`
	config := pfv.Decode(s)

	file, err := os.Open("image.psd")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	img, _, err := psd.Decode(file, &psd.DecodeOptions{})
	imgs := ppi.CreateImage(img, config)
	for _, v := range imgs {
		out, err := os.Create(fmt.Sprintf("%s.png", strings.Replace(v.Name, "/", "_", -1)))
		if err != nil {
			panic(err)
		}
		defer out.Close()
		err := png.Encode(out, v.Image)
		if err != nil {
			panic(err)
		}
	}
}
```
## Requirements
go 1.16 or later

## Installation
```
go get github.com/raa0121/ppi
```

## License
Apache License 2.0

## Author
raa0121 <raa0121@gmail.com>
