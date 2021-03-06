// Copyright ©2012 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mafft implements interaction with the MAFFT multiple alignment tool.
// MAFFT is available from http://mafft.cbrc.jp/alignment/software/
package mafft

import (
	"os/exec"

	"github.com/biogo/external"
)

type Mafft struct {
	// Usage: mafft <inputfile> > <outputfile>
	//
	// For details relating to options and parameters, see the MAFFT manual.
	//
	Cmd string `buildarg:"{{if .}}{{.}}{{else}}mafft{{end}}"` // mafft

	// Algorithm:
	Auto          bool    `buildarg:"{{if .}}--auto{{end}}"`                     // --auto
	HexamerPair   bool    `buildarg:"{{if .}}--6merpair{{end}}"`                 // --6merpair
	GlobalPair    bool    `buildarg:"{{if .}}--globalpair{{end}}"`               // --globalpair
	LocalPair     bool    `buildarg:"{{if .}}--localpair{{end}}"`                // --localpair
	GenafPair     bool    `buildarg:"{{if .}}--genafpair{{end}}"`                // --genafpair
	FastaPair     bool    `buildarg:"{{if .}}--fastapair{{end}}"`                // --fastapair
	Weighting     float64 `buildarg:"{{if .}}--weighti{{split}}{{.}}{{end}}"`    // --weighti <f.>
	ReTree        int     `buildarg:"{{if .}}--retree{{split}}{{.}}{{end}}"`     // --retree <n>
	MaxIterate    int     `buildarg:"{{if .}}--maxiterate{{split}}{{.}}{{end}}"` // --maxiterate <n>
	Fft           bool    `buildarg:"{{if .}}--fft{{end}}"`                      // --fft
	NoFft         bool    `buildarg:"{{if .}}--nofft{{end}}"`                    // --nofft
	NoScore       bool    `buildarg:"{{if .}}--noscore{{end}}"`                  // --noscore
	MemSave       bool    `buildarg:"{{if .}}--memsave{{end}}"`                  // --memsave
	Partree       bool    `buildarg:"{{if .}}--parttree{{end}}"`                 // --parttree
	DPPartTree    bool    `buildarg:"{{if .}}--dpparttree{{end}}"`               // --dpparttree
	FastaPartTree bool    `buildarg:"{{if .}}--fastaparttree{{end}}"`            // --fastaparttree
	PartSize      int     `buildarg:"{{if .}}--partsize{{split}}{{.}}{{end}}"`   // --partsize <n>
	GroupSize     int     `buildarg:"{{if .}}--groupsize{{split}}{{.}}{{end}}"`  // --groupsize <n>

	// Parameter:
	GapOpenCost          float64 `buildarg:"{{if .}}--op{{split}}{{.}}{{end}}"`       // --op <f.>
	ExtensionCost        float64 `buildarg:"{{if .}}--ep{{split}}{{.}}{{end}}"`       // --ep <f.>
	LocalOpenCost        float64 `buildarg:"{{if .}}--lop{{split}}{{.}}{{end}}"`      // --lop <f.>
	LocalPairOffset      float64 `buildarg:"{{if .}}--lep{{split}}{{.}}{{end}}"`      // --lep <f.>
	LocalExtensionCost   float64 `buildarg:"{{if .}}--lexp{{split}}{{.}}{{end}}"`     // --lexp <f.>
	GapOpenSkipCost      float64 `buildarg:"{{if .}}--LOP{{split}}{{.}}{{end}}"`      // --LOP <f.>
	GapExtensionSkipCost float64 `buildarg:"{{if .}}--LEXP{{split}}{{.}}{{end}}"`     // --LEXP <f.>
	Blosum               byte    `buildarg:"{{if .}}--bl{{split}}{{.}}{{end}}"`       // --bl <n>
	JttPAM               uint    `buildarg:"{{if .}}--jtt{{split}}{{.}}{{end}}"`      // --jtt <n>
	TransMembranePAM     uint    `buildarg:"{{if .}}--tm{{split}}{{.}}{{end}}"`       // --tm <n>
	AminoMatrix          string  `buildarg:"{{if .}}--aamatrix{{split}}{{.}}{{end}}"` // --aamatrix <file>
	FModel               bool    `buildarg:"{{if .}}--fmodel{{end}}"`                 // --fmodel

	// Output:
	ClustalOut bool `buildarg:"{{if .}}--clustalout{{end}}"` // --clustalout
	InputOrder bool `buildarg:"{{if .}}--inputorder{{end}}"` // --inputorder
	Reorder    bool `buildarg:"{{if .}}--reorder{{end}}"`    // --reorder
	TreeOut    bool `buildarg:"{{if .}}--treeout{{end}}"`    // --treeout
	Quiet      bool `buildarg:"{{if .}}--quiet{{end}}"`      // --quiet

	// Input:
	Nucleic bool     `buildarg:"{{if .}}--nuc{{end}}"`                                 // --nuc
	Amino   bool     `buildarg:"{{if .}}--amino{{end}}"`                               // --amino
	Seed    []string `buildarg:"{{if .}}{{mprintf \"--seed\x00%s\" . | args}}{{end}}"` // --seed <file>...

	// Performance:
	Threads int `buildarg:"{{if .}}--thread{{split}}{{.}}{{end}}"` // --thread <n>

	// Files:
	InFile string `buildarg:"{{if .}}{{.}}{{else}}-{{end}}"` // <inputfile> - default to Stdin.
}

func (m Mafft) BuildCommand() (*exec.Cmd, error) {
	cl := external.Must(external.Build(m))
	return exec.Command(cl[0], cl[1:]...), nil
}
