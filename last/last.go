// Copyright ©2012 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package last implements interaction with the last alignment tool.
// The last tool is available from http://last.cbrc.jp.
package last

import (
	"errors"
	"os/exec"

	"github.com/biogo/external"
)

var ErrMissingRequired = errors.New("last: missing required argument")

type DB struct {
	// Usage: lastdb [options] output-name fasta-sequence-file(s)
	// Prepare sequences for subsequent alignment with lastal.
	//
	// Main Options:
	//  -p: interpret the sequences as proteins
	//  -c: soft-mask lowercase letters
	//
	// Advanced Options (default settings):
	//  -Q: input format: 0=fasta, 1=fastq-sanger, 2=fastq-solexa, 3=fastq-illumina (0)
	//  -s: volume size (unlimited)
	//  -m: spaced seed pattern
	//  -u: subset seed file (yass.seed)
	//  -w: index step (1)
	//  -a: user-defined alphabet
	//  -b: bucket depth
	//  -x: just count sequences and letters
	//  -v: be verbose: write messages about what lastdb is doing
	//
	Cmd string `buildarg:"{{if .}}{{.}}{{else}}lastdb{{end}}"` // lastdb

	// Main Options:
	Protein  bool `buildarg:"{{if .}}-p{{end}}"` // -p: interpret the sequences as proteins
	Softmask bool `buildarg:"{{if .}}-c{{end}}"` // -c: soft-mask lowercase letters

	// Advanced Options:
	VolumeSize  int    `buildarg:"{{if .}}-s{{split}}{{.}}{{end}}"` // -s: volume size
	SeedPattern string `buildarg:"{{if .}}-m{{split}}{{.}}{{end}}"` // -m: spaced seed pattern
	HeaderFile  string `buildarg:"{{if .}}-u{{split}}{{.}}{{end}}"` // -u: subset seed file
	IndexStep   int    `buildarg:"{{if .}}-w{{split}}{{.}}{{end}}"` // -w: index step
	Alphabet    string `buildarg:"{{if .}}-a{{split}}{{.}}{{end}}"` // -a: user-defined alphabet
	BucketDepth int    `buildarg:"{{if .}}-b{{split}}{{.}}{{end}}"` // -b: bucket depth
	OnlyCount   bool   `buildarg:"{{if .}}-x{{end}}"`               // -x: just count sequences and letters
	Verbose     bool   `buildarg:"{{if .}}-v{{end}}"`               // -v: be verbose

	// Files:
	OutFile string   `buildarg:"{{.}}"`      // "<lastdb>"
	InFiles []string `buildarg:"{{args .}}"` // "<in.fa>"...
}

func (db DB) BuildCommand() (*exec.Cmd, error) {
	if db.OutFile == "" || len(db.InFiles) == 0 {
		return nil, ErrMissingRequired
	}
	cl := external.Must(external.Build(db))
	return exec.Command(cl[0], cl[1:]...), nil
}

