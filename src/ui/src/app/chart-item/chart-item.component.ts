import { Component, OnInit } from '@angular/core';
import { Chart } from '../shared/models/chart';
import { ChartsService } from '../shared/services/charts.service';

@Component({
  selector: 'app-chart-item',
  templateUrl: './chart-item.component.html',
  styleUrls: ['./chart-item.component.scss'],
  inputs: ['chart', 'showVersion', 'showDescription']
})
export class ChartItemComponent implements OnInit {
  public iconUrl: string;
  // Chart to represent
  public chart: Chart;
  // Show version form by default
  public showVersion: boolean = true;
  // Truncate the description
  public showDescription: boolean = true;

  constructor(private chartsService: ChartsService) {}

  ngOnInit() {
    this.iconUrl = this.chartsService.getChartIconURL(this.chart);
  }

  goToDetailUrl(): string {
    return `/charts/${this.chart.attributes.repo.name}/${this.chart.attributes
      .name}`;
  }

  goToRepoUrl(): string {
    return `/charts/${this.chart.attributes.repo.name}`;
  }
}
