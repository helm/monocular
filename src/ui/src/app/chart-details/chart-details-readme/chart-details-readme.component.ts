import { Component, OnInit, Input } from '@angular/core';
import { Chart } from '../../chart';
import { ChartsService } from '../../charts.service';

@Component({
  selector: 'app-chart-details-readme',
  templateUrl: './chart-details-readme.component.html',
  styleUrls: ['./chart-details-readme.component.scss']
})
export class ChartDetailsReadmeComponent implements OnInit {
  @Input() chart: Chart
  readmeContent: String
  markdown = require( "markdown" ).markdown;

  constructor(private chartsService: ChartsService) { }

  ngOnInit() {
    this.getMarkDownMock()
  }

  getMarkDownMock(): void {
    this.chartsService.getMockedReadme().forEach((response) => {
      this.readmeContent = this.markdown.toHTML(response.text())
    })
  }
}
