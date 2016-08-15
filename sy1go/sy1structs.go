package sy1go

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jhoonb/archivex"
)

//Constants ...
const (
	FTWindows   = 1
	FTUnix      = 2
	DEFAULTDATA = `METAL GTR
color=magenta
ver=112
0,1
45,0
76,5
1,3
2,39
3,8
4,1
5,40
6,1
7,1
8,89
9,0
10,0
11,123
12,50
13,77
71,2
72,78
91,0
95,0
96,1
97,1
14,2
15,49
16,35
17,59
18,82
19,72
20,70
21,68
22,95
23,119
24,0
25,23
26,105
27,100
28,71
29,75
30,11
59,0
31,4
32,1
33,12
34,54
65,0
82,0
35,9
83,105
36,105
98,64
37,38
66,0
64,1
52,100
53,110
54,10
55,122
56,90
60,83
61,72
62,46
63,46
90,77
77,0
78,1
79,40
80,70
81,45
38,0
94,16
39,40
74,0
73,0
93,4
75,127
84,4
85,2
92,0
40,2
86,45057
50,64
87,44
88,45057
51,91
89,43
57,0
41,3
42,0
43,50
44,125
67,0
68,0
58,0
46,6
47,1
48,67
49,112
69,0
70,0`
)

var colors = []string{"cyan", "red", "green", "blue", "yellow", "magenta"}
var vow = []string{"a", "e", "i", "o", "u", "y"}
var con = []string{"b", "c", "d", "f", "g", "k", "l", "m", "n", "p", "r", "s", "t", "v", "x", "z"}

// GetColor ...
func GetColor() string { return "color=" + colors[rand.Intn(len(colors))] }

// GetVer ...
func GetVer() string { return "ver=113" }

// Word ...
func Word() string {
	var x = rand.Intn(2)
	var lung = 3 + rand.Intn(5)
	var s = ""
	var gu = ""
	for idx := 0; idx < lung; idx++ {
		if x%2 == 0 {
			gu = vow[rand.Intn(len(vow))]
		} else {
			gu = con[rand.Intn(len(con))]
		}
		if rand.Intn(9) == 2 {
			x = x + 1
		}
		s = s + gu
	}
	return s
}

//SyParam ...
type SyParam struct {
	Description string
	index       int
	vmin        int
	vmax        int
	sdefault    int
	current     int
}

// InitSyParam ...
func InitSyParam(_index int, _vmin int, _vmax int, _desc string) *SyParam {
	var sy = SyParam{
		Description: _desc,
		index:       _index,
		vmin:        _vmin,
		vmax:        _vmax,
		sdefault:    _vmin,
		current:     _vmin,
	}
	return &sy
}

// RandomizeParam ...
func (p *SyParam) RandomizeParam(genetic bool) {
	switch p.index {
	case 9:
		p.current = 0
	case 29:
		p.current = 100
	default:
		if p.vmax > 1 {
			var d = p.vmax - p.vmin
			var h = p.vmin + rand.Intn(d+1)
			p.current = h
		} else {
			p.current = rand.Intn(100) % 2
		}
	}
	if genetic {
		p.sdefault = p.current
	}
}

// ParamString ...
func (p *SyParam) ParamString() (res string) {
	res = fmt.Sprintf("%d,%d", p.index, p.current)
	return
}

// Deviation ...
func (p *SyParam) Deviation(percentage float32, genetic bool) {
	switch p.index {
	case 9:
		p.current = 0
	case 29:
		p.current = 100
	default:
		var res int
		if p.vmax > 4 {
			d := (int)(percentage * float32(p.vmax-p.vmin))
			if d < 1 {
				d = 1
			}
			var rng = rand.Intn(1+d*2) - d
			if ((p.vmin + rng + p.sdefault) < p.vmin) || ((p.vmin + p.sdefault + rng) > p.vmax) {
				res = p.vmin + p.sdefault - rng
			} else {
				res = p.vmin + p.sdefault + rng
			}
			p.current = res

		} else {
			v1 := percentage * 100.0
			v2 := int(v1)
			if rand.Intn(100) < v2 {
				p.RandomizeParam(genetic)
			}
		}
	}
	if genetic {
		p.sdefault = p.current
	}
}

