import { Component, OnInit, Input } from '@angular/core';
import { Chart } from '../../shared/models/chart';
import { ChartsService } from '../../shared/services/charts.service';
import { ActivatedRoute, Params } from '@angular/router';

@Component({
  selector: 'app-chart-details-readme',
  templateUrl: './chart-details-readme.component.html',
  styleUrls: ['./chart-details-readme.component.scss']
})
export class ChartDetailsReadmeComponent implements OnInit {
  @Input() chart: Chart
  readmeContent: String
  markdown = require( "markdown" ).markdown;

  constructor(
    private route: ActivatedRoute,
    private chartsService: ChartsService,
  ) { }

  ngOnInit() {
    this.getReadme()
  }

  getReadme(): void {
    this.route.params.forEach((params: Params) => {
      let repo = params['repo'];
      let chartName = params['chartName']
      let latestVersion = this.chart.relationships.latestChartVersion.data.version
      this.chartsService.getChartReadme(repo, chartName, latestVersion)
        .subscribe(chart => {
          this.readmeContent = this.markdown.toHTML(chart.attributes.content)
        })
      })
  }
}
