package imgen

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

func GenerateImg(pathToFile string, pathToImg string ) (err error) {
	binaryData, err := os.ReadFile(pathToFile)
	if err != nil {
		log.Printf("Ошибка при чтении файла: %v", err)
		return err
	}

	length := len(binaryData)
	side := int(math.Ceil(math.Sqrt(float64(length))))

	img := image.NewRGBA(image.Rect(0, 0, side, side))

	for i, b := range binaryData {
		x := i % side
		y := i / side
		img.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: b}) 
	}

	outFile, err := os.Create(pathToImg)
	if err != nil {
		log.Fatalf("Ошибка при создании файла: %v", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		log.Fatalf("Error encoding image: %v", err)
	}

	log.Println("Image successfully created")
	return 
}
