/*
Copyright (c) 2019 The Helm Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"bytes"
	"image/color"
	"github.com/helm/monocular/cmd/chartsvc/models"

	"github.com/disintegration/imaging"
)

//ChartsList a list of charts used in unit tests
var ChartsList []*models.Chart

//IconBytes the bytes of a chart icon image used in unit tests
func IconBytes() []byte {
	var b bytes.Buffer
	img := imaging.New(1, 1, color.White)
	imaging.Encode(&b, img, imaging.PNG)
	return b.Bytes()
}

const TestChartReadme = "# Quickstart\n\n```bash\nhelm install my-repo/my-chart\n```"
const TestChartValues = "image:\n  registry: docker.io\n  repository: my-repo/my-chart\n  tag: 0.1.0"
const TestChartSchema = `{"properties": {"type": "object"}}`
