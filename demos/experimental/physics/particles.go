package physics

import (
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
	app.DemoMap["physics-experimental.particles"] = &ParticleDemo{
		workGroups: gls.NewNumWorkGroups(8, 8, 8),
	}
}

type ParticleDemo struct {
	a                *app.App
	workGroups       *gls.NumWorkGroups
	particleGraphics *graphic.ParticleSim
	particleGeometry *geometry.ParticleGeometry
	mat              *material.ParticleMaterial
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

	var numParticles uint32 = 3000
	//var particlesBufferSize uint32
	//{
	//	var particles_ [30]Particle
	//	numParticles = uint32(len(particles_))
	//	//particlesBufferSize = uint32(unsafe.Sizeof(particles_))
	//}

	// Creates bounding box for all particles
	t.mat = material.NewParticleMaterial(math32.Color4{R: 0.2, G: 0.4, B: 0.6, A: 1})
	t.mat.SetParticleSize(20)
	t.particleGeometry = geometry.NewParticles(numParticles, 2, 2, 2)
	t.particleGraphics = graphic.NewParticleSim(t.particleGeometry, t.mat)
	t.particleGraphics.SetPosition(0.2, 1, 0)
	a.Scene().Add(t.particleGraphics)
}

// Update is called every frame.
func (t *ParticleDemo) Update(a *app.App, deltaTime time.Duration) {
	// called just before the call to render
	//a.Renderer().Coman().SetProgram(t.particleGeometry.ComputeSpecs())
	//a.Renderer().Coman().Compute(*t.workGroups, deltaTime)
}

// Cleanup is called once at the end of the demo.
func (t *ParticleDemo) Cleanup(a *app.App) {
	//a.Renderer().Coman().DeleteProgram(&t.computeSpecs)
}
