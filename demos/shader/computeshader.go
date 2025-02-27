package shader

import (
	"fmt"
	"os"
	"time"

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

// This struct reflects an SSBO in the compute shader.
type DataBuffer3 struct {
	loc   math32.Vector3
	i     uint32
	speed float32
	_     [2]int32 // Using such offset, we can account for OpenGL padding if there is any
	data  [5]math64.Vector4
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

	var vectors []math32.Vector3 = make([]math32.Vector3, 30, 30)
	callback := func(b_ *gls.BufferRaw, deltaTime time.Duration) {
		// This callback is called everytime we DispatchCompute using Coman
		// (assuming the program using this buffer is dispatched)
		// The loop below demonstrates reading out an SSBO after the compute
		// shader has worked on it.

		//b := b_.Typed()
		//for i, v := range gls.BufferAsT[math32.Vector3](b) {
		// Commit changes to the buffer:
		//v.AddScalar(1)
		//gls.BufferSet[math32.Vector3](b, i, v)

		//// Also consider BufferGet when not using BufferAsT
		////gls.BufferGet[math32.Vector3](b, i)

		//fmt.Printf("%dth Vector: %+v\n", i, v)
		//}
	}
	bufferObjects := gls.NewBufferObjects()

	{
		vectorsBuffer := gls.SliceAsBuffer[math32.Vector3](vectors)
		ssbo := gls.NewSSBO(gs, 0,
			gls.BO_DYNAMIC_COPY, gls.BO_READ_WRITE, callback, vectorsBuffer.Size).SetInitialBuffer(&vectorsBuffer.BufferRaw)
		bufferObjects.Set(ssbo)
	}

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
	// Make sure to set the desired program
	a.Renderer().Coman().SetProgram(t.computeSpecs)
	// Call it once, call it twice - the number of times to repeat the compute
	// shader between frame renderings is up to the user. - It could also be
	// done outside of the update loop!
	a.Renderer().Coman().Compute(*t.workGroups, deltaTime)
}

// Cleanup is called once at the end of the demo.
func (t *ComputeDemo) Cleanup(a *app.App) {}
