import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Chart } from '../chart';

@Component({
  selector: 'app-chart-item',
  templateUrl: './chart-item.component.html',
  styleUrls: ['./chart-item.component.scss'],
  inputs: ['chart', 'showVersion', 'truncateDescription', 'fixHeight']
})
export class ChartItemComponent implements OnInit {
  // Chart to represent
  public chart: Chart;
  // Show version form by default
  public showVersion: boolean = true;
  // Truncate the description
  public truncateDescription: boolean = true;
  // Fix the height of the elements
  public fixHeight: boolean = false;

  constructor(private router: Router) { }

  ngOnInit() {
  }

	goToDetail(chart: Chart): void {
    let link = ['/charts', chart.attributes.repo, chart.attributes.name];
    this.router.navigate(link);
  }
}
