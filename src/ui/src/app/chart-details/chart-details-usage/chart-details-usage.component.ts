import { Component, OnInit, Input } from '@angular/core';
import { Chart } from '../../shared/models/chart';

@Component({
  selector: 'app-chart-details-usage',
  templateUrl: './chart-details-usage.component.html',
  styleUrls: ['./chart-details-usage.component.scss']
})
export class ChartDetailsUsageComponent implements OnInit {
  @Input() chart: Chart
  @Input() currentVersion: string

  constructor() { }

  ngOnInit() {}

  // Deletes /stable prefix not needed for stable repos
  get cmdChartId(): string {
    return this.chart.id.replace("stable/", "")
  }

  get showRepoInstructions(): boolean {
    return this.chart.attributes.repo != 'stable'
  }

  get repoAddInstructions(): string {
    return `helm repo add incubator ${this.chart.attributes.repoURL}`
  }
}
