import { Component, OnInit } from '@angular/core';
import { ChartsService } from '../charts.service';
import { Chart } from '../chart';

@Component({
  selector: 'app-chart-index',
  templateUrl: './chart-index.component.html',
  styleUrls: ['./chart-index.component.scss']
})
export class ChartIndexComponent implements OnInit {
	charts: Chart[]
  constructor(private chartsService: ChartsService) { }

  ngOnInit() {
		this.loadCharts();
  }

  loadCharts(): void {
		this.chartsService.getCharts().then(charts => this.charts = charts);
  }
}
