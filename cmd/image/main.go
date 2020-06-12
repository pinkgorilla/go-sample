package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// color pallete
var (
	green  = color.RGBA{R: 0x18, G: 0xBA, B: 0x62, A: 255}
	grey33 = color.RGBA{R: 0x33, G: 0x33, B: 0x33, A: 255}
	grey4B = color.RGBA{R: 0x4B, G: 0x4B, B: 0x4B, A: 255}
	blue   = color.RGBA{R: 0x3B, G: 0x8E, B: 0xF3, A: 255}
	red    = color.RGBA{R: 0xDB, G: 0x2B, B: 0x39, A: 255}
	blueA  = color.RGBA{R: 0x3B, G: 0x8E, B: 0xF3, A: 128}
)

func main() {
	canvas, err := GetCanvas()
	if err != nil {
		panic(err)
	}
	drawSummary(canvas, &Summary{
		Title:       "TOTAL UTANG",
		Amount:      -574000,
		Name:        "Gerry Pahabol",
		Note:        "Utang Anda per 11 JUNI 2020, 13:00",
		StoreName:   "TOKO SUMBER REJEKI",
		PhoneNumber: "081222212336",
	})
	// uncomment this lines to drawSingle image
	// drawSingle(canvas, &Single{
	// 	Title:       "Gerry Pahabol memberi",
	// 	Amount:      574000,
	// 	Date:        "Pada 11 JUNI 2020",
	// 	Note:        "Pembelian Ginjal Segar",
	// 	Name:        "TOKO SUMBER REJEKI",
	// 	PhoneNumber: "081222212336",
	// })

	// writer can be anything implements io.Writer, such as http.ResponseWriter
	// here we use file as writer, because we want to write the PNG to file
	writer, err := os.Create("single_with_png_template.png")
	if err != nil {
		log.Fatalf("failed to create: %s", err)
	}
	png.Encode(writer, canvas)
	// close writer accordingly, if you use http.ResponseWriter as writer
	// you might consider to omit the writer Close method because it will be called by Go
	defer writer.Close()

	for i := 10; i <= 100; i += 10 {
		writer, err := os.Create(fmt.Sprintf("single_with_png_template_%v.jpg", i))
		if err != nil {
			log.Fatalf("failed to create: %s", err)
		}
		jpeg.Encode(writer, canvas, &jpeg.Options{Quality: i})
		writer.Close()
	}
}

func getBase64Template() string {
	template, err := os.Open("template.jpg")
	if err != nil {
		fmt.Println(err)
	}
	defer template.Close()
	templateBytes, err := ioutil.ReadAll(template)
	base64Data := base64.StdEncoding.EncodeToString([]byte(templateBytes))
	return base64Data
}

func GetCanvas() (draw.Image, error) {

	bs, err := base64.StdEncoding.DecodeString(base64TemplatePNG)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(bs)
	template, err := png.Decode(buffer)
	if err != nil {
		return nil, err
	}
	canvas := image.NewRGBA(template.Bounds())
	draw.Draw(canvas, canvas.Bounds(), template, image.ZP, draw.Over)
	return canvas, nil
}

type Single struct {
	Title       string
	Amount      float64
	Date        string
	Name        string
	PhoneNumber string
	Note        string
}
type Summary struct {
	Title       string
	Amount      float64
	Name        string
	StoreName   string
	PhoneNumber string
	Note        string
}

