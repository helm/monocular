import { Component, OnInit } from '@angular/core';
import { ChartsService } from '../charts.service';
import { Chart } from '../chart';

@Component({
  selector: 'app-charts',
  templateUrl: './charts.component.html',
  styleUrls: ['./charts.component.scss']
})
export class ChartsComponent implements OnInit {
  charts: Chart[]
  constructor(private chartsService: ChartsService) { }

  ngOnInit() {
		this.loadCharts();
  }

  loadCharts(): void {
		this.chartsService.getCharts().subscribe(charts => this.charts = charts);
  }
}
