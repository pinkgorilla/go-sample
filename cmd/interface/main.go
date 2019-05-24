package main

import "log"

type AlatTulis interface {
	Write(text string)
}

type Pencil struct {
}

func (p Pencil) Write(text string) {
	log.Println("The hand writes elegantly using a Pen, and it writes: ", text)
}

type Pen struct {
	Color string
}

func (p Pen) Write(text string) {
	log.Println("The hand writes elegantly using a Pen, and it writes: ", text, "in ", p.Color, " color")
}

type Hand struct {
}

func (h Hand) Write(a AlatTulis) {
	a.Write("hello interface")
}

func main() {
	h := Hand{}
	pen := Pen{}
	h.Write(pen)
	pencil := Pencil{}
	h.Write(pencil)
}
