import { Component, OnInit, Input } from '@angular/core';
import { Chart } from '../../shared/models/chart';

@Component({
  selector: 'app-chart-list-item',
  templateUrl: './chart-list-item.component.html',
  styleUrls: ['./chart-list-item.component.scss']
})
export class ChartListItemComponent implements OnInit {
  @Input() chart: Chart;
  constructor() { }

  ngOnInit() {
  }

}
