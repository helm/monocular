import { Component, OnInit } from '@angular/core';
import { Chart } from '../shared/models/chart';
import { ConfigService } from '../shared/services/config.service';

@Component({
  selector: 'app-chart-item',
  templateUrl: './chart-item.component.html',
  styleUrls: ['./chart-item.component.scss'],
  inputs: ['chart', 'showVersion', 'showDescription', 'fixHeight']
})
export class ChartItemComponent implements OnInit {
  // Chart to represent
  public chart: Chart;
  // Show version form by default
  public showVersion: boolean = true;
  // Truncate the description
  public showDescription: boolean = true;
  // Fix the height of the elements
  public fixHeight: boolean = false;

  constructor(
    private config: ConfigService
  ) {}

  ngOnInit() {
  }

	goToDetailUrl(): string {
    return `/charts/${this.chart.attributes.repo.name}/${this.chart.attributes.name}`;
  }

  goToRepoUrl(): string {
    return `/charts/${this.chart.attributes.repo.name}`;
  }

  /**
   * Display the icon of the application if it's provided. In the other case,
   * It will return an string for a placeholder.
   *
   * @return {string} The URL of the icon or a placeholder
   */
  getIconUrl(): string {
    let icons = this.chart.relationships.latestChartVersion.data.icons;
    if (icons !== undefined && icons.length > 0) {
      return this.config.backendHostname + icons.find(icon => icon.name === '160x160-fit').path;
    } else {
      return '/assets/images/placeholder.png';
    }
  }
}
