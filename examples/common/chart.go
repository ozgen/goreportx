package common

import (
	"bytes"
	"encoding/base64"
	"github.com/wcharczuk/go-chart/v2"
)

func GenerateChartBase64() string {
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: []float64{1, 2, 3, 4, 5},
				YValues: []float64{1, 2, 1, 3, 4},
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	_ = graph.Render(chart.PNG, buffer)

	encoded := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return "data:image/png;base64," + encoded
}
