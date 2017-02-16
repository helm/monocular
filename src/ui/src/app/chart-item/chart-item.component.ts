import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Chart } from '../shared/models/chart';
import { ConfigService } from '../shared/services/config.service';

@Component({
  selector: 'app-chart-item',
  templateUrl: './chart-item.component.html',
  styleUrls: ['./chart-item.component.scss'],
  inputs: ['chart', 'showVersion', 'showDescription']
})
export class ChartItemComponent implements OnInit {
  // Chart to represent
  public chart: Chart;
  // Show version form by default
  public showVersion: boolean = true;
  // Truncate the description
  public showDescription: boolean = true;

  constructor(
    private router: Router,
    private config: ConfigService
  ) {}

  ngOnInit() {
  }

	goToDetail(chart: Chart): void {
    let link = ['/charts', chart.attributes.repo.name, chart.attributes.name];
    this.router.navigate(link);
  }

  goToRepo(repo: string): void {
    let link = ['/charts', repo];
    this.router.navigate(link);
  }

  /**
   * Display the icon of the application if it's provided. In the other case,
   * It will return an string for a placeholder.
   *
   * @return {string} The URL of the icon or a placeholder
   */
  getIconUrl(chart: Chart): string {
    let icons = chart.relationships.latestChartVersion.data.icons;
    if (icons !== undefined && icons.length > 0) {
      return this.config.backendHostname + icons.find(icon => icon.name === '160x160-fit').path;
    } else {
      return '/assets/images/placeholder.png';
    }
  }
}
