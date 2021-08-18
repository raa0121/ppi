package ppi

import (
	"bytes"
	"image"
	"image/draw"

	"github.com/oov/psd"
	"github.com/raa0121/pfv"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

type Image struct {
	Name string
	Image draw.Image
}

func generateLayersName(layer psd.Layer, dirName string) map[string]psd.Layer {
	names := map[string]psd.Layer{}
	layerName := detectLayerNameEncoding([]byte(layer.Name))
	if layer.Folder() {
		for _, l := range layer.Layer {
			tmpMap := map[string]psd.Layer{}
			if dirName == "" {
				tmpMap = generateLayersName(l, layerName)
			} else {
				tmpMap = generateLayersName(l, dirName + "/" + layerName)
			}
			names = merge(names, tmpMap)
		}
	} else {
		if dirName == "" {
			names[layerName] = layer
		} else {
			names[dirName + "/" + layerName] = layer
		}
	}
	return names
}

func merge(m ...map[string]psd.Layer) map[string]psd.Layer {
	ans := make(map[string]psd.Layer, 0)

	for _, c := range m {
		for k, v := range c {
			ans[k] = v
		}
	}
	return ans
}

func CreateImage(img *psd.PSD, conf *pfv.Pfv) []Image {
	output := map[string][]psd.Layer{}
	imgs := []Image{}
	canvas := &image.RGBA{}
	layerNames := map[string]psd.Layer{}

	for _, l := range img.Layer {
		layerNames = merge(layerNames, generateLayersName(l, ""))
	}
	for _, v := range conf.Items {
		for _, vv := range v.Elements {
			for layerName, l := range layerNames {
				if vv.Path == layerName {
					output[v.Name] = append(output[v.Name], l)
				}
			}
		}
	}
	for i, o := range output {
		output[i] = reverse(o)
	}
	for i, o := range output {
		for j, oo := range o {
			if j == 0 {
				canvas = image.NewRGBA(img.Picker.Bounds())
				draw.Draw(canvas, oo.Rect, oo.Picker, oo.Rect.Min, draw.Src)
			} else {
				draw.Draw(canvas, oo.Rect, oo.Picker, oo.Rect.Min, draw.Over)
			}
		}
		im := Image{
			Name: i,
			Image: canvas,
		}
		imgs = append(imgs, im)
	}
	return imgs
}

func detectLayerNameEncoding(body []byte) string {
	encodings := []string{"sjis", "euc-jp", "utf-8"}
	var f []byte
	for _, enc := range encodings {
		ee, _ := charset.Lookup(enc)
		if ee == nil {
			continue
		}
		var buf bytes.Buffer
		ic := transform.NewWriter(&buf, ee.NewDecoder())
		_, err := ic.Write(body)
		if err != nil {
			continue
		}
		err = ic.Close()
		if err != nil {
			continue
		}
		f = buf.Bytes()
		break
	}
	return string(f)
}

func reverse(p []psd.Layer) []psd.Layer {
	reversed := []psd.Layer{}
	for i := range p {
		n := p[len(p)-1-i]
		reversed = append(reversed, n)
	}
	return reversed
}
