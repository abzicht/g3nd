package physics

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"
	"unsafe"

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
	app.DemoMap["physics-experimental.particles"] = &ParticleDemo{
		computeShaderFile: "shaders/particle_compute.glsl",
	}
}

func loadFile(filepath string) string {
	b, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type ParticleDemo struct {
	a                       *app.App
	workGroups              *gls.NumWorkGroups
	particleGraphics        *graphic.ParticleSim
	particleGeometry        *geometry.ParticleGeometry
	mat                     *material.ParticleMaterial
	computeSpecs            *gls.ComputeSpecs
	computeShaderFile       string
	gridWidthU, gridHeightU gls.Uniform
	numParticles            uint32
	gridWidth, gridHeight   uint32
}

type Particle struct {
	pos      math32.Vector3
	velocity math32.Vector2 //rotation
}

// Start is called once at the start of the demo.
func (t *ParticleDemo) Start(a *app.App) {
	// Adds directional front light
	dir1 := light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, 0.6)
	dir1.SetPosition(0, 0, 100)
	a.Scene().Add(dir1)

	// Add axes helper
	axes := helper.NewAxes(1)
	a.Scene().Add(axes)

	a.Renderer().Coman().AddShader("particle_compute_demo", loadFile(a.DirData()+"/"+t.computeShaderFile))
	a.Renderer().Coman().AddProgram("ParticleDemoProg", "particle_compute_demo")
	callback := func(b_ *gls.BufferRaw, deltaTime time.Duration) {
		//b := b_.Typed()
		//for i, p := range b.AsVec3() {
		//	fmt.Printf("part. %d (pos: %+v); ", i, p)
		//}
	}
	{
		t.gridWidth, t.gridHeight = 100, 10
		var gridDepth uint32 = 10
		t.numParticles = uint32(t.gridWidth * t.gridHeight * gridDepth)
		localSize := uint32(8) // Workgroup size (8x8x8)
		t.workGroups = new(gls.NumWorkGroups)
		t.workGroups.X = uint32(t.gridWidth+localSize-1) / localSize
		t.workGroups.Y = uint32(t.gridHeight+localSize-1) / localSize
		t.workGroups.Z = uint32(gridDepth+localSize-1) / localSize
	}
	positionsBO := gls.NewSSBO(a.Renderer().Coman().GLS(), 0,
		gls.BO_DYNAMIC_COPY, gls.BO_READ_WRITE, callback,
		uint32(gls.TypeSize(t.numParticles)*gls.SizeofT[math32.Vector3]()))
	particles := math32.NewArrayF32(int(3*t.numParticles), int(3*t.numParticles))
	for i := 0; i < len(particles); i += 3 {
		particles[i] = (rand.Float32() - 0.5) * 2
		particles[i+1] = (rand.Float32() - 0.5) * 2
		particles[i+2] = (rand.Float32() - 0.5) * 2
	}
	positionsBO.SetInitialBuffer(gls.NewBufferRaw(
		unsafe.Pointer(unsafe.SliceData(particles)), uint32(gls.TypeSize(t.numParticles)*gls.SizeofT[math32.Vector3]())))
	{
		ssbos := gls.NewBufferObjects()
		ssbos.Set(positionsBO)
		t.computeSpecs = gls.NewComputeSpecs("ParticleDemoProg", "4_3", *gls.NewShaderDefines(), ssbos)
		_, err := a.Renderer().Coman().SetProgram(t.computeSpecs)
		if err != nil {
			fmt.Printf("Failed to set the shader program: %s\n", err)
			return
		}
	}

	t.mat = material.NewParticleMaterial(math32.Color4{R: 0.2, G: 0.4, B: 0.6, A: 0.8})
	t.mat.SetParticleSize(10)
	// Create geometry that takes its data from the created buffer
	t.particleGeometry = geometry.NewParticles(t.numParticles, math32.NewVector3(5, 5, 5))
	t.particleGraphics = graphic.NewParticleSim(t.particleGeometry, t.mat)
	t.particleGraphics.SetPosition(0, 1, 0)
	a.Scene().Add(t.particleGraphics)
	if false {
		geom2 := geometry.NewSphere(.5, 32, 16)
		mat2 := material.NewStandard(&math32.Color{R: 1, G: 1, B: 1})
		mat2.SetWireframe(false)
		mat2.SetSide(material.SideDouble)
		sphere2 := graphic.NewMesh(geom2, mat2)
		sphere2.SetPosition(0.5, 0, 0)
		a.Scene().Add(sphere2)
	}

	t.gridWidthU.Init("gridWidth")
	t.gridHeightU.Init("gridHeight")
}

// Update is called every frame.
func (t *ParticleDemo) Update(a *app.App, deltaTime time.Duration) {
	// called just before the call to render
	a.Renderer().Coman().SetProgram(t.computeSpecs)
	{
		gl := a.Renderer().Coman().GLS()
		gl.Uniform1i(t.gridWidthU.Location(gl), int32(t.gridWidth))
		gl.Uniform1i(t.gridHeightU.Location(gl), int32(t.gridHeight))
	}
	a.Renderer().Coman().Compute(*t.workGroups, deltaTime)
}

// Cleanup is called once at the end of the demo.
func (t *ParticleDemo) Cleanup(a *app.App) {
	a.Renderer().Coman().DeletePrograms()
}
