import { Component, OnInit, Input } from '@angular/core';
import { ChartsService } from '../../shared/services/charts.service';
import { Chart } from '../../shared/models/chart';
import { ChartVersion } from '../../shared/models/chart-version';
import { Router } from '@angular/router';

@Component({
  selector: 'app-chart-details-info',
  templateUrl: './chart-details-info.component.html',
  styleUrls: ['./chart-details-info.component.scss']
})
export class ChartDetailsInfoComponent implements OnInit {
  @Input() chart: Chart
  @Input() currentVersion: string
  versions: ChartVersion[]
  constructor(
    private chartsService: ChartsService,
    private router: Router
  ) { }

  ngOnInit() {
    this.loadVersions(this.chart)
  }

  get sources() {
    return this.chart.attributes.sources || [];
  }

  loadVersions(chart: Chart): void {
    this.chartsService.getVersions(chart.attributes.repo, chart.attributes.name)
      .subscribe(versions => { this.versions = versions })
  }
}
