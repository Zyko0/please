package patch

import (
	"reflect"

	"github.com/agiledragon/gomonkey/v2"
)

type applyFunc func(p *gomonkey.Patches)

type Patch struct {
	apply   applyFunc
	patches *gomonkey.Patches
}

func NewPatchFunc(target, replace any) *Patch {
	p := &Patch{
		patches: gomonkey.NewPatches(),
	}
	p.apply = func(mp *gomonkey.Patches) {
		mp.ApplyFunc(target, replace)
	}

	return p
}

func NewPatchMethod(obj any, method string, replace any) *Patch {
	p := &Patch{
		patches: gomonkey.NewPatches(),
	}
	p.apply = func(mp *gomonkey.Patches) {
		mp.ApplyMethod(reflect.TypeOf(obj), method, replace)
	}

	return p
}

func (p *Patch) Enable() {
	if p != nil {
		p.apply(p.patches)
	}
}

func (p *Patch) Disable() {
	if p != nil {
		p.patches.Reset()
	}
}
