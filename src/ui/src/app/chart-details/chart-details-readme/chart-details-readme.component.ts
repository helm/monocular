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
  @Input() currentVersion: String
  readmeContent: String
  markdown = require('marked')

  constructor(
    private chartsService: ChartsService,
  ) { }

  // Detect if input changed
  ngOnChanges() {
    this.getReadme()
  }

  getReadme(): void {
    this.chartsService.getChartReadme(this.chart.attributes.repo, this.chart.attributes.name, this.currentVersion)
      .subscribe(chart => {
        this.readmeContent = this.markdown(chart.attributes.content)
      })
  }
}
