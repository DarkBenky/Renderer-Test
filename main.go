package main

import (
	"image/color"
	"log"
	"math"
	"sort"

	"golang.org/x/image/colornames"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

type Triangle struct {
	vertices [3]Vector2
	color    color.Color
	z        float64 // Z value for depth sorting
}

type Vector2 struct {
	x, y float64
}

type Game struct {
	triangles []Triangle
	cameraX   float64
	cameraY   float64
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.cameraY -= 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.cameraY += 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.cameraX -= 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.cameraX += 0.1
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen
	screen.Fill(colornames.Black)

	// Sort triangles by depth (z value)
	sortedTriangles := make([]Triangle, len(g.triangles))
	copy(sortedTriangles, g.triangles)
	sortTriangles(sortedTriangles, g.cameraX, g.cameraY)

	// Draw sorted triangles
	for _, tri := range sortedTriangles {
		drawTriangle(screen, tri, g.cameraX, g.cameraY)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func sortTriangles(triangles []Triangle, cameraX, cameraY float64) {
	sort.Slice(triangles, func(i, j int) bool {
		distanceI := math.Sqrt(math.Pow(triangles[i].vertices[0].x-cameraX, 2) + math.Pow(triangles[i].vertices[0].y-cameraY, 2))
		distanceJ := math.Sqrt(math.Pow(triangles[j].vertices[0].x-cameraX, 2) + math.Pow(triangles[j].vertices[0].y-cameraY, 2))
		return distanceI > distanceJ // Sort by distance from camera (far to near)
	})
}

func drawTriangle(screen *ebiten.Image, tri Triangle, cameraX, cameraY float64) {
	// Calculate screen positions based on triangle vertices and camera position
	v0 := screenPosition(tri.vertices[0], cameraX, cameraY)
	v1 := screenPosition(tri.vertices[1], cameraX, cameraY)
	v2 := screenPosition(tri.vertices[2], cameraX, cameraY)

	// Create a new image to hold the triangle
	triangleImage := ebiten.NewImage(100, 100) // size is arbitrary
	triangleImage.Fill(color.Transparent)

	// Draw the triangle on the new image
	ebitenutil.DrawLine(triangleImage, v0.x, v0.y, v1.x, v1.y, tri.color)
	ebitenutil.DrawLine(triangleImage, v1.x, v1.y, v2.x, v2.y, tri.color)
	ebitenutil.DrawLine(triangleImage, v2.x, v2.y, v0.x, v0.y, tri.color)

	// Draw the triangle image on the screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(v0.x, v0.y)
	screen.DrawImage(triangleImage, op)
}

func screenPosition(v Vector2, cameraX, cameraY float64) Vector2 {
	// Scale and translate the vertices based on the camera position
	scale := 100.0 // Scale factor
	return Vector2{
		x: (v.x - cameraX) * scale + screenWidth/2,
		y: (v.y - cameraY) * scale + screenHeight/2,
	}
}

func main() {
	// Create some sample triangles
	triangles := []Triangle{
		{
			vertices: [3]Vector2{
				{x: 0, y: 0},
				{x: 1, y: 1},
				{x: -1, y: 1},
			},
			color: color.RGBA{255, 0, 0, 255},
			z:     2,
		},
		{
			vertices: [3]Vector2{
				{x: 1, y: 0},
				{x: 2, y: 1},
				{x: 0.5, y: 1},
			},
			color: color.RGBA{0, 255, 0, 255},
			z:     1,
		},
		{
			vertices: [3]Vector2{
				{x: -1, y: -1},
				{x: 0, y: 0},
				{x: -1.5, y: 0.5},
			},
			color: color.RGBA{0, 0, 255, 255},
			z:     3,
		},
	}

	game := &Game{
		triangles: triangles,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ebiten Triangle Rendering Example")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
