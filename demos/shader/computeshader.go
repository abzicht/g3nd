package shader

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/g3nd/app"
)

func init() {
	app.DemoMap["shader.computeshader"] = &ComputeDemo{filepath: "shaders/compute2.glsl", workGroups: gls.NewNumWorkGroups(16, 1, 1)}
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
	filepath     string
	workGroups   *gls.NumWorkGroups
	computeSpecs *renderer.ComputeSpecs
}

// Start is called once at the start of the demo.
func (t *ComputeDemo) Start(a *app.App) {

	// Adds directional front light
	dir1 := light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, 0.6)
	dir1.SetPosition(0, 0, 100)
	a.Scene().Add(dir1)

	// Add axes helper
	axes := helper.NewAxes(1)
	a.Scene().Add(axes)

	// Create custom shader
	a.Renderer().Coman().AddShader("computeShader", loadFile(a.DirData()+"/"+t.filepath))

	a.Renderer().Coman().AddProgram("ComputeProgram", "computeShader")

	callback := func(b *gls.BufferRAM, deltaTime time.Duration) {
		vec := b.GetVector4(0)
		fmt.Printf("0 %f %f %f %f\n", vec.X, vec.Y, vec.Z, vec.W)
		vec.X = 1.0
		err := b.SetVector4(0, vec)
		if err != nil {
			panic(err)
		}
	}
	ssbo := gls.NewSSBO(a.Renderer().Coman().GetGLS(), 0, gls.BO_DYNAMIC_COPY, gls.BO_READ_WRITE, callback, t.workGroups.X*16)
	ssbo.SetInitialData(slices.Repeat([]byte{0x11}, int(t.workGroups.X*16)))
	bufferObjects := gls.NewBufferObjects()
	bufferObjects.Set(ssbo)

	t.computeSpecs = renderer.NewComputeSpecs("ComputeProgram", "4_3", *gls.NewShaderDefines(), *bufferObjects)
	_, err := a.Renderer().Coman().SetProgram(t.computeSpecs)
	if err != nil {
		fmt.Printf("Failed to set the shader program: %s\n", err)
		return
	}
	a.Renderer().Coman().Compute(*t.workGroups, time.Duration(1))
}

// Update is called every frame.
func (t *ComputeDemo) Update(a *app.App, deltaTime time.Duration) {
	a.Renderer().Coman().SetProgram(t.computeSpecs)
	a.Renderer().Coman().Compute(*t.workGroups, time.Duration(1))
}

// Cleanup is called once at the end of the demo.
func (t *ComputeDemo) Cleanup(a *app.App) {}
