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
	fmt.Printf("Size: %d\n", gls.SizeMat4x3Std430)

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
		//for _, by := range b.AsBytes() {
		//	fmt.Printf("%02x", by)
		//}
		vec, err := b.GetDvec4(0)
		if err != nil {
			panic(err)
		}
		fmt.Printf("0 %f %f %f %f\n", vec.X, vec.Y, vec.Z, vec.W)
		vec.X = 1.0
		err = b.SetDvec4(0, vec)
		if err != nil {
			panic(err)
		}
	}
	callbackIntegers := func(b *gls.BufferRAM, deltaTime time.Duration) {
		//for _, by := range b.AsBytes() {
		//	fmt.Printf("%02x", by)
		//}
		i, err := b.GetInt(0)
		if err != nil {
			panic(err)
		}
		//fmt.Printf("\t0 %d\n", i)
		i += 1
		err = b.SetInt(0, i)
		if err != nil {
			panic(err)
		}
	}
	callbackBools := func(b *gls.BufferRAM, deltaTime time.Duration) {
		//for _, by := range b.AsBytes() {
		//	fmt.Printf("%02x", by)
		//}
		i, err := b.GetBool(0)
		if err != nil {
			panic(err)
		}
		s := "true"
		if !i {
			s = "false"
		}
		fmt.Printf("\t0 %s\n", s)
		err = b.SetBool(0, !i)
		if err != nil {
			panic(err)
		}
		for _, b_ := range b.AsBool() {
			s := "1"
			if !b_ {
				s = "0"
			}
			fmt.Printf("%s", s)
		}
		fmt.Println("")
	}
	bufferObjects := gls.NewBufferObjects()
	bufferObjects.Set(
		gls.NewSSBO(a.Renderer().Coman().GetGLS(), 0,
			gls.BO_DYNAMIC_COPY, gls.BO_READ_WRITE, callback, gls.SizeVec4Std430*16).SetInitialData(slices.Repeat([]byte{0x00}, int(gls.SizeVec4Std430*16))))
	bufferObjects.Set(
		gls.NewSSBO(a.Renderer().Coman().GetGLS(), 1,
			gls.BO_DYNAMIC_COPY, gls.BO_READ_WRITE, callbackIntegers, gls.SizeIntStd430*16).SetInitialData(slices.Repeat([]byte{0x00}, int(gls.SizeIntStd430*16))))
	bufferObjects.Set(
		gls.NewSSBO(a.Renderer().Coman().GetGLS(), 2,
			gls.BO_DYNAMIC_COPY, gls.BO_READ_WRITE, callbackBools, gls.SizeBoolStd430*16).SetInitialData(slices.Repeat([]byte{0x00}, int(gls.SizeBoolStd430*16))))

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
