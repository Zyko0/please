package assets

import (
	"bytes"
	_ "embed"
	"image/color"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	WhiteImage *ebiten.Image

	//go:embed images/gopher_enemy.png
	gopherEnemySrc   []byte
	GopherEnemyImage *ebiten.Image

	//go:embed images/gopher.png
	gopherSrc   []byte
	GopherImage *ebiten.Image

	//go:embed images/ebiten.png
	ebitenSrc   []byte
	EbitenImage *ebiten.Image

	//go:embed images/ebiten_block0.png
	ebitenBlock0Src   []byte
	EbitenBlock0Image *ebiten.Image

	//go:embed images/ebiten_block1.png
	ebitenBlock1Src   []byte
	EbitenBlock1Image *ebiten.Image

	//go:embed images/ebiten_block2.png
	ebitenBlock2Src   []byte
	EbitenBlock2Image *ebiten.Image

	//go:embed images/ebiten_block3.png
	ebitenBlock3Src   []byte
	EbitenBlock3Image *ebiten.Image

	//go:embed images/ebiten_block4.png
	ebitenBlock4Src   []byte
	EbitenBlock4Image *ebiten.Image

	//go:embed images/ebiten_res0.png
	ebitenRes0Src   []byte
	EbitenRes0Image *ebiten.Image

	//go:embed images/ebiten_res1.png
	ebitenRes1Src   []byte
	EbitenRes1Image *ebiten.Image

	//go:embed images/ebiten_res2.png
	ebitenRes2Src   []byte
	EbitenRes2Image *ebiten.Image

	//go:embed images/ebiten_res3.png
	ebitenRes3Src   []byte
	EbitenRes3Image *ebiten.Image

	//go:embed images/hajimehoshi.png
	hajimehoshiSrc   []byte
	HajimeHoshiImage *ebiten.Image

	//go:embed images/solarlune.png
	solarLuneSrc   []byte
	SolarLuneImage *ebiten.Image

	//go:embed images/dogcow.png
	dogcowSrc   []byte
	DogcowImage *ebiten.Image
)

func init() {
	var err error

	WhiteImage = ebiten.NewImage(3, 3)
	WhiteImage.Fill(color.White)

	img, err := png.Decode(bytes.NewReader(gopherEnemySrc))
	if err != nil {
		log.Fatal(err)
	}
	GopherEnemyImage = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(gopherSrc))
	if err != nil {
		log.Fatal(err)
	}
	GopherImage = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(ebitenSrc))
	if err != nil {
		log.Fatal(err)
	}
	EbitenImage = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(ebitenBlock0Src))
	if err != nil {
		log.Fatal(err)
	}
	EbitenBlock0Image = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(ebitenBlock1Src))
	if err != nil {
		log.Fatal(err)
	}
	EbitenBlock1Image = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(ebitenBlock2Src))
	if err != nil {
		log.Fatal(err)
	}
	EbitenBlock2Image = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(ebitenBlock3Src))
	if err != nil {
		log.Fatal(err)
	}
	EbitenBlock3Image = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(ebitenBlock4Src))
	if err != nil {
		log.Fatal(err)
	}
	EbitenBlock4Image = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(ebitenRes0Src))
	if err != nil {
		log.Fatal(err)
	}
	EbitenRes0Image = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(ebitenRes1Src))
	if err != nil {
		log.Fatal(err)
	}
	EbitenRes1Image = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(ebitenRes2Src))
	if err != nil {
		log.Fatal(err)
	}
	EbitenRes2Image = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(ebitenRes3Src))
	if err != nil {
		log.Fatal(err)
	}
	EbitenRes3Image = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(hajimehoshiSrc))
	if err != nil {
		log.Fatal(err)
	}
	HajimeHoshiImage = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(solarLuneSrc))
	if err != nil {
		log.Fatal(err)
	}
	SolarLuneImage = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(dogcowSrc))
	if err != nil {
		log.Fatal(err)
	}
	DogcowImage = ebiten.NewImageFromImage(img)
}
