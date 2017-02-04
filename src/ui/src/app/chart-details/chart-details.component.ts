import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { ChartsService } from '../shared/services/charts.service';
import { Chart } from '../shared/models/chart';

@Component({
  selector: 'app-chart-details',
  templateUrl: './chart-details.component.html',
  styleUrls: ['./chart-details.component.scss']
})
export class ChartDetailsComponent implements OnInit {
  /* This resource will be different, probably ChartVersion */
  chart: Chart
  currentVersion: String

  constructor(
    private route: ActivatedRoute,
    private chartsService: ChartsService
  ) { }

  ngOnInit() {
    this.route.params.forEach((params: Params) => {
      let repo = params['repo'];
      let chartName = params['chartName']
      this.chartsService.getChart(repo, chartName)
        .subscribe(chart => {
          this.chart = chart
          this.currentVersion = params['version'] || this.chart.relationships.latestChartVersion.data.version
        })
      })
  }
}
