import { Component, OnInit } from '@angular/core';
import { ChartsService } from '../shared/services/charts.service';
import { Chart } from '../shared/models/chart';
import { ConfigService } from '../shared/services/config.service';
import { MetaService } from 'ng2-meta';

@Component({
  selector: 'app-chart-index',
  templateUrl: './chart-index.component.html',
  styleUrls: ['./chart-index.component.scss']
})
export class ChartIndexComponent implements OnInit {
	charts: Chart[]
  loading: boolean = true;
  totalChartsNumber: number

  constructor(
    private chartsService: ChartsService,
    private config: ConfigService,
    private metaService: MetaService
  ) {}

  ngOnInit() {
		this.loadCharts();
    this.updateMetaTags();
  }

  loadCharts(): void {
		this.chartsService.getCharts().subscribe(charts => {
      this.loading = false;
      this.charts = charts;
      this.totalChartsNumber = charts.length;
    });
  }

  updateMetaTags(): void {
    let title: string = this.config.appName;
    this.metaService.setTitle(title, "");
    this.metaService.setTag('og:title', title);
  }
}