//SyPatch ...
type SyPatch struct {
	Description string
	Version     string
	Color       string
	Params      map[int]*SyParam
	Percentage  float32
	Deviation   float32
	Genetic     bool
	Impact      int
	Directory   string
	Patchset    *SyPatchset
}

// InitSyPatch ...
func InitSyPatch(filename string) *SyPatch {
	if filename != "" {
		fmt.Println("Initializing with ", filename)
	}
	var sp = SyPatch{
		Description: Word(),
		Version:     GetVer(),
		Color:       GetColor(),
		Params:      initWithDefaultParams(filename),
		Percentage:  -1.0,
		Genetic:     false,
	}
	return &sp
}

// GenerateRandomPatch ...
func (p *SyPatch) GenerateRandomPatch() {
	for _, v := range p.Params {
		//v1 := v.current
		v.RandomizeParam(p.Genetic)
		//v2 := v.current
		//fmt.Println("Patch ", p.Description, " param ", v.Description, "Orig val=", v1, " ch val=", v2)
	}
}

// CopyPatchValues copies the values of a Patch to another Patch
func (p *SyPatch) CopyPatchValues(pDest *SyPatch) {
	for k, v := range p.Params {
		_, ok := pDest.Params[k]
		if ok {
			pDest.Params[k].sdefault = v.sdefault
			pDest.Params[k].current = v.current
		}
	}
}

// GenerateParametricPatch ...
func (p *SyPatch) GenerateParametricPatch() {
	if p.Impact > 0 {
		var s []int
		for k := range p.Params {
			s = append(s, k)
		}
		list := rand.Perm(len(s))
		for w := 0; w < p.Impact; w++ {
			ix := list[w]
			newk := s[ix]
			pa, ok := p.Params[newk]
			if ok {
				pa.Deviation(p.Percentage, p.Genetic)
			}
		}
	} else {
		for _, v := range p.Params {
			v.Deviation(p.Percentage, p.Genetic)
		}
	}
}

// GeneratePatchText ...
func (p *SyPatch) GeneratePatchText() string {
	newLine := "\r\n"
	sbuf := p.Description + newLine + p.Color + newLine + p.Version + newLine
	for h := 0; h < len(CorrectSequence); h++ {
		k := CorrectSequence[h]
		v, ok := p.Params[k]
		if ok {
			sbuf = sbuf + v.ParamString() + newLine
		}
	}
	return sbuf
}

// GetPatchFileName ...
func GetPatchFileName(i int) string {
	return fmt.Sprintf("%03d.sy1", i)
}

func checkFileErr(e error) {
	if e != nil {
		panic(e)
	}
}

// GeneratePatchFile ...
func (p *SyPatch) GeneratePatchFile(index int) string {
	fp := GetPatchFileName(index)
	fullPath := p.Directory + string(os.PathSeparator) + fp
	pt := p.GeneratePatchText()
	err := ioutil.WriteFile(fullPath, []byte(pt), 0644)
	checkFileErr(err)
	return fp
}

//SyPatchset ...
type SyPatchset struct {
	Name         string
	Directory    string
	Percentage   float32
	Genetic      bool
	Patches      []*SyPatch
	Impact       int
	FullRandom   bool
	DefaultPatch *SyPatch
	MorphPatch   *SyPatch
}

// InitNewPatchset ...
func InitNewPatchset(filename string) *SyPatchset {
	p := SyPatchset{
		Directory:    "./",
		Percentage:   0.5,
		Genetic:      false,
		Patches:      make([]*SyPatch, 10),
		Name:         Word(),
		Impact:       -1,
		DefaultPatch: InitSyPatch(filename),
	}
	tmpPatches := make([]*SyPatch, 128)
	for w := 0; w < 128; w++ {
		pa := InitSyPatch("")
		pa.Patchset = &p
		p.DefaultPatch.CopyPatchValues(pa)
		//CopyPatchValues(p.DefaultPatch,pa)
		tmpPatches[w] = pa
	}
	p.Patches = tmpPatches
	return &p
}

