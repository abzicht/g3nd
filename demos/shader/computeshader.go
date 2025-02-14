package shader

import (
	"encoding/hex"
	"fmt"
	"os"
	"time"
	"unsafe"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/g3nd/app"
)

func init() {
	app.DemoMap["shader.computeshader"] = &ComputeDemo{filepath: "shaders/compute2.glsl", workGroups: gls.NewNumWorkGroups(64, 1, 1)}
}

func loadFile(filepath string) string {
	b, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type ComputeDemo struct {
	a          *app.App
	filepath   string
	workGroups *gls.NumWorkGroups
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
	a.Coman().AddShader("computeShader", loadFile(a.DirData()+"/"+t.filepath))

	a.Coman().AddProgram("ComputeProgram", "computeShader")

	callback := func(ssbo *gls.SSBO, p unsafe.Pointer, deltaTime time.Duration) {
		fmt.Print("Raw data: ")
		for i := 0; i < ssbo.Length*4; i++ {
			var element byte = *(*byte)(unsafe.Pointer(uintptr(p) + uintptr(i)))
			fmt.Print(hex.EncodeToString([]byte{element}))
		}
		fmt.Print("\n")
		for i := 0; i < ssbo.Length; i++ {
			// are writing go or are we writing c???
			// from the unsafe docs:
			// > e := unsafe.Pointer(&p[i] + i * sizeof(p))
			var element float32 = *(*float32)(unsafe.Pointer(uintptr(p) + uintptr(i)*4))
			fmt.Printf("Value at %d: %f\n", i, element)
		}
		fmt.Print("\n")
	}
	ssbo := gls.NewSSBO(a.Coman().GetGLS(), 0, gls.SSBO_READ_ONLY, callback, 0, 4)
	bufferObjects := gls.NewBufferObjects()
	bufferObjects.Set(ssbo)

	computeSpecs := renderer.NewComputeSpecs("ComputeProgram", "4_3", *gls.NewShaderDefines(), *bufferObjects)
	_, err := a.Coman().SetProgram(computeSpecs)
	if err != nil {
		fmt.Printf("Failed to set the shader program: %s\n", err)
		return
	}
	a.Coman().Compute(*t.workGroups, time.Duration(1))
}

// Update is called every frame.
func (t *ComputeDemo) Update(a *app.App, deltaTime time.Duration) {
}

// Cleanup is called once at the end of the demo.
func (t *ComputeDemo) Cleanup(a *app.App) {}
