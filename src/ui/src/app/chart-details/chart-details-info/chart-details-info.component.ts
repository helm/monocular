import { Component, OnInit, Input } from '@angular/core';
import { Chart } from '../../shared/models/chart';

@Component({
  selector: 'app-chart-details-info',
  templateUrl: './chart-details-info.component.html',
  styleUrls: ['./chart-details-info.component.scss']
})
export class ChartDetailsInfoComponent implements OnInit {
  @Input() chart: Chart
  constructor() { }

  ngOnInit() {
  }

  get lastVersion() {
    return this.chart.relationships.latestChartVersion;
  }

  get sources() {
    return this.chart.attributes.sources || [];
  }

}
