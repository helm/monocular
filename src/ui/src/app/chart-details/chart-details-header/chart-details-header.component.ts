import { Component, OnInit, Input } from '@angular/core';
import { Chart } from '../../chart';

@Component({
  selector: 'app-chart-details-header',
  templateUrl: './chart-details-header.component.html',
  styleUrls: ['./chart-details-header.component.scss']
})
export class ChartDetailsHeaderComponent implements OnInit {
  @Input() chart: Chart
  // TODO, remove
  latestVersion: String = '1.2.3-mocked'

  constructor() { }

  ngOnInit() {
  }

}
