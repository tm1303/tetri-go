package main


var (

	black   = "\033[40m " //[]byte{keyEscape, '[', '3', '1', 'm'},
	red     = "\033[41m " //[]byte{keyEscape, '[', '3', '1', 'm'},
	green   = "\033[42m " //[]byte{keyEscape, '[', '3', '2', 'm'},
	yellow  = "\033[43m " //[]byte{keyEscape, '[', '3', '3', 'm'},
	blue    = "\033[44m " //[]byte{keyEscape, '[', '3', '4', 'm'},
	magenta = "\033[45m " //[]byte{keyEscape, '[', '3', '5', 'm'},
	cyan    = "\033[46m " //[]byte{keyEscape, '[', '3', '6', 'm'},
	white   = "\033[47m " //[]byte{keyEscape, '[', '3', '7', 'm'},
)

var oShape = shape{
	name:      "O",
	block:     red,
	gridIndex: 0,
	grids: []shapeGrid{
		{
			{true, true},
			{true, true},
		},
	},
	top:  0,
	left: 4,
}

var lShape = shape{
	name:      "L",
	block:     green,
	gridIndex: 0,
	grids: []shapeGrid{
		{
			{false, true, false},
			{false, true, false},
			{false, true, true},
		},
		{
			{false, false, false},
			{true, true, true},
			{true, false, false},
		},
		{
			{true, true, false},
			{false, true, false},
			{false, true, false},
		},
		{
			{false, false, true},
			{true, true, true},
			{false, false, false},
		},
	},
	top:  0,
	left: 4,
}

var jShape = shape{
	name:      "J",
	block:     yellow,
	gridIndex: 0,
	grids: []shapeGrid{
		{
			{false, true, false},
			{false, true, false},
			{true, true, false},
		},
		{
			{true, false, false},
			{true, true, true},
			{false, false, false},
		},
		{
			{false, true, true},
			{false, true, false},
			{false, true, false},
		},
		{
			{false, false, false},
			{true, true, true},
			{false, false, true},
		},
	},
	top:  0,
	left: 4,
}

var iShape = shape{
	name:      "I",
	block:     blue,
	gridIndex: 0,
	grids: []shapeGrid{
		{
			{false, false, false, false},
			{true, true, true, true},
			{false, false, false, false},
		},
		{
			{false, true, false},
			{false, true, false},
			{false, true, false},
			{false, true, false},
		},
	},
	top:  0,
	left: 4,
}

var sShape = shape{
	name:      "S",
	block:     cyan,
	gridIndex: 0,
	grids: []shapeGrid{
		{
			{false, true, false},
			{false, true, true},
			{false, false, true},
		},
		{
			{false, false, false},
			{false, true, true},
			{true, true, false},
		},
		{
			{true, false, false},
			{true, true, false},
			{false, true, false},
		},
		{
			{false, true, true},
			{true, true, false},
			{false, false, false},
		},
	},
	top:  0,
	left: 4,
}

var zShape = shape{
	name:      "S",
	block:     magenta,
	gridIndex: 0,
	grids: []shapeGrid{
		{
			{false, true, false},
			{true, true, false},
			{true, false, false},
		},
		{
			{true, true, false},
			{false, true, true},
			{false, false, false},
		},
		{
			{false, false, true},
			{false, true, true},
			{false, true, false},
		},
		{
			{false, false, false},
			{true, true, false},
			{false, true, true},
		},
	},
	top:  0,
	left: 4,
}


var tShape = shape{
	name:      "T",
	block:     black,
	gridIndex: 0,
	grids: []shapeGrid{
		{
			{false, true, false},
			{true, true, true},
			{false, false, false},
		},
		{
			{false, true, false},
			{false, true, true},
			{false, true, false},
		},
		{
			{false, false, false},
			{true, true, true},
			{false, true, false},
		},
		{
			{false, true, false},
			{true, true, false},
			{false, true, false},
		},
	},
	top:  0,
	left: 4,
}

var shapeLib = []shape{
	sShape,
	tShape,
	zShape,
	iShape,
	jShape,
	lShape,
	oShape,
}


type shapeGrid [][]bool

type shape struct {
	name      string
	block     string
	grids     []shapeGrid
	gridIndex int
	top       int
	left      int
}