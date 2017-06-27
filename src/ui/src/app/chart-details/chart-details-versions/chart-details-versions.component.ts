import { Component, OnInit, Input } from '@angular/core';
import { ChartVersion } from '../../shared/models/chart-version';
import { ChartAttributes } from '../../shared/models/chart';

@Component({
  selector: 'app-chart-details-versions',
  templateUrl: './chart-details-versions.component.html',
  styleUrls: ['./chart-details-versions.component.scss']
})
export class ChartDetailsVersionsComponent implements OnInit {
  @Input() versions: ChartVersion[]
  @Input() currentVersion: ChartVersion
  constructor() { }

  ngOnInit() { }

  goToVersionUrl(version: ChartVersion): string {
    let chart: ChartAttributes = version.relationships.chart.data
    return `/charts/${chart.repo.name}/${chart.name}/${version.attributes.version}`;
  }

  isSelected(version: ChartVersion): boolean {
    return version.attributes.version == this.currentVersion.attributes.version;
  }
}
