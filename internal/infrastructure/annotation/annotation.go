package annotation

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	// embed is used to store fonts
	_ "embed"

	"github.com/lonepeon/sport/internal/domain"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

//go:embed Montserrat-Regular.ttf
var montSerratRegularFile []byte
var montSerratRegularFont *opentype.Font

func init() {
	font, err := opentype.Parse(montSerratRegularFile)
	if err != nil {
		panic(fmt.Sprintf("can't parse MontSerrat-Regular font: %v", err))
	}

	montSerratRegularFont = font
}

type Annotation struct {
}

func (a Annotation) AnnotateMapWithStats(ctx context.Context, file domain.MapFile, distance domain.Distance, speed domain.Speed) (domain.ShareableMapFile, error) {
	src, err := png.Decode(file.File())
	if err != nil {
		return domain.ShareableMapFile{}, fmt.Errorf("can't decode image from png: %v", err)
	}

	overlayedImage := image.NewRGBA(src.Bounds())
	draw.Draw(overlayedImage, overlayedImage.Bounds(), src, image.Pt(0, 0), draw.Src)
	bg := &image.Uniform{color.RGBA{0xFF, 0xFF, 0xFF, 0x77}}
	draw.Draw(overlayedImage, image.Rect(0, src.Bounds().Max.Y-200, src.Bounds().Max.X, src.Bounds().Max.Y), bg, image.Pt(0, 0), draw.Over)

	face, err := opentype.NewFace(montSerratRegularFont, &opentype.FaceOptions{
		Size:    100,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	if err != nil {
		return domain.ShareableMapFile{}, fmt.Errorf("can't setup font: %v", err)
	}

	hpadding := 30
	drawing := font.Drawer{
		Dst:  overlayedImage,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(hpadding, src.Bounds().Max.Y-60),
	}

	drawing.DrawString(fmt.Sprintf("%.2fkm", distance.Kilometers()))
	speedLabel := fmt.Sprintf("%.2fkm/h", speed.KilometersPerHour())
	length := font.MeasureString(face, speedLabel)
	drawing.Dot = fixed.P(src.Bounds().Max.X-length.Round()-hpadding, drawing.Dot.Y.Round())
	drawing.DrawString(speedLabel)

	var buf bytes.Buffer
	if err := png.Encode(&buf, overlayedImage); err != nil {
		return domain.ShareableMapFile{}, fmt.Errorf("can't encode image to png: %v", err)
	}

	return domain.NewSharableMapFile(buf.Bytes()), nil
}
