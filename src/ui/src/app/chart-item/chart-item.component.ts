import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Chart } from '../shared/models/chart';
import { CONFIG } from '../config';

@Component({
  selector: 'app-chart-item',
  templateUrl: './chart-item.component.html',
  styleUrls: ['./chart-item.component.scss'],
  inputs: ['chart', 'showVersion', 'truncateDescription', 'fixHeight']
})
export class ChartItemComponent implements OnInit {
  // Chart to represent
  public chart: Chart;
  // Show version form by default
  public showVersion: boolean = true;
  // Truncate the description
  public truncateDescription: boolean = true;
  // Fix the height of the elements
  public fixHeight: boolean = false;

  constructor(private router: Router) { }

  ngOnInit() {
  }

	goToDetail(chart: Chart): void {
    let link = ['/charts', chart.attributes.repo, chart.attributes.name];
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
      return CONFIG.backendHostname + icons.find(icon => icon.name === '160x160-fit').path;
    } else {
      return '/assets/images/placeholder.png';
    }
  }
}
