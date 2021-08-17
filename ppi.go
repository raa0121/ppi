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

func CreateImage(img *psd.PSD, conf *pfv.Pfv) []Image {
	output := map[string][]psd.Layer{}
	imgs := []Image{}
	canvas := &image.RGBA{}

	for _, v := range conf.Items {
		for _, vv := range v.Elements {
			for _, l := range img.Layer {
				if l.Folder() {
					folder := detectLayerNameEncoding([]byte(l.Name))
					for _, ll := range l.Layer {
						name := detectLayerNameEncoding([]byte(ll.Name))
						path := string(folder) + "/" + string(name) 
						if vv.Path == path {
							output[v.Name] = append(output[v.Name], ll)
						}
					}
				} else {
					name := detectLayerNameEncoding([]byte(l.Name))
					path := string(name) 
					if vv.Path == path {
						output[v.Name] = append(output[v.Name], l)
					}
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