// CreateMorphing ...
func (p *SyPatchset) CreateMorphing(steps int) {
	var s []int
	dict := make(map[int]float32)
	var flVal float32
	var curDif float32
	var p1 *SyParam
	var p2 *SyParam
	var ok bool
	for k := range p.DefaultPatch.Params {
		s = append(s, k)
	}
	for k := range s {
		p1, ok = p.DefaultPatch.Params[k]
		if ok {
			p2 = p.MorphPatch.Params[k]
			curDif = float32(p2.current-p1.current) / 24.0
			dict[k] = curDif
		}
	}
	for st := 0; st < steps+2; st++ {
		pDest := p.Patches[st]
		for k, v := range dict {
			val, ok := pDest.Params[k]
			if ok {
				flVal = float32(p.DefaultPatch.Params[k].current) + v*float32(st)
				val.current = int(flVal)

			}
		}
	}
}

// UpdateValues ...
func (p *SyPatchset) UpdateValues() {
	for _, pa := range p.Patches {
		pa.Directory = p.Directory
		pa.Percentage = p.Percentage
		pa.Genetic = p.Genetic
		pa.Impact = p.Impact
		pa.Description = Word()
	}
}

// GeneratePatches ...
func (p *SyPatchset) GeneratePatches() {
	for _, patch := range p.Patches {
		if p.FullRandom {
			patch.GenerateRandomPatch()
		} else {
			patch.GenerateParametricPatch()
		}
	}

}

// GenerateZip ...
func (p *SyPatchset) GenerateZip(filename string) []string {
	zip := new(archivex.ZipFile)
	defer zip.Close()
	zipFileName := p.Directory + string(os.PathSeparator) + filename
	zip.Create(zipFileName)
	var filesToDelete []string
	for h, patch := range p.Patches {
		fp := patch.GeneratePatchFile(1 + h)
		zip.AddFile(fp)
		filesToDelete = append(filesToDelete, p.Directory+string(os.PathSeparator)+fp)
	}
	// deleting file
	return filesToDelete
}

// DeleteTempSyFiles ...
func (p *SyPatchset) DeleteTempSyFiles(filesToDelete []string) {
	for _, f2d := range filesToDelete {
		//fmt.Println("Removing ", f2d)
		err := os.Remove(f2d)
		if err != nil {
			fmt.Println("Error deleting file ", f2d, ". Error: ", err)
		}

	}
}

// ---------------------------------- helper methods --------------------------

