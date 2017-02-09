import { Component, Input, OnChanges } from '@angular/core';
import { Chart } from '../../shared/models/chart';
import { ChartsService } from '../../shared/services/charts.service';

@Component({
  selector: 'app-chart-details-readme',
  templateUrl: './chart-details-readme.component.html',
  styleUrls: ['./chart-details-readme.component.scss']
})
export class ChartDetailsReadmeComponent implements OnChanges {
  @Input() chart: Chart
  @Input() currentVersion: string
  readmeContent: string
  markdown = require('marked')

  constructor(
    private chartsService: ChartsService,
  ) { }

  // Detect if input changed
  ngOnChanges() {
    this.getReadme()
  }

  // TODO. This should not require loading the specific version and then the readme
  getReadme(): void {
    this.chartsService.getVersion(this.chart.attributes.repo, this.chart.attributes.name, this.currentVersion)
      .subscribe(chartVersion => {
        this.chartsService.getChartReadme(chartVersion)
          .subscribe(resp => {
            this.readmeContent = this.markdown(resp.text())
          })
      })
  }
}