// negative data.Amount value will cause the amount to be drawn with color red.
// otherwise data.Amount will be drawn with green
func drawSummary(canvas draw.Image, data *Summary) {
	ctx := freetype.NewContext()
	ctx.SetDst(canvas)
	ctx.SetClip(canvas.Bounds())
	ctx.SetDPI(72)

	drawText(ctx, image.Pt(60, 60), data.Title, 42, fazzneu.Bold(), grey33)
	p := message.NewPrinter(language.Indonesian)

	amountColor := green
	if data.Amount < 0 {
		amountColor = red
	}
	drawText(ctx, image.Pt(60, 120), p.Sprintf("Rp %.f", math.Abs(data.Amount)), 60, fazzneu.Bold(), amountColor)

	// draws the transparent blue area
	draw.DrawMask(canvas, image.Rect(40, 220, 960, 377), image.NewUniform(blue), image.Pt(40, 220), image.NewUniform(color.Alpha{25}), image.ZP, draw.Over)
	// draws text inside transparant blue area
	p1, _ := drawText(ctx, image.Pt(80, 250), data.Note, 36, fazzneu.Regular(), blue)
	p2, _ := drawText(ctx, image.Pt(80, p1.Y.Ceil()+18), "atas nama", 36, fazzneu.Regular(), blue)
	drawText(ctx, image.Pt(p2.X.Ceil()+10, p1.Y.Ceil()+18), data.Name, 36, fazzneu.Bold(), blue)

	drawText(ctx, image.Pt(60, 568), data.StoreName, 42, fazzneu.Bold(), blue)
	drawText(ctx, image.Pt(60, 632), data.PhoneNumber, 36, fazzneu.Regular(), grey4B)
}

// negative data.Amount value will cause the amount to be drawn with color red.
// otherwise data.Amount will be drawn with green
func drawSingle(canvas draw.Image, data *Single) {
	ctx := freetype.NewContext()
	ctx.SetDst(canvas)
	ctx.SetClip(canvas.Bounds())
	ctx.SetDPI(72)

	drawText(ctx, image.Pt(60, 60), data.Title, 42, fazzneu.Regular(), grey33)

	amountColor := green
	if data.Amount < 0 {
		amountColor = red
	}
	// drawText(ctx, image.Pt(60, 120), fmt.Sprintf("Rp %v", data.Amount), 60, fazzneu.Bold(), green)
	p := message.NewPrinter(language.Indonesian)
	drawText(ctx, image.Pt(60, 120), p.Sprintf("Rp %.f", math.Abs(data.Amount)), 60, fazzneu.Bold(), amountColor)
	drawText(ctx, image.Pt(60, 200), data.Date, 36, fazzneu.Regular(), grey4B)
	drawText(ctx, image.Pt(60, 315), "Catatan", 42, fazzneu.Bold(), grey33)
	// TODO: support note when note require multiple line
	drawText(ctx, image.Pt(60, 375), data.Note, 36, fazzneu.Regular(), grey4B)
	drawText(ctx, image.Pt(60, 568), data.Name, 42, fazzneu.Bold(), blue)
	drawText(ctx, image.Pt(60, 632), data.PhoneNumber, 36, fazzneu.Regular(), grey4B)
}

func drawText(ctx *freetype.Context, point image.Point, text string, fontSize float64, font *truetype.Font, color color.Color) (fixed.Point26_6, error) {
	ctx.SetSrc(image.NewUniform(color))
	ctx.SetFont(font)
	ctx.SetFontSize(fontSize)
	pt := freetype.Pt(point.X, point.Y+int(ctx.PointToFixed(fontSize)>>6))
	return ctx.DrawString(text, pt)
}

var fazzneu Fonts

// Fonts helps font operation, switching between bold and regular font
type Fonts struct {
	regular *truetype.Font
	bold    *truetype.Font
}

// Bold returns bold font
func (f *Fonts) Bold() *truetype.Font {
	if f.bold == nil {
		f.bold = f.decodeFont(base64FazzNeuBold)
	}
	return f.bold
}

// Regular returns regular font
func (f *Fonts) Regular() *truetype.Font {
	if f.regular == nil {
		f.regular = f.decodeFont(base64FazzNeuRegular)
	}
	return f.regular
}
func (f *Fonts) decodeFont(base64Font string) *truetype.Font {
	bs, err := base64.StdEncoding.DecodeString(base64Font)
	if err != nil {
		panic(err)
	}
	font, err := freetype.ParseFont(bs)
	if err != nil {
		panic(err)
	}
	return font
}
