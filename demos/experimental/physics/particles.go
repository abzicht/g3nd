package physics

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"

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
		numParticles:      1e4,
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

// Start is called once at the start of the demo.
func (t *ParticleDemo) Start(a *app.App) {
	gs := a.Renderer().Coman().GLS()
	// Adds directional front light
	dir1 := light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, 1)
	dir1.SetPosition(0, 30, 100)
	a.Scene().Add(dir1)

	a.Renderer().Coman().AddShader("particle_compute_demo", loadFile(a.DirData()+"/"+t.computeShaderFile))
	a.Renderer().Coman().AddProgram("ParticleDemoProg", "particle_compute_demo")
	ssbos := gls.NewBufferObjects()
	{
		localSize := uint32(512) // NumWorkGroups: (512x1x1)
		t.workGroups = new(gls.NumWorkGroups)
		t.workGroups.X = uint32(t.numParticles+localSize-1) / localSize
		t.workGroups.Y = 1
		t.workGroups.Z = 1
	}
	var positions []math32.Vector3 // initial particle positions
	for i := uint32(0); i < t.numParticles; i++ {
		randVec := math32.NewVector3(rand.Float32(), rand.Float32(), rand.Float32())
		randVec.MultiplyScalar(0.05).AddScalar(0.5)
		var position math32.Vector3
		position.Add(t.boundsMax)
		position.Sub(t.boundsMin)
		position.Multiply(randVec)
		position.Add(t.boundsMin)
		positions = append(positions, position)
	}
	positionsBuffer := gls.SliceAsBuffer[math32.Vector3](positions)
	positionsBO := gls.NewSSBO(gs, "ParticlePos",
		gls.BO_DYNAMIC_COPY, gls.BO_READ_WRITE, nil, positionsBuffer.Size).SetInitialBuffer(&positionsBuffer.BufferRaw)
	ssbos.Set(positionsBO)
	colorBO := gls.NewSSBO(gs, "ParticleColor", gls.BO_DYNAMIC_COPY,
		gls.BO_READ_WRITE, nil, t.numParticles*uint32(gls.StrideofT[math32.Vector4]()))
	ssbos.Set(colorBO)
	{
		var velocities []math32.Vector3
		for i := uint32(0); i < t.numParticles; i++ {
			velocity := math32.NewVector3(rand.Float32(), rand.Float32(), rand.Float32())
			velocity.SubScalar(0.5)
			velocity.MultiplyScalar(0.01)
			velocities = append(velocities, *velocity)
		}
		velocitiesBuffer := gls.SliceAsBuffer[math32.Vector3](velocities)
		ssbos.Set(gls.NewSSBO(gs, "ParticleVel",
			gls.BO_DYNAMIC_COPY, gls.BO_READ_WRITE, nil, velocitiesBuffer.Size).SetInitialBuffer(&velocitiesBuffer.BufferRaw))
	}
	{
		t.computeSpecs = gls.NewComputeSpecs("ParticleDemoProg", "4_3", *gls.NewShaderDefines(), ssbos)
		_, err := a.Renderer().Coman().SetProgram(t.computeSpecs)
		if err != nil {
			fmt.Printf("Failed to set the shader program: %s\n", err)
			return
		}
	}

	t.particleMat = material.NewParticleMaterial(math32.Color4{R: 0.2, G: 0.4, B: 0.6, A: 0.8}, colorBO) // replace colorBO with nil to use a standard material with the given color
	t.particleMat.SetParticleSize(3)                                                                     // Set to -1 in conjunction with setting no shape to get true pixel particles instead of quads or geometries

	// Create geometry that takes its data from the created buffer
	t.particleGeometry = geometry.NewParticles(t.numParticles, positionsBO, t.boundsMax.Sub(t.boundsMin))
	t.particleGeometry.SetParticleShape(geometry.NewSphere(0.01, 8, 8)) // Leave out to deal with pixels only

	t.particleGraphics = graphic.NewParticleSim(t.particleGeometry, t.particleMat)
	t.particleGraphics.SetPosition(t.boundsMin.X, t.boundsMin.Y, t.boundsMin.Z)
	a.Scene().Add(t.particleGraphics)

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
	// Switch to our custom program
	a.Renderer().Coman().SetProgram(t.computeSpecs)
	// Dispatch the compute shader
	a.Renderer().Coman().Compute(*t.workGroups, deltaTime)
}

// Cleanup is called once at the end of the demo.
func (t *ParticleDemo) Cleanup(a *app.App) {
	// Be careful with deleting all programs, we can also just keep them
	//a.Renderer().Coman().DeletePrograms()
}
