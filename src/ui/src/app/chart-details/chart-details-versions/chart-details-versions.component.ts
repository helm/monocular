import { Component, OnInit, Input } from '@angular/core';
import { ChartVersion } from '../../shared/models/chart-version';
import { ChartAttributes } from '../../shared/models/chart';
import { Router } from '@angular/router';

@Component({
  selector: 'app-chart-details-versions',
  templateUrl: './chart-details-versions.component.html',
  styleUrls: ['./chart-details-versions.component.scss']
})
export class ChartDetailsVersionsComponent implements OnInit {
  @Input() versions: ChartVersion[]
  @Input() currentVersion: string
  constructor(
    private router: Router,
  ) { }

  ngOnInit() { }

	goToVersion(version: ChartVersion): void {
    var chart: ChartAttributes = version.relationships.chart.data
    let link = ['/charts', chart.repo.name, chart.name, version.attributes.version];
    this.router.navigate(link);
  }

  isSelected(version: ChartVersion): boolean {
    return version.attributes.version == this.currentVersion
  }
}
