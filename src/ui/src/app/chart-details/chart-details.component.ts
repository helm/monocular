import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { ChartsService } from '../shared/services/charts.service';
import { Chart } from '../shared/models/chart';
import { ChartVersion } from '../shared/models/chart-version';
import { SeoService } from '../shared/services/seo.service';
import { ConfigService } from '../shared/services/config.service';
import ColorThief from 'color-thief-browser';

@Component({
  selector: 'app-chart-details',
  templateUrl: './chart-details.component.html',
  styleUrls: ['./chart-details.component.scss']
})
export class ChartDetailsComponent implements OnInit {
  /* This resource will be different, probably ChartVersion */
  chart: Chart;
  loading: boolean = true;
  currentVersion: ChartVersion;
  iconUrl: string;
  chartColor: string;
  titleVersion: string;

  constructor(
    private route: ActivatedRoute,
    private chartsService: ChartsService,
    private config: ConfigService,
    private seo: SeoService
  ) {}

  ngOnInit() {
    this.route.params.forEach((params: Params) => {
      let repo = params['repo'];
      let chartName = params['chartName'];
      this.chartsService.getChart(repo, chartName).subscribe(chart => {
        this.loading = false;
        this.chart = chart;
        let version =
          params['version'] ||
          this.chart.relationships.latestChartVersion.data.version;
        this.chartsService
          .getVersion(repo, chartName, version)
          .subscribe(chartVersion => {
            this.currentVersion = chartVersion;
          });
        this.titleVersion = params['version'] || '';
        this.updateMetaTags();
        this.iconUrl = this.getIconUrl();
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

  goToRepoUrl(): string {
    return `/charts/${this.chart.attributes.repo.name}`;
  }

  getIconUrl(): string {
    let icons = this.chart.relationships.latestChartVersion.data.icons;
    if (icons !== undefined && icons.length > 0) {
      const icon =
        this.config.backendHostname +
        icons.find(icon => icon.name === '160x160-fit').path;
      if (!this.chartColor) {
        const imgObj = new Image();
        imgObj.crossOrigin = 'Anonymous';
        imgObj.src = icon;
        imgObj.addEventListener('load', e => {
          const ct = new ColorThief();
          const palette = ct.getPalette(imgObj, 2);
          if (palette.length > 0) {
            const rgb = palette[0];
            this.chartColor = `rgba(${rgb[0]}, ${rgb[1]}, ${rgb[2]}, 1)`;
          }
        });
      }

      return icon;
    } else {
      return '/assets/images/placeholder.png';
    }
  }
}
