/*
Copyright (c) 2023-2026 Microbus LLC and various contributors

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

package showcase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/microbus-io/bespa/website/shared"
	"github.com/microbus-io/bespa/website/storage"
)

func HandleCharts(w http.ResponseWriter, r *http.Request) {
	unitedStates, _, _ := storage.USQuery("", "", 0, -1)
	sort.Slice(unitedStates, func(i, j int) bool {
		return unitedStates[i].GDPPerCapita < unitedStates[j].GDPPerCapita
	})

	echartsBar := wf.Chart(`{
		title: {
			text: 'Monthly Average Rainfall',
			subtext: 'Source: WorldClimate.com'
		},
		tooltip: {
			trigger: 'axis',
			axisPointer: {type: 'shadow'}
		},
		legend: {top: 'bottom'},
		grid: {top: 80, bottom: 60, left: 50, right: 20},
		xAxis: {
			type: 'category',
			data: ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec']
		},
		yAxis: {
			type: 'value',
			name: 'Rainfall (mm)'
		},
		series: [
			{name: 'Tokyo',    type: 'bar', data: [49.9, 71.5, 106.4, 129.2, 144.0, 176.0, 135.6, 148.5, 216.4, 194.1, 95.6, 54.4]},
			{name: 'New York', type: 'bar', data: [83.6, 78.8, 98.5, 93.4, 106.0, 84.5, 105.0, 104.3, 91.2, 83.5, 106.6, 92.3]},
			{name: 'London',   type: 'bar', data: [48.9, 38.8, 39.3, 41.4, 47.0, 48.3, 59.0, 59.6, 52.4, 65.2, 59.3, 51.2]},
			{name: 'Berlin',   type: 'bar', data: [42.4, 33.2, 34.5, 39.7, 52.6, 75.5, 57.4, 60.4, 47.6, 39.1, 46.8, 51.1]}
		]
	}`, nil)

	echartsPie := wf.Chart(`{
		title: {
			text: 'Browser market shares. January, 2022',
			subtext: 'Source: statcounter.com',
			left: 'left'
		},
		tooltip: {trigger: 'item', formatter: '{b}: {c}% ({d}%)'},
		legend: {orient: 'vertical', right: 0, top: 'middle'},
		series: [{
			name: 'Browsers',
			type: 'pie',
			radius: ['35%', '65%'],
			center: ['40%', '55%'],
			avoidLabelOverlap: true,
			label: {formatter: '{b}: {c}%'},
			data: [
				{name: 'Chrome',  value: 61.04},
				{name: 'Safari',  value: 9.47},
				{name: 'Edge',    value: 9.32},
				{name: 'Firefox', value: 8.15},
				{name: 'Other',   value: 11.02}
			]
		}]
	}`, nil)

	// Column + map demos driven by the US-states GDP/Capita dataset.
	abbrevs := make([]string, 0, len(unitedStates))
	gdps := make([]float64, 0, len(unitedStates))
	type mapPoint struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}
	mapData := make([]mapPoint, 0, len(unitedStates))
	var minGDP, maxGDP float64
	for i, us := range unitedStates {
		abbrevs = append(abbrevs, us.Abbrev)
		gdps = append(gdps, us.GDPPerCapita)
		mapData = append(mapData, mapPoint{Name: us.Name, Value: us.GDPPerCapita})
		if i == 0 || us.GDPPerCapita < minGDP {
			minGDP = us.GDPPerCapita
		}
		if us.GDPPerCapita > maxGDP {
			maxGDP = us.GDPPerCapita
		}
	}
	abbrevsJSON, _ := json.Marshal(abbrevs)
	gdpsJSON, _ := json.Marshal(gdps)
	mapDataJSON, _ := json.Marshal(mapData)

	echartsColumn := wf.Chart(fmt.Sprintf(`{
		title: { text: 'GDP/Capita' },
		grid: { top: 60, bottom: 80, left: 70, right: 20 },
		tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
		xAxis: {
			type: 'category',
			data: %s,
			axisLabel: { rotate: 75, interval: 0 }
		},
		yAxis: { type: 'value' },
		series: [{ type: 'bar', data: %s }]
	}`, abbrevsJSON, gdpsJSON), nil)

	echartsMap := wf.Chart(fmt.Sprintf(`{
		title: { text: 'GDP/Capita' },
		tooltip: {
			trigger: 'item',
			formatter: function(p) { return p.name + ': ' + (p.value ? p.value.toFixed(0) : 'N/A'); }
		},
		visualMap: {
			min: %f, max: %f,
			left: 'left', bottom: 20,
			text: ['High', 'Low'],
			calculable: true,
			inRange: {
				color: ['rgb(var(--md-sys-color-secondary-container))', 'rgb(var(--md-sys-color-primary))']
			}
		},
		series: [{
			type: 'map',
			map: 'US-MAINLAND',
			roam: false,
			emphasis: { label: { show: true } },
			data: %s
		}]
	}`, minGDP, maxGDP, mapDataJSON), nil).WithHeight("600px")

	page := wf.Page().Add(
		wf.AppBar("Chart widgets"),
		echartsBar,
		echartsPie,
		echartsColumn,
		echartsMap,
	)
	shared.Render(w, r, page)
}