type Align struct {
	// Usage: lastal [options] lastdb-name fasta-sequence-file(s)
	// Find local sequence alignments.
	//
	// Score options (default settings):
	//  -r: match score   (DNA: 1, protein: blosum62, 0<Q<5:  6)
	//  -q: mismatch cost (DNA: 1, protein: blosum62, 0<Q<5: 18)
	//  -p: file for residue pair scores
	//  -a: gap existence cost (DNA: 7, protein: 11, 0<Q<5: 21)
	//  -b: gap extension cost (DNA: 1, protein:  2, 0<Q<5:  9)
	//  -A: insertion existence cost (a)
	//  -B: insertion extension cost (b)
	//  -c: unaligned residue pair cost (100000)
	//  -F: frameshift cost (off)
	//  -x: maximum score drop for gapped alignments (max[y, a+b*20])
	//  -y: maximum score drop for gapless alignments (t*10)
	//  -z: maximum score drop for final gapped alignments (x)
	//  -d: minimum score for gapless alignments (e if j<2, else e*3/5)
	//  -e: minimum score for gapped alignments (DNA: 40, protein: 100, 0<Q<5: 180)
	//
	// Cosmetic options (default settings):
	//  -v: be verbose: write messages about what lastal is doing
	//  -o: output file
	//  -f: output format: 0=tabular, 1=maf (1)
	//
	// Miscellaneous options (default settings):
	//  -s: strand: 0=reverse, 1=forward, 2=both (2 for DNA, 1 for protein)
	//  -m: maximum multiplicity for initial matches (10)
	//  -l: minimum length for initial matches (1)
	//  -n: maximum number of gapless alignments per query position (infinity)
	//  -k: step-size along the query sequence (1)
	//  -i: query batch size (8 KiB, unless there are multiple lastdb volumes)
	//  -u: mask lowercase during extensions: 0=never, 1=gapless,
	//     2=gapless+gapped but not final, 3=always (2 if lastdb -c and Q<5, else 0)
	//  -w: supress repeats inside exact matches, offset by this distance or less (1000)
	//  -G: genetic code file
	//  -t: 'temperature' for calculating probabilities (1/lambda)
	//  -g: 'gamma' parameter for gamma-centroid and LAMA (1)
	//  -j: output type:
	//      0=match counts,
	//      1=gapless,
	//      2=redundant gapped,
	//      3=gapped,
	//      4=column ambiguity estimates,
	//      5=gamma-centroid,
	//      6=LAMA (3)
	//  -Q: input format:
	//      0=fasta,
	//      1=fastq-sanger,
	//      2=fastq-solexa,
	//      3=fastq-illumina,
	//      4=prb,
	//      5=PSSM (0)
	//
	Cmd string `buildarg:"{{if .}}{{.}}{{else}}lastal{{end}}"` // lastal

	// Score options:
	MatchScore     int    `buildarg:"{{if .}}-r{{split}}{{.}}{{end}}"` // -r: match score
	MismatchCost   int    `buildarg:"{{if .}}-q{{split}}{{.}}{{end}}"` // -q: mismatch cost
	ScoreFile      string `buildarg:"{{if .}}-p{{split}}{{.}}{{end}}"` // -p: file for residue pair scores
	GapCost        int    `buildarg:"{{if .}}-a{{split}}{{.}}{{end}}"` // -a: gap existence cost
	ExtendCost     int    `buildarg:"{{if .}}-b{{split}}{{.}}{{end}}"` // -b: gap extension cost
	UnalignedCost  int    `buildarg:"{{if .}}-c{{split}}{{.}}{{end}}"` // -c: unaligned residue pair cost
	FrameShiftCost int    `buildarg:"{{if .}}-F{{split}}{{.}}{{end}}"` // -F: frameshift cost (off)
	MaxGapDrop     int    `buildarg:"{{if .}}-x{{split}}{{.}}{{end}}"` // -x: max score drop for gapped
	MaxGaplessDrop int    `buildarg:"{{if .}}-y{{split}}{{.}}{{end}}"` // -y: max score drop for gapless
	MaxFinalDrop   int    `buildarg:"{{if .}}-z{{split}}{{.}}{{end}}"` // -z: max score drop for final gapped
	MinGapless     int    `buildarg:"{{if .}}-d{{split}}{{.}}{{end}}"` // -d: min score for gapless
	MinGapped      int    `buildarg:"{{if .}}-e{{split}}{{.}}{{end}}"` // -e: min score for gapped

	// Cosmetic options:
	Verbose bool   `buildarg:"{{if .}}-v{{end}}"`               // -v: be verbose
	OutFile string `buildarg:"{{if .}}-o{{split}}{{.}}{{end}}"` // -o: output file
	Tabular bool   `buildarg:"{{if .}}-f{{split}}0{{end}}"`     // -f: output format

	// Miscellaneous options:
	Strand      int     `buildarg:"{{if .}}-s{{split}}{{.}}{{end}}"` // -s: strand
	MaxMultiple int     `buildarg:"{{if .}}-m{{split}}{{.}}{{end}}"` // -m: max multiplicity for init matches
	MinSeed     int     `buildarg:"{{if .}}-l{{split}}{{.}}{{end}}"` // -l: min length for init matches
	MaxGapless  int     `buildarg:"{{if .}}-n{{split}}{{.}}{{end}}"` // -n: max number of gapless per query pos
	StepSize    int     `buildarg:"{{if .}}-k{{split}}{{.}}{{end}}"` // -k: step-size along the query seq
	BatchSize   int     `buildarg:"{{if .}}-i{{split}}{{.}}{{end}}"` // -i: query batch size
	MaskLower   int     `buildarg:"{{if .}}-u{{split}}{{.}}{{end}}"` // -u: mask lowercase during extensions
	SupressRep  int     `buildarg:"{{if .}}-w{{split}}{{.}}{{end}}"` // -w: supress repeats inside exact matches
	GenCodeFile string  `buildarg:"{{if .}}-G{{split}}{{.}}{{end}}"` // -G: genetic code file
	Temperature float64 `buildarg:"{{if .}}-t{{split}}{{.}}{{end}}"` // -t: 'temperature' for calculating probabilities
	Gamma       float64 `buildarg:"{{if .}}-g{{split}}{{.}}{{end}}"` // -g: 'gamma' parameter for gamma-centroid and LAMA
	OutputType  int     `buildarg:"{{if .}}-j{{split}}{{.}}{{end}}"` // -j: output type
	InFormat    int     `buildarg:"{{if .}}-Q{{split}}{{.}}{{end}}"` // -Q: input format

	// Files:
	DB      string   `buildarg:"{{.}}"`      // "<lastdb>"
	InFiles []string `buildarg:"{{args .}}"` // "<in.fa>"...
}

