package shader

import (
	"os"
	"time"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/g3nd/app"
)

func init() {
	app.DemoMap["shader.shaderfiles"] = &ShaderFiles{filepaths: &Filepaths{"shaders/vertex.glsl", "shaders/geometry.glsl", "shaders/fragment.glsl"}}
}

type Filepaths struct {
	vertex   string
	geometry string
	fragment string
}

func (f *Filepaths) loadFile(filepath string) string {
	b, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (f *Filepaths) LoadGeometry(path string) string {
	return f.loadFile(path + "/" + f.geometry)
}
func (f *Filepaths) LoadFragment(path string) string {
	return f.loadFile(path + "/" + f.fragment)
}
func (f *Filepaths) LoadVertex(path string) string {
	return f.loadFile(path + "/" + f.vertex)
}

type ShaderFiles struct {
	a         *app.App
	filepaths *Filepaths
	plane1    *graphic.Mesh
	box1      *graphic.Mesh
	sphere1   *graphic.Mesh
}

// Start is called once at the start of the demo.
func (t *ShaderFiles) Start(a *app.App) {

	// Adds directional front light
	dir1 := light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, 0.6)
	dir1.SetPosition(0, 0, 100)
	a.Scene().Add(dir1)

	// Add axes helper
	axes := helper.NewAxes(1)
	a.Scene().Add(axes)

	// Create custom shader
	a.Renderer().AddShader("shaderFilesVertex", t.filepaths.LoadVertex(a.DirData()))
	a.Renderer().AddShader("shaderFilesGeometry", t.filepaths.LoadGeometry(a.DirData()))
	a.Renderer().AddShader("shaderFilesFrag", t.filepaths.LoadFragment(a.DirData()))
	a.Renderer().AddProgram("shaderFiles", "shaderFilesVertex", "shaderFilesFrag", "shaderFilesGeometry")

	// Creates plane 1
	geom1 := geometry.NewPlane(2, 2)
	mat1 := NewMaterial(&math32.Color{R: 0.8, G: 0.2, B: 0.1})
	mat1.SetSide(material.SideDouble)
	mat1.SetShininess(10)
	mat1.SetSpecularColor(&math32.Color{R: 0, G: 0, B: 0})
	t.plane1 = graphic.NewMesh(geom1, mat1)
	t.plane1.SetPosition(-1.2, 1, 0)
	a.Scene().Add(t.plane1)

	// Creates box1
	geom2 := geometry.NewBox(2, 2, 1)
	mat2 := NewMaterial(&math32.Color{0.2, 0.4, 0.8})
	t.box1 = graphic.NewMesh(geom2, mat2)
	t.box1.SetPosition(1.2, 1, 0)
	a.Scene().Add(t.box1)

	// Creates sphere 1
	geom3 := geometry.NewSphere(1, 32, 16)
	mat3 := NewMaterial(&math32.Color{0.5, 0.6, 0.7})
	t.sphere1 = graphic.NewMesh(geom3, mat3)
	t.sphere1.SetPosition(0, -1.2, 0)
	a.Scene().Add(t.sphere1)
}

// Update is called every frame.
func (t *ShaderFiles) Update(a *app.App, deltaTime time.Duration) {
	t.plane1.RotateY(-0.005)
	t.box1.RotateY(-0.005)
	t.sphere1.RotateY(-0.005)
}

// Cleanup is called once at the end of the demo.
func (t *ShaderFiles) Cleanup(a *app.App) {}

//
// Custom material
//

type Material struct {
	material.Standard // Embedded standard material
	color             math32.Color
	size              math32.Vector2
	uniColor          gls.Uniform
	uniSize           gls.Uniform
}

func NewMaterial(color *math32.Color) *Material {
	m := new(Material)
	m.Standard.Init("shaderFiles", color)

	// Creates uniforms
	m.uniColor.Init("Color")
	m.uniSize.Init("Size")

	// Set initial values
	m.color = *color
	m.size.Set(0.5, 0.2)
	return m
}

func (m *Material) RenderSetup(gl *gls.GLS) {
	m.Standard.RenderSetup(gl)
	gl.Uniform3fv(m.uniColor.Location(gl), 1, &m.color.R)
	gl.Uniform2fv(m.uniSize.Location(gl), 1, &m.size.X)
}
