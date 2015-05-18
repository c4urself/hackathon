# hackathon
Go webserver and image parser ?


- [ ] Write a matrix of image.Image regions that represent the original
   - `func genMatrixRegions(regionLength int) []Region`
- [ ] Write a function that returns base color for each region (findBaseColour)  
   - `func findBaseColor(image image.Image) color.RGBA` 
- [ ] Write a function that calculates color distance between two colors
   - `func getDistanceBetweenColors(color1 color.Color, color2 color.Color) float`
- [ ] Write a function that opens up image set (openImageSet)  
- [ ] Write a function that goes through each region and assigns an image (assignImage)  
- [ ] Write a function that generates the image (writeMosaic)  
