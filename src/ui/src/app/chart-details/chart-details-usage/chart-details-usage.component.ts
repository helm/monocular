import { Component, OnInit, Input } from '@angular/core';
import { Chart } from '../../shared/models/chart';

@Component({
  selector: 'app-chart-details-usage',
  templateUrl: './chart-details-usage.component.html',
  styleUrls: ['./chart-details-usage.component.scss']
})
export class ChartDetailsUsageComponent implements OnInit {
  @Input() chart: Chart
  installCommand: String

  constructor() { }

  ngOnInit() {
    let latestVersion: String = this.chart.relationships.latestChartVersion.data.version
    this.installCommand = `helm install ${ this.chart.id } --version ${ latestVersion }`
  }

}