func readLines(path string) ([]string, error) {
	fmt.Println("Opening ", path)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

var CorrectSequence = []int{0,
	45,
	76,
	1,
	2,
	3,
	4,
	5,
	6,
	7,
	8,
	9,
	10,
	11,
	12,
	13,
	71,
	72,
	91,
	95,
	96,
	97,
	14,
	15,
	16,
	17,
	18,
	19,
	20,
	21,
	22,
	23,
	24,
	25,
	26,
	27,
	28,
	29,
	30,
	59,
	31,
	32,
	33,
	34,
	65,
	82,
	35,
	83,
	36,
	98,
	37,
	66,
	64,
	52,
	53,
	54,
	55,
	56,
	60,
	61,
	62,
	63,
	90,
	77,
	78,
	79,
	80,
	81,
	38,
	94,
	39,
	74,
	73,
	93,
	75,
	84,
	85,
	92,
	40,
	86,
	50,
	87,
	88,
	51,
	89,
	57,
	41,
	42,
	43,
	44,
	67,
	68,
	58,
	46,
	47,
	48,
	49,
	69,
	70,
}

func initWithDefaultParams(filename string) (result map[int]*SyParam) {
	dict := make(map[int]*SyParam, 128)

	dict[0] = InitSyParam(0, 0, 3, "OSCILLATOR 1 WAVE")
	dict[45] = InitSyParam(45, 0, 127, "FM")
	dict[76] = InitSyParam(76, 0, 127, "DETUNE")
	dict[1] = InitSyParam(1, 1, 4, "OSCILLATOR 2 WAVE")
	dict[2] = InitSyParam(2, 0, 127, "PITCH")
	dict[3] = InitSyParam(3, 0, 127, "FINE")
	dict[4] = InitSyParam(4, 0, 1, "TRACK")
	dict[5] = InitSyParam(5, 0, 127, "MIX")
	dict[6] = InitSyParam(6, 0, 1, "SYNC")
	dict[7] = InitSyParam(7, 0, 1, "RING")
	dict[8] = InitSyParam(8, 0, 127, "PW")
	dict[9] = InitSyParam(9, -24, 24, "TRANSPOSE")
	dict[10] = InitSyParam(10, 0, 1, "M. ENV SWITCH")
	dict[11] = InitSyParam(11, 0, 127, "AMOUNT")
	dict[12] = InitSyParam(12, 0, 127, "ATTACK")
	dict[13] = InitSyParam(13, 0, 127, "DECAY")
	dict[71] = InitSyParam(71, 0, 2, "DEST")
	dict[72] = InitSyParam(72, 0, 127, "TUNE")
	dict[14] = InitSyParam(14, 0, 3, "FILTER TYPE")
	dict[15] = InitSyParam(15, 0, 127, "ATTACK")
	dict[16] = InitSyParam(16, 0, 127, "DECAY")
	dict[17] = InitSyParam(17, 0, 127, "SUSTAIN")
	dict[18] = InitSyParam(18, 0, 127, "RELEASE")
	dict[19] = InitSyParam(19, 0, 127, "FREQUENCY")
	dict[20] = InitSyParam(20, 0, 127, "RESONANCE")
	dict[21] = InitSyParam(21, 0, 127, "AMOUNT")
	dict[22] = InitSyParam(22, 0, 127, "TRACK")
	dict[23] = InitSyParam(23, 0, 127, "SATURATION")
	dict[24] = InitSyParam(24, 0, 1, "VELOCITY AMOUNT")
	dict[25] = InitSyParam(25, 0, 127, "ATTACK")
	dict[26] = InitSyParam(26, 0, 127, "DECAY")
	dict[27] = InitSyParam(27, 0, 127, "SUSTAIN")
	dict[28] = InitSyParam(28, 0, 127, "RELEASE")
	dict[29] = InitSyParam(29, 0, 127, "GAIN")
	dict[30] = InitSyParam(30, 0, 127, "VELOCITY AMOUNT")
	dict[59] = InitSyParam(59, 0, 1, "ARP SWITCH")
	dict[31] = InitSyParam(31, 1, 4, "TYPE")
	dict[32] = InitSyParam(32, 0, 3, "RANGE")
	dict[33] = InitSyParam(33, 0, 18, "BEAT")
	dict[34] = InitSyParam(34, 5, 127, "GATE")
	dict[65] = InitSyParam(65, 0, 1, "DELAY SWITCH")
	dict[82] = InitSyParam(82, 0, 2, "TYPE")
	dict[35] = InitSyParam(35, 0, 19, "TIME")
	dict[83] = InitSyParam(83, 0, 127, "SPREAD")
	dict[36] = InitSyParam(36, 1, 120, "FEEDBACK")
	dict[37] = InitSyParam(37, 0, 127, "DRY/WET")
	dict[66] = InitSyParam(66, 0, 1, "CHORUS SWITCH")
	dict[64] = InitSyParam(64, 1, 4, "TYPE")
	dict[52] = InitSyParam(52, 0, 127, "TIME")
	dict[53] = InitSyParam(53, 0, 127, "DEPTH")
	dict[54] = InitSyParam(54, 0, 127, "RATE")
	dict[55] = InitSyParam(55, 0, 127, "FEEDBACK")
	dict[56] = InitSyParam(56, 0, 127, "LEVEL")
	dict[60] = InitSyParam(60, 0, 127, "EQ TONE")
	dict[61] = InitSyParam(61, 0, 127, "FREQUENCY")
	dict[62] = InitSyParam(62, 0, 127, "LEVEL")
	dict[63] = InitSyParam(63, 0, 127, "Q")
	dict[90] = InitSyParam(90, 32, 96, "L-R")
	dict[77] = InitSyParam(77, 0, 1, "EFFECT SWITCH")
	dict[78] = InitSyParam(78, 0, 9, "TYPE")
	dict[79] = InitSyParam(79, 0, 127, "CTRL1")
	dict[80] = InitSyParam(80, 0, 127, "CTRL2")
	dict[81] = InitSyParam(81, 0, 127, "LEVEL")
	dict[38] = InitSyParam(38, 0, 2, "PLAY MODE")
	dict[39] = InitSyParam(39, 0, 127, "PORTAMENTO")
	dict[74] = InitSyParam(74, 0, 1, "AUTO")
	dict[40] = InitSyParam(40, 0, 24, "PB RANGE")
	dict[73] = InitSyParam(73, 0, 1, "UNISON")
	dict[75] = InitSyParam(75, 0, 127, "DETUNE")
	dict[84] = InitSyParam(84, 0, 127, "SPREAD")
	dict[85] = InitSyParam(85, 0, 48, "PITCH")
	dict[50] = InitSyParam(50, 0, 127, "LFO1 WHEEL SENS")
	dict[51] = InitSyParam(51, 0, 127, "SPEED")
	dict[57] = InitSyParam(57, 0, 1, "LFO1 SWITCH")
	dict[41] = InitSyParam(41, 1, 7, "DEST")
	dict[42] = InitSyParam(42, 0, 4, "WAVEFORM")
	dict[43] = InitSyParam(43, 0, 127, "SPEED")
	dict[44] = InitSyParam(44, 0, 127, "AMOUNT")
	dict[67] = InitSyParam(67, 0, 1, "TEMPO SYNC")
	dict[68] = InitSyParam(68, 0, 1, "KEY SYNC")
	dict[58] = InitSyParam(58, 0, 1, "LFO2 SWITCH")
	dict[46] = InitSyParam(46, 1, 7, "DEST")
	dict[47] = InitSyParam(47, 0, 4, "WAVEFORM")
	dict[48] = InitSyParam(48, 0, 127, "SPEED")
	dict[49] = InitSyParam(49, 0, 127, "AMOUNT")
	dict[69] = InitSyParam(69, 0, 1, "TEMPO SYNC")
	dict[70] = InitSyParam(70, 0, 1, "KEY SYNC")
	dict[91] = InitSyParam(91, 0, 0, "AA")
	dict[95] = InitSyParam(94, 1, 16, "FF")
	dict[96] = InitSyParam(93, 1, 2, "XX")
	dict[95] = InitSyParam(95, 0, 0, "BB")
	dict[96] = InitSyParam(96, 1, 1, "CC")
	dict[97] = InitSyParam(97, 1, 1, "DD")
	dict[98] = InitSyParam(97, 1, 64, "EE")
	dict[98] = InitSyParam(92, 0, 1, "WW")
	dict[86] = InitSyParam(86, 45057, 45057, "E1E")
	dict[88] = InitSyParam(88, 45057, 45057, "W1W")

	if filename != "" {
		fmt.Println("Initializing with ", filename)
		sylines, err := readLines(filename)
		if err == nil {
			//if err != nil {
			syData := sylines[3:]
			for _, line := range syData {
				//fmt.Println(line)
				line = strings.TrimSpace(line)
				tokens := strings.Split(line, ",")
				if len(tokens) > 0 {
					//fmt.Println("line: ", line, "tokens: ", tokens)
					iKey, err := strconv.Atoi(tokens[0])
					if err == nil {
						iVal, err := strconv.Atoi(tokens[1])
						if err == nil {
							_, ok := dict[iKey]
							if ok {
								dict[iKey].sdefault = iVal
								dict[iKey].current = iVal
							}
						}
					}

				}
			}
		} else {
			fmt.Println("Error reading patch file ", err)
		}
	} else {
		sylines := strings.Split(DEFAULTDATA, "\n")
		syData := sylines[3:]
		for _, line := range syData {
			line = strings.TrimSpace(line)
			tokens := strings.Split(line, ",")
			if len(tokens) > 0 {
				iKey, err := strconv.Atoi(tokens[0])
				if err == nil {
					iVal, err := strconv.Atoi(tokens[1])
					if err == nil {
						_, ok := dict[iKey]
						if ok {
							dict[iKey].sdefault = iVal
							dict[iKey].current = iVal
						}
					}
				}

			}
		}

	}
	result = dict
	return
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
