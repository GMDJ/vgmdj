package chars

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConversion(t *testing.T) {
	ast := assert.New(t)

	c1 := NewConversion(10, "32")
	ast.Equal(c1.BaseValue(), float64(32))
	ast.Equal(c1.ZoomOut().ToString(1), "320")
	ast.Equal(c1.ZoomOut().ToInt(), 320)
	ast.Equal(c1.ZoomIn().ToString(c1.multiples), "3.2")
	ast.Equal(c1.ZoomIn().ToString(100), "3.20")
	ast.Equal(c1.ZoomIn().ToInt(), 3)
	ast.Equal(c1.ZoomIn(320).ToFloat64(), 32.00)

	c2 := NewConversion(100, "32")
	ast.Equal(c2.BaseValue(), float64(32))
	ast.Equal(c2.ZoomOut().ToString(1), "3200")
	ast.Equal(c2.ZoomOut().ToInt(), 3200)
	ast.Equal(c2.ZoomIn().ToString(c2.multiples), "0.32")
	ast.Equal(c2.ZoomIn().ToString(10), "0.3")
	ast.Equal(c2.ZoomIn().ToInt(), 0)
	ast.Equal(c2.ZoomIn().ToFloat64(), 0.32)

}