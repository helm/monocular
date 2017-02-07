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

  // TODO, remove hardcoded code once https://github.com/helm/monocular/issues/86 is implemented
  get repoAddInstructions(): string {
    if (this.chart.attributes.repo == "incubator") {
      return "helm repo add incubator https://kubernetes-charts-incubator.storage.googleapis.com/"
    }
  }
}
