import { Component, OnInit, Input } from '@angular/core';
import { ChartsService } from '../../shared/services/charts.service';
import { Chart } from '../../shared/models/chart';
import { ChartVersion } from '../../shared/models/chart-version';

@Component({
  selector: 'app-chart-details-info',
  templateUrl: './chart-details-info.component.html',
  styleUrls: ['./chart-details-info.component.scss']
})
export class ChartDetailsInfoComponent implements OnInit {
  @Input() chart: Chart
  @Input() currentVersion: string
  versions: ChartVersion[]
  constructor(
    private chartsService: ChartsService,
  ) { }

  ngOnInit() {
    this.loadVersions(this.chart)
  }

  get sources() {
    return this.chart.attributes.sources || [];
  }

  get sourceUrl(): string {
    var chartSource = this.chart.attributes.repo.source;
    if (!chartSource) return

    // Used to handle possible trailing URLs
    var urljoin = require('url-join');
    return urljoin(chartSource, this.chart.attributes.name);
  }

  get sourceName(): string {
    var parser = document.createElement('a');
    parser.href = this.chart.attributes.repo.source;
    return parser.hostname;
  }

  loadVersions(chart: Chart): void {
    this.chartsService.getVersions(chart.attributes.repo.name, chart.attributes.name)
      .subscribe(versions => { this.versions = versions })
  }
}
