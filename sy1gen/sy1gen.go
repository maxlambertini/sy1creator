package main

import (
	"fmt"
	"github.com/maxlambertini/sy1creator/sy1go"
)

func main() {
	var p *sy1go.SyPatchset      // my patchset
	var directory = "./"         // default output directory (--directory, -d )
	var percentage float32 = 0.5 // deviation from current value (--percentage, -p)
	var genetic = true           // set result as new default (--genetic, -g)
	var fullrandom = false       // randomization or parametric (--fullrandom -f)
	var impact = -1              // number of parameters to process, -1 all (--impact, -i)
	var name = sy1go.Word()      // patchset name (--name, -n)
	//var initFileName = ""        // sy1 file to use as initialization
	//var morphFileName = ""       // sy1 file to use for morphing
	//var morphSteps = 1

	p = sy1go.InitNewPatchset("")
	fmt.Println("Generating patch ")
	p.Directory = directory
	p.Percentage = percentage
	p.Genetic = genetic
	p.FullRandom = fullrandom
	p.Impact = impact
	p.Name = name
	p.UpdateValues()
	p.GeneratePatches()

	p.UpdateValues()
	p.GeneratePatches()
	/*
	   if morphFileName != "":
	       p.MorphPatch = newSyPatch(morphFileName);
	       p.createMorphing (morphSteps)
	*/

	fmt.Println("Generating patchset:\n", p)
	p.GenerateZip(name + ".zip")
	fmt.Println("o \n\nThat's all!")
}
