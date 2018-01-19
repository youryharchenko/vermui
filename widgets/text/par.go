// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package text

import "github.com/verdverm/vermui/lib/render"

// Par displays a paragraph.
/*
  par := vermui.NewPar("Simple Text")
  par.Height = 3
  par.Width = 17
  par.BorderLabel = "Label"
*/
type Par struct {
	render.Block
	Text        string
	TextFgColor render.Attribute
	TextBgColor render.Attribute
	WrapLength  int // words wrap limit. Note it may not work properly with multi-width char
}

// NewPar returns a new *Par with given text as its content.
func NewPar(s string) *Par {
	return &Par{
		Block:       *render.NewBlock(),
		Text:        s,
		TextFgColor: render.ThemeAttr("par.text.fg"),
		TextBgColor: render.ThemeAttr("par.text.bg"),
		WrapLength:  0,
	}
}

// Buffer implements Bufferer interface.
func (p *Par) Buffer() render.Buffer {
	buf := p.Block.Buffer()

	fg, bg := p.TextFgColor, p.TextBgColor
	cs := render.DefaultTxBuilder.Build(p.Text, fg, bg)

	// wrap if WrapLength set
	if p.WrapLength < 0 {
		cs = render.WrapTx(cs, p.Width-2)
	} else if p.WrapLength > 0 {
		cs = render.WrapTx(cs, p.WrapLength)
	}

	y, x, n := 0, 0, 0
	for y < p.InnerArea().Dy() && n < len(cs) {
		w := cs[n].Width()
		if cs[n].Ch == '\n' || x+w > p.InnerArea().Dx() {
			y++
			x = 0 // set x = 0
			if cs[n].Ch == '\n' {
				n++
			}

			if y >= p.InnerArea().Dy() {
				buf.Set(p.InnerArea().Min.X+p.InnerArea().Dx()-1,
					p.InnerArea().Min.Y+p.InnerArea().Dy()-1,
					render.Cell{Ch: '…', Fg: p.TextFgColor, Bg: p.TextBgColor})
				break
			}
			continue
		}

		buf.Set(p.InnerArea().Min.X+x, p.InnerArea().Min.Y+y, cs[n])

		n++
		x += w
	}

	return buf
}
