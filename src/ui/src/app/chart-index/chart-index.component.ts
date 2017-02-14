import { Component, OnInit } from '@angular/core';
import { ChartsService } from '../shared/services/charts.service';
import { Chart } from '../shared/models/chart';

@Component({
  selector: 'app-chart-index',
  templateUrl: './chart-index.component.html',
  styleUrls: ['./chart-index.component.scss']
})
export class ChartIndexComponent implements OnInit {
	charts: Chart[]
  loading: boolean = true;

  constructor(private chartsService: ChartsService) { }

  ngOnInit() {
		this.loadCharts();
  }

  loadCharts(): void {
		this.chartsService.getCharts().subscribe(charts => {
      this.loading = false;
      this.charts = charts;
    });
  }
}
