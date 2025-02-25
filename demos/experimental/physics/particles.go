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
	"github.com/g3n/g3nd/app"
)

func init() {
	app.DemoMap["physics-experimental.particles"] = &ParticleDemo{
		computeShaderFile: "shaders/particle_compute.glsl",
		boundsMin:         math32.NewVector3(0, 0, 0),
		boundsMax:         math32.NewVector3(2, 2, 2),
		numParticles:      1e5,
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
	a                 *app.App
	workGroups        *gls.NumWorkGroups
	particleGraphics  *graphic.ParticleSim
	particleGeometry  *geometry.ParticleGeometry
	particleMat       *material.ParticleMaterial
	computeSpecs      *gls.ComputeSpecs
	computeShaderFile string
	numParticles      uint32
	boundsMin         *math32.Vector3
	boundsMax         *math32.Vector3
}

type Particle struct {
	pos      math32.Vector3
	velocity math32.Vector2 //rotation
}

// Start is called once at the start of the demo.
func (t *ParticleDemo) Start(a *app.App) {
	gs := a.Renderer().Coman().GLS()
	// Adds directional front light
	dir1 := light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, 1)
	dir1.SetPosition(0, 30, 100)
	a.Scene().Add(dir1)

	//// Add axes helper
	//axes := helper.NewAxes(1)
	//a.Scene().Add(axes)

	a.Renderer().Coman().AddShader("particle_compute_demo", loadFile(a.DirData()+"/"+t.computeShaderFile))
	//a.Renderer().Coman().AddChunk("", loadFile(a.DirData()+"/"+t.computeShaderFile))
	a.Renderer().Coman().AddProgram("ParticleDemoProg", "particle_compute_demo")
	ssbos := gls.NewBufferObjects()
	{
		localSize := uint32(512) // Workgroup size (8x8x8)
		t.workGroups = new(gls.NumWorkGroups)
		t.workGroups.X = uint32(t.numParticles+localSize-1) / localSize
		t.workGroups.Y = 1
		t.workGroups.Z = 1
	}
	var particles []math32.Vector3
	for i := uint32(0); i < t.numParticles; i++ {
		randVec := math32.NewVector3(rand.Float32(), rand.Float32(), rand.Float32())
		randVec.MultiplyScalar(0.05).AddScalar(0.5)
		var particle math32.Vector3
		particle.Add(t.boundsMax)
		particle.Sub(t.boundsMin)
		particle.Multiply(randVec)
		particle.Add(t.boundsMin)
		particles = append(particles, particle)
	}
	bufferSize := t.numParticles * uint32(gls.SizeofT[math32.Vector3]())
	positionsBO := gls.NewSSBO(gs, geometry.ParticlePositionBinding,
		gls.BO_DYNAMIC_COPY, gls.BO_READ_WRITE, nil, bufferSize).SetInitialBuffer(gls.NewBufferRaw(
		unsafe.Pointer(unsafe.SliceData(particles)), bufferSize))
	ssbos.Set(positionsBO)
	colorBO := gls.NewSSBO(gs, material.ParticleColorBinding, gls.BO_DYNAMIC_COPY,
		gls.BO_READ_WRITE, nil, t.numParticles*uint32(gls.SizeofT[math32.Vector4]()))
	ssbos.Set(colorBO)
	{
		var velocities []math32.Vector3
		for i := uint32(0); i < t.numParticles; i++ {
			velocity := math32.NewVector3(rand.Float32(), rand.Float32(), rand.Float32())
			velocity.SubScalar(0.5)
			velocity.MultiplyScalar(0.01)
			velocities = append(velocities, *velocity)
		}
		bufferSize = t.numParticles * uint32(gls.SizeofT[math32.Vector3]())
		ssbos.Set(gls.NewSSBO(gs, 3,
			gls.BO_DYNAMIC_COPY, gls.BO_READ_WRITE, nil, bufferSize).SetInitialBuffer(gls.NewBufferRaw(
			unsafe.Pointer(unsafe.SliceData(velocities)), bufferSize)))
	}
	{
		t.computeSpecs = gls.NewComputeSpecs("ParticleDemoProg", "4_3", *gls.NewShaderDefines(), ssbos)
		_, err := a.Renderer().Coman().SetProgram(t.computeSpecs)
		if err != nil {
			fmt.Printf("Failed to set the shader program: %s\n", err)
			return
		}
	}

	t.particleMat = material.NewParticleMaterial(math32.Color4{R: 0.2, G: 0.4, B: 0.6, A: 0.8}, nil)
	t.particleMat.SetParticleSize(30)

	// Create geometry that takes its data from the created buffer
	t.particleGeometry = geometry.NewParticles(t.numParticles, positionsBO, t.boundsMax.Sub(t.boundsMin))
	//t.particleGeometry.SetParticleShape(geometry.NewCone(0.04, 0.09, 3, 3, true))
	t.particleGeometry.SetParticleShape(geometry.NewSphere(0.01, 8, 8))

	t.particleGraphics = graphic.NewParticleSim(t.particleGeometry, t.particleMat)
	t.particleGraphics.SetPosition(t.boundsMin.X, t.boundsMin.Y, t.boundsMin.Z)
	a.Scene().Add(t.particleGraphics)
	if false {
		geom2 := geometry.NewSphere(.5, 32, 16)
		mat2 := material.NewStandard(&math32.Color{R: 1, G: 1, B: 1})
		mat2.SetWireframe(false)
		mat2.SetSide(material.SideDouble)
		sphere2 := graphic.NewMesh(geom2, mat2)
		sphere2.SetPosition(2.5, 1, 2)
		a.Scene().Add(sphere2)
	}

	var uniBoundsMin gls.Uniform
	uniBoundsMin.Init("BoundsMin")
	gs.Uniform3fv(uniBoundsMin.Location(gs), 1, &t.boundsMin.X)
	var uniBoundsMax gls.Uniform
	uniBoundsMax.Init("BoundsMax")
	gs.Uniform3fv(uniBoundsMax.Location(gs), 1, &t.boundsMax.X)
}

// Update is called every frame.
func (t *ParticleDemo) Update(a *app.App, deltaTime time.Duration) {
	// called just before the call to render
	a.Renderer().Coman().SetProgram(t.computeSpecs)
	a.Renderer().Coman().Compute(*t.workGroups, deltaTime)
}

// Cleanup is called once at the end of the demo.
func (t *ParticleDemo) Cleanup(a *app.App) {
	//a.Renderer().Coman().DeletePrograms()
}
