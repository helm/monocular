import { Component, OnInit, Input } from '@angular/core';
import { Router } from '@angular/router';
import { Chart } from '../../chart';

@Component({
  selector: 'app-chart-item',
  templateUrl: './chart-item.component.html',
  styleUrls: ['./chart-item.component.scss']
})
export class ChartItemComponent implements OnInit {

  @Input() chart: Chart;
  constructor(private router: Router) { }

  ngOnInit() {
  }

	goToDetail(chart: Chart): void {
    let link = ['/charts', chart.attributes.repo, chart.attributes.name];
    this.router.navigate(link);
  }
}
