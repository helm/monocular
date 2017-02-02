import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { ChartsService } from '../shared/services/charts.service';
import { Chart } from '../shared/models/chart';
import { MetaService } from 'ng2-meta';

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
    private chartsService: ChartsService,
    private metaService: MetaService
  ) { }

  ngOnInit() {
    this.route.params.forEach((params: Params) => {
      let repo = params['repo'];
      let chartName = params['chartName']
      this.chartsService.getChart(repo, chartName)
        .subscribe(chart => {
          this.chart = chart
          this.currentVersion = params['version'] || this.chart.relationships.latestChartVersion.data.version
          this.updateMetaTags(chart);
        });
    });
  }

  /**
   * Update the metatags with the name and the description of the application.
   *
   * @param {Chart} chart The chart to get the name and the description.
   */
  updateMetaTags(chart: Chart): void {
    this.metaService.setTitle(chart.attributes.name);
    this.metaService.setTag('description', chart.attributes.description);
    this.metaService.setTag('og:title', chart.attributes.name);
    this.metaService.setTag('og:description', chart.attributes.description);
  }
}
