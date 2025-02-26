package shader

import (
	"fmt"
	"os"
	"time"
	"unsafe"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/math64"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/g3nd/app"
)

func init() {
	app.DemoMap["shader.computeshader"] = &ComputeDemo{computefile: "shaders/demo-compute.glsl", vertexfile: "shaders/demo-vertex.glsl", fragmentfile: "shaders/demo-fragment.glsl", workGroups: gls.NewNumWorkGroups(8, 8, 8)}
}

func loadFile(filepath string) string {
	b, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type ComputeDemo struct {
	a            *app.App
	computefile  string
	vertexfile   string
	fragmentfile string
	plane1       *graphic.Mesh
	workGroups   *gls.NumWorkGroups
	computeSpecs *gls.ComputeSpecs
	shaderSpecs  renderer.ShaderSpecs
}

type DataBuffer3 struct {
	loc   math32.Vector3
	i     uint32
	speed float32
	_     [2]int32
	data  [5]math64.Vector4
}

func (d *DataBuffer3) ToString() string {
	//return fmt.Sprintf("i %d, loc %p, speed %f, data %s", d.i, d.loc, d.speed, d.data)
	return fmt.Sprintf("%+v", *d)
}

// Start is called once at the start of the demo.
func (t *ComputeDemo) Start(a *app.App) {
	gs := a.Renderer().Coman().GLS()
	// Adds directional front light
	dir1 := light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, 0.6)
	dir1.SetPosition(0, 0, 100)
	a.Scene().Add(dir1)

	// Add axes helper
	axes := helper.NewAxes(1)
	a.Scene().Add(axes)

	// Create custom shader
	a.Renderer().AddShader("vertex-shader", loadFile(a.DirData()+"/"+t.vertexfile))
	a.Renderer().AddShader("fragment-shader", loadFile(a.DirData()+"/"+t.fragmentfile))
	a.Renderer().Coman().AddShader("compute-shader", loadFile(a.DirData()+"/"+t.computefile))

	a.Renderer().Coman().AddProgram("ComputeProgram", "compute-shader")
	a.Renderer().AddProgram("FragProgram", "vertex-shader", "fragment-shader")

	var vectors [30]math32.Vector3
	vec3length := len(vectors)
	callback := func(b_ *gls.BufferRaw, deltaTime time.Duration) {
		b := b_.Typed()
		for i := 0; i < vec3length; i++ {
			v, err := gls.Get[math32.Vector3](b, 0)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%dth Vector: %+v\n", i, v)
		}
	}
	bufferObjects := gls.NewBufferObjects()

	ssbo := gls.NewSSBO(gs, 0,
		gls.BO_DYNAMIC_COPY, gls.BO_READ_WRITE, callback, uint32(unsafe.Sizeof(vectors)))
	bufferObjects.Set(ssbo)

	t.computeSpecs = gls.NewComputeSpecs("ComputeProgram", "4_3", *gls.NewShaderDefines(), bufferObjects)
	_, err := a.Renderer().Coman().SetProgram(t.computeSpecs)
	if err != nil {
		fmt.Printf("Failed to set the shader program: %s\n", err)
		return
	}

	t.shaderSpecs.Defines = *gls.NewShaderDefines()
	t.shaderSpecs.Name = "FragProgram"
	a.Renderer().SetProgram(&t.shaderSpecs)
	_, err = a.Renderer().Coman().SetProgram(t.computeSpecs)
	if err != nil {
		fmt.Printf("Failed to set the shader program: %s\n", err)
		return
	}

	geom1 := geometry.NewPlane(2, 2)
	mat1 := NewComputeMaterial(&math32.Color{R: 0.0, G: 0.8, B: 1.0})
	mat1.SetSide(material.SideDouble)
	mat1.SetShininess(10)
	mat1.SetSpecularColor(&math32.Color{R: 0, G: 0, B: 0})
	t.plane1 = graphic.NewMesh(geom1, mat1)
	t.plane1.SetPosition(0, 0, 0)
	a.Scene().Add(t.plane1)
	t.plane1.RotateY(-0.001)
}

type ComputeMaterial struct {
	material.Standard // Embedded standard material
	color             math32.Color
	uniColor          gls.Uniform
}

func NewComputeMaterial(color *math32.Color) *ComputeMaterial {

	m := new(ComputeMaterial)
	m.Standard.Init("FragProgram", color)

	m.uniColor.Init("Color")
	m.color = *color
	return m
}

func (m *ComputeMaterial) RenderSetup(gl *gls.GLS) {

	m.Standard.RenderSetup(gl)
	gl.Uniform3fv(m.uniColor.Location(gl), 1, &m.color.R)
}

// Update is called every frame.
func (t *ComputeDemo) Update(a *app.App, deltaTime time.Duration) {
	t.plane1.RotateY(-0.005)
	a.Renderer().Coman().SetProgram(t.computeSpecs)
	a.Renderer().Coman().Compute(*t.workGroups, deltaTime)
}

// Cleanup is called once at the end of the demo.
func (t *ComputeDemo) Cleanup(a *app.App) {}
