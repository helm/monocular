import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { ChartsService } from '../shared/services/charts.service';
import { Chart } from '../shared/models/chart';
import { SeoService } from '../shared/services/seo.service';
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
    private seo: SeoService
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
   * Update the metatags with the name and the description of the application.
   */
  updateMetaTags(): void {
    if (this.titleVersion.length > 0) {
      this.seo.setMetaTags('detailsWithVersion', {
        name: this.chart.attributes.name,
        description: this.chart.attributes.description,
        version: this.titleVersion
      });
    } else {
      this.seo.setMetaTags('details', {
        name: this.chart.attributes.name,
        description: this.chart.attributes.description
      });
    }
  }
}
