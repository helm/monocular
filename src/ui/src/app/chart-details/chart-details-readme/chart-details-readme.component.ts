import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { Chart } from '../../shared/models/chart';
import { ChartsService } from '../../shared/services/charts.service';
import { ChartVersion } from '../../shared/models/chart-version';

@Component({
  selector: 'app-chart-details-readme',
  templateUrl: './chart-details-readme.component.html',
  styleUrls: ['./chart-details-readme.component.scss']
})
export class ChartDetailsReadmeComponent implements OnChanges {
  @Input() chart: Chart;
  @Input() currentVersion: ChartVersion;

  loading: boolean = true;
  readmeContent: string;
  markdown = require('marked');

  constructor(
    private chartsService: ChartsService,
  ) { }

  // Detect if input changed
  ngOnChanges(changes: SimpleChanges) {
    this.getReadme();
  }

  // TODO. This should not require loading the specific version and then the readme
  getReadme(): void {
    if (!this.currentVersion) return;
    this.chartsService.getChartReadme(this.currentVersion)
      .subscribe(resp => {
        this.loading = false;
        this.readmeContent = this.markdown(resp.text());
      });
  }
}
