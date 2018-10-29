package goChartjs

// Bool is a convenience typedef for pointer to bool so that we can differentiate between unset
// and false.
type Bool *bool

var (
	t = true
	f = false
	// True is a convenience for pointer to true
	True = Bool(&t)

	// False is a convenience for pointer to false
	False = Bool(&f)
)

type Color struct {
	// Red
	R int
	// Green
	G int
	// Blue
	B int
	//Alpha
	A float64
}

type DisplayFormats struct {
	Quarter string `json:"quarter,omitempty"`
	Minute  string `json:"minte,omitempty"`
	Second  string `json:"second,omitempty"`
}

type Time struct {
	DisplayFormats DisplayFormats `json:"displayFormats,omitempty"`
	Unit           string         `json:"unit,omitempty"`
}

type Ticks struct {
	AxesType    string `json:"type,omitempty"`
	BeginAtZero bool   `json:"beginAtZero,omitempty"`
}

type Axes struct {
	Type  string `json:"type,omitempty"`
	Time  *Time  `json:"time,omitempty"`
	Ticks *Ticks `json:"ticks,omitempty"`
}

type Scales struct {
	YAxes []Axes `json:"yAxes,omitempty"`
	XAxes []Axes `json:"xAxes,omitempty"`
}

type Legend struct {
	Display Bool `json:"display,omitempty"`
}

type Title struct {
	Display Bool   `json:"display,omitempty"`
	Text    string `json:"text,omitempty"`
}

type Options struct {
	Scales              `json:"scales,omitempty"`
	Legend              *Legend `json:"legend,omitempty"`
	Responsive          Bool    `json:"responsive, omitempty"`
	Title               *Title  `json:"title,omitempty"`
	MaintainAspectRatio Bool    `json:"maintainAspectRatio,omitempty"`
}

type Point struct {
	X string  `json:"x, omitempty"`
	Y float64 `json:"y, omitempty"`
}

type Dataset struct {
	Label           string   `json:"label,omitempty"`
	Data            []Point  `json:"data,omitempty"`
	BackgroundColor string   `json:"backgroundColor,omitempty"`
	BorderColor     []string `json:"borderColor,omitempty"`
	BorderWidth     int      `json:"borderWidth,omitempty"`
	// A pointer is used since line tension can be 0.
	LineTension            *int   `json:"lineTension,omitempty"`
	CubicInterpolationMode string `json:"cubicInterpolationMode,omitempty"`
}

type ChartData struct {
	Labels   []string  `json:"labels,omitempty"`
	Datasets []Dataset `json:"datasets,omitempty"`
}

type Chart struct {
	Name      string                    `json:"-"` // omit
	ChartType string                    `json:"type,omitempty"`
	Data      ChartData                 `json:"data,omitempty"`
	Options   *Options                  `json:"options,omitempty"`
	GetData   func() ([]Dataset, error) `json:"-"`
}
