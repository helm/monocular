import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { ChartsService } from '../charts.service';
import { Chart } from '../chart';

@Component({
  selector: 'app-chart-details',
  templateUrl: './chart-details.component.html',
  styleUrls: ['./chart-details.component.scss']
})
export class ChartDetailsComponent implements OnInit {
  /* This resource will be different, probably ChartVersion */
  chart: Chart

  constructor(
    private route: ActivatedRoute,
    private chartsService: ChartsService
  ) { }

  ngOnInit() {
    /*TODO: Move this to resolver */
    this.route.params.forEach((params: Params) => {
      let repo = params['repo'];
      let chartName = params['chartName']
      this.chartsService.getChart(repo, chartName)
        .then(chart => this.chart = chart)
      })
    }
}