func (a Align) BuildCommand() (*exec.Cmd, error) {
	if a.DB == "" || len(a.InFiles) == 0 {
		return nil, ErrMissingRequired
	}
	cl := external.Must(external.Build(a))
	return exec.Command(cl[0], cl[1:]...), nil
}

type Expect struct {
	// Usage: lastex [options] reference-counts-file query-counts-file [alignments-file]
	// Calculate expected numbers of alignments for random sequences.
	//
	// Options (default settings):
	//  -s: strands (2 for DNA, 1 for protein)
	//  -r: match score   (DNA: 1, protein: blosum62)
	//  -q: mismatch cost (DNA: 1, protein: blosum62)
	//  -p: file for residue pair scores
	//  -a: gap existence cost (DNA: 7, protein: 11)
	//  -b: gap extension cost (DNA: 1, protein:  2)
	//  -g: do calculations for gapless alignments
	//  -y: find the expected number of alignments with score >= this
	//  -E: maximum expected number of alignments
	//  -z: calculate the expected number of alignments per:
	//      0 = reference counts file / query counts file
	//      1 = reference counts file / each query sequence
	//      2 = each reference sequence / query counts file
	//      3 = each reference sequence / each query sequence (0)
	//
	Cmd string `buildarg:"{{if .}}{{.}}{{else}}lastex{{end}}"` // lastex

	// Options:
	Strand       int    `buildarg:"{{if .}}-s{{split}}{{.}}{{end}}"` // -s: strands
	MatchScore   int    `buildarg:"{{if .}}-r{{split}}{{.}}{{end}}"` // -r: match score
	MismatchCost int    `buildarg:"{{if .}}-q{{split}}{{.}}{{end}}"` // -q: mismatch cost
	ScoreFile    string `buildarg:"{{if .}}-p{{split}}{{.}}{{end}}"` // -p: file for residue pair scores
	GapCost      int    `buildarg:"{{if .}}-a{{split}}{{.}}{{end}}"` // -a: gap existence cost
	ExtendCost   int    `buildarg:"{{if .}}-b{{split}}{{.}}{{end}}"` // -b: gap extension cost
	DoGapless    bool   `buildarg:"{{if .}}-g{{end}}"`               // -g: do calculations for gapless
	FindThresh   int    `buildarg:"{{if .}}-y{{split}}{{.}}{{end}}"` // -y: find alignments with score >= this
	MaxExpected  int    `buildarg:"{{if .}}-E{{split}}{{.}}{{end}}"` // -E: maximum expected number
	Calculate    int    `buildarg:"{{if .}}-z{{split}}{{.}}{{end}}"` // -z: calculate expected alignments

	// Files:
	Ref        string   `buildarg:"{{.}}"`      // "<lastdb>"
	Query      string   `buildarg:"{{.}}"`      // "<lastdb>"
	AlignFiles []string `buildarg:"{{args .}}"` // "<in.maf>"...
}

func (e Expect) BuildCommand() (*exec.Cmd, error) {
	if e.Ref == "" || e.Query == "" {
		return nil, ErrMissingRequired
	}
	cl := external.Must(external.Build(e))
	return exec.Command(cl[0], cl[1:]...), nil
}
