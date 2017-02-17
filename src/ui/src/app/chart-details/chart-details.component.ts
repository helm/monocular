import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { ChartsService } from '../shared/services/charts.service';
import { Chart } from '../shared/models/chart';
import { MetaService } from 'ng2-meta';
import { ConfigService } from '../shared/services/config.service';

@Component({
  selector: 'app-chart-details',
  templateUrl: './chart-details.component.html',
  styleUrls: ['./chart-details.component.scss']
})
export class ChartDetailsComponent implements OnInit {
  /* This resource will be different, probably ChartVersion */
  chart: Chart;
  loading: boolean = true;
  currentVersion: string;
  titleVersion: string;

  constructor(
    private route: ActivatedRoute,
    private chartsService: ChartsService,
    private config: ConfigService,
    private metaService: MetaService
  ) { }

  ngOnInit() {
    this.route.params.forEach((params: Params) => {
      let repo = params['repo'];
      let chartName = params['chartName']
      this.chartsService.getChart(repo, chartName)
        .subscribe(chart => {
          this.loading = false;
          this.chart = chart;
          this.currentVersion = params['version'] || this.chart.relationships.latestChartVersion.data.version;
          this.titleVersion = params['version'] || '';
          this.updateMetaTags();
        });
    });
  }

  /**
   * Content for the title version tag
   *
   * @return {string} Title to display in the site
   */
  contentTitleVersion(): string {
    if (this.titleVersion.length > 0) {
      return `${this.chart.attributes.name} ${this.titleVersion}`;
    } else {
      return this.chart.attributes.name;
    }
  }

  /**
   * Update the metatags with the name and the description of the application.
   */
  updateMetaTags(): void {
    let title: string = this.contentTitleVersion();
    this.metaService.setTitle(title, ` | ${this.config.appName}`);
    this.metaService.setTag('description', this.chart.attributes.description);
    this.metaService.setTag('og:title', title);
    this.metaService.setTag('og:description', this.chart.attributes.description);
  }
}
