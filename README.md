# Tubular Hackathon


GOAL:
An app that takes an Instagram profile as an input and generates a mosaic of their most famous photo's commenters.

### Web application

- [ ] Write a basic web application that takes and processes an input (Instagram username)


BONUS/IDEAS:
- Allow choice of top 5 most famous photos as a 2nd step
- Authentication
- Autocomplete Instagram profile name?

### Mosaic builder


- [x] Write a matrix of image.Image regions that represent the original
   - `func genMatrixRegions(regionLength int) []Region`
- [x] Write a function that returns base color for each region (findBaseColour)  
   - `func findBaseColor(image image.Image) color.RGBA` 
- [x] Write a function that calculates color distance between two colors
   - `func getDistanceBetweenColors(color1 color.Color, color2 color.Color) float`
- [x] Write a function that opens up image set (openImageSet)  
- [x] Write a function that goes through each region and assigns an image (assignImage)  
- [x] Write a function that generates the image (writeMosaic)  
