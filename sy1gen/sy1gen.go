package main

import (
	"flag"
	"fmt"

	"github.com/maxlambertini/sy1creator/sy1go"
)

func main() {
	//var initFileName = ""        // sy1 file to use as initialization
	//var morphFileName = ""       // sy1 file to use for morphing
	//var morphSteps = 1

	var dirPtr string

	var p *sy1go.SyPatchset // my patchset
	var filesToDelete []string

	var percentagePtr float64
	var geneticPtr bool
	var fullrandomPtr bool
	var impactPtr int
	var namePtr string
	var usagePtr bool

	flag.StringVar(&dirPtr, "directory", ".", "Output path")
	flag.StringVar(&dirPtr, "d", ".", "Output path (shorthand)")
	flag.Float64Var(&percentagePtr, "percentage", 0.5, "Deviation from current value, a float between 0 and 1")
	flag.Float64Var(&percentagePtr, "p", 0.5, "Deviation from current value, a float between 0 and 1 (shorthand)")
	flag.BoolVar(&geneticPtr, "genetic", false, "Activate genetic evolution from base path")
	flag.BoolVar(&geneticPtr, "g", false, "Activate genetic evolution from base path (shorthand)")
	flag.BoolVar(&fullrandomPtr, "fullrandom", false, "Enable full random mode")
	flag.BoolVar(&fullrandomPtr, "r", false, "Enable full random mode (shorthand)")
	flag.IntVar(&impactPtr, "impact", -1, "Number of parameter to change -1 all, max 80")
	flag.IntVar(&impactPtr, "i", -1, "Number of parameter to change -1 all, max 80 (shorthand)")
	flag.StringVar(&namePtr, "name", sy1go.Word(), "Define name of patchset")
	flag.StringVar(&namePtr, "n", sy1go.Word(), "Define name of patchset (shorthand)")
	flag.BoolVar(&usagePtr, "h", false, "Show help message (shorthand)")
	flag.BoolVar(&usagePtr, "help", false, "Show help message")

	flag.Parse()

	if !usagePtr {

		p = sy1go.InitNewPatchset("")

		defer func() {
			fmt.Println("Deleting temp files")
			p.DeleteTempSyFiles(filesToDelete)
			fmt.Println("\n\nThat's all!")
		}()

		fmt.Println("...Generating patch ")
		p.Directory = dirPtr
		p.Percentage = float32(percentagePtr)
		p.Genetic = geneticPtr
		p.FullRandom = fullrandomPtr
		p.Impact = impactPtr
		p.Name = namePtr
		p.UpdateValues()
		p.GeneratePatches()

		/*
		   if morphFileName != "":
		       p.MorphPatch = newSyPatch(morphFileName);
		       p.createMorphing (morphSteps)
		*/

		fmt.Println("...Generating patchset")
		filesToDelete = p.GenerateZip(namePtr + ".zip")
	} else {
		flag.Usage()
	}
}
