import { Component, OnInit, Input } from '@angular/core';
import { Chart } from '../../chart';

@Component({
  selector: 'app-chart-details-usage',
  templateUrl: './chart-details-usage.component.html',
  styleUrls: ['./chart-details-usage.component.scss']
})
export class ChartDetailsUsageComponent implements OnInit {
  @Input() chart: Chart
  installCommand: String
  // TODO, remove
  latestVersion: String = '1.2.3-mocked'

  constructor() { }

  ngOnInit() {
    this.installCommand = `helm install ${ this.chart.id }-${ this.latestVersion }.tgz`
  }

}
