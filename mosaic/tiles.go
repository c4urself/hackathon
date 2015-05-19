package mosaic

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Tile struct {
	img       image.Image
	offset    image.Point
	baseColor color.NRGBA
}

// split image on multiple tiles
func BreakToTiles(img image.Image, tileSize int) []Tile {
	imgSize := img.Bounds().Size()
	xPointsCount := int(math.Ceil(float64(imgSize.X) / float64(tileSize)))
	yPointsCount := int(math.Ceil(float64(imgSize.Y) / float64(tileSize)))
	rgbImage := img.(*image.RGBA)
	startPoint := image.Point{X: 0, Y: 0}

	log.Printf("Creating tiles for %vx%v\n", xPointsCount, yPointsCount)

	var tiles []Tile

	for i := 0; i < yPointsCount; i++ {
		for j := 0; j < xPointsCount; j++ {
			// get the x,y start and x,y end for the region and make a rectangle
			tileMinOffset := image.Point{X: j * tileSize, Y: i * tileSize}
			tileMaxOffset := image.Point{X: (j + 1) * tileSize, Y: (i + 1) * tileSize}
			originalRect := image.Rectangle{Min: tileMinOffset, Max: tileMaxOffset}

			// create a rectangle based off of the old one with start at 0,0
			tileRectangle := image.Rectangle{startPoint, startPoint.Add(originalRect.Size())}
			tileImg := image.NewRGBA(tileRectangle)
			draw.Draw(tileImg, tileRectangle, rgbImage, tileMinOffset, draw.Src)

			// add to array of Tiles
			tile := Tile{img: tileImg, offset: tileMinOffset, baseColor: GetColorProfile(tileImg)}
			tiles = append(tiles, tile)
		}
	}

	log.Printf("Created %v tiles", len(tiles))

	return tiles
}

// Generate final image from the set of tiles
func CollectFromTiles(tiles []Tile) image.Image {
	// Get size of the final image
	width := 0
	height := 0

	for _, tile := range tiles {
		xOffset := tile.offset.X + tile.img.Bounds().Size().X
		yOffset := tile.offset.Y + tile.img.Bounds().Size().Y

		if width < xOffset {
			width = xOffset
		}

		if height < yOffset {
			height = yOffset
		}
	}

	log.Printf("Drawing mosaic")
	// Draw the image
	destinationImage := image.NewRGBA(image.Rect(0, 0, width, height))
	for _, tile := range tiles {
		draw.Draw(destinationImage,
			image.Rectangle{Min: tile.offset, Max: tile.offset.Add(tile.img.Bounds().Size())},
			tile.img,
			image.Point{X: 0, Y: 0},
			draw.Src)
	}

	return destinationImage
}

// Imports tiles from thumbnails directory
func ImportTiles(baseDir string, tileSize int) []Tile {
	var tiles []Tile
	var totalFound int
	filepath.Walk(baseDir, func(path string, _ os.FileInfo, _ error) error {

		// only accept PNG for now
		if strings.HasSuffix(path, ".png") {
			totalFound += 1
			fImg, _ := os.Open(path)
			defer fImg.Close()
			img, _, err := image.Decode(fImg)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			img = resize.Resize(uint(tileSize), 0, img, resize.NearestNeighbor)
			tiles = append(tiles, Tile{img: img, offset: image.Point{X: 0, Y: 0}, baseColor: GetColorProfile(img)})
		}
		return nil
	})

	if totalFound == 0 {
		log.Panicf("No thumbnails found!")
	}
	log.Printf("Found %v thumbnails", totalFound)

	return tiles
}

// Matches original tiles with list of suggested
func FindSimilarTiles(originals []Tile, candidates []Tile) []Tile {

	var similar []Tile

	log.Printf("Finding matching tiles")

	for _, original := range originals {
		var pick TileDistance
		var chosenTile Tile
		var distanceMap = make(map[Tile]float64)

		for _, candidate := range candidates {
			distance := GetColorDistance(original.baseColor, candidate.baseColor)
			distanceMap[candidate] = distance
		}
		orderedCandidates := sortMapByValue(distanceMap)
		position := rand.Intn(5)
		pick = orderedCandidates[position]
		chosenTile = Tile{img: pick.tile.img, offset: original.offset, baseColor: pick.tile.baseColor}
		similar = append(similar, chosenTile)
	}

	return similar
}

type TileDistance struct {
	tile     Tile
	distance float64
}

type OrderedTileDistanceList []TileDistance

func (p OrderedTileDistanceList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p OrderedTileDistanceList) Len() int           { return len(p) }
func (p OrderedTileDistanceList) Less(i, j int) bool { return p[i].distance < p[j].distance }

func sortMapByValue(m map[Tile]float64) OrderedTileDistanceList {
	var p OrderedTileDistanceList = make(OrderedTileDistanceList, len(m))
	i := 0
	for k, v := range m {
		p[i] = TileDistance{k, v}
		i++
	}
	sort.Sort(p)
	return p[0:10]
}
