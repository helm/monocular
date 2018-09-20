import { Component, OnInit, Input } from '@angular/core';
import { ChartsService } from '../../shared/services/charts.service';
import { Chart } from '../../shared/models/chart';
import { Maintainer } from '../../shared/models/maintainer';
import { ChartVersion } from '../../shared/models/chart-version';
import * as urljoin from 'url-join';

@Component({
  selector: 'app-chart-details-info',
  templateUrl: './chart-details-info.component.html',
  styleUrls: ['./chart-details-info.component.scss']
})
export class ChartDetailsInfoComponent implements OnInit {
  @Input() chart: Chart;
  @Input() currentVersion: ChartVersion;
  versions: ChartVersion[];
  constructor(private chartsService: ChartsService) {}

  ngOnInit() {
    this.loadVersions(this.chart);
  }

  get sources() {
    return this.chart.attributes.sources || [];
  }

  get maintainers(): Maintainer[] {
    return this.chart.attributes.maintainers || [];
  }

  loadVersions(chart: Chart): void {
    this.chartsService
      .getVersions(chart.attributes.repo.name, chart.attributes.name)
      .subscribe(versions => {
        this.versions = versions;
      });
  }

  maintainerUrl(maintainer: Maintainer): string {
    // Use GitHub URL with maintainer name if this is an upstream Helm repo from
    // github.com/helm/charts (i.e. stable or incubator)
    if (this.isUpstreamHelmRepo(this.chart.attributes.repo.url)) {
      return `https://github.com/${maintainer.name}`;
    } else {
      return `mailto:${maintainer.email}`;
    }
  }

  private isUpstreamHelmRepo(repoURL: string): boolean {
    return (
      repoURL === "https://kubernetes-charts.storage.googleapis.com" ||
      repoURL === "https://kubernetes-charts-incubator.storage.googleapis.com"
    );
  }
}
