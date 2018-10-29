package goChartjs

import(
    "fmt"
)

func (c *Color)String()string{
        return fmt.Sprintf("rgba(%d, %d, %d, %f)", c.R, c.G, c.B, c.A);
}
