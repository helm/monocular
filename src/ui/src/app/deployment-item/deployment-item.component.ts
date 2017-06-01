import { Component, OnInit } from '@angular/core';
import { Deployment } from '../shared/models/deployment';
import { ConfigService } from '../shared/services/config.service';

@Component({
  selector: 'app-deployment-item',
  templateUrl: './deployment-item.component.html',
  styleUrls: ['./deployment-item.component.scss'],
  inputs: ['deployment']
})
export class DeploymentItemComponent implements OnInit {
  // Chart to represent
  public deployment: Deployment;

  constructor(
    private config: ConfigService
  ) {}

  ngOnInit() {
    console.log(this.deployment)
  }

  /**
   * Display the icon of the application if it's provided. In the other case,
   * It will return a string for a placeholder.
   *
   * @return {string} The URL of the icon or a placeholder
   */
  getIconUrl(): string {
    // let icons = this.deployment.chart.relationships.latestChartVersion.data.icons;
    // if (icons !== undefined && icons.length > 0) {
    //   return this.config.backendHostname + icons.find(icon => icon.name === '160x160-fit').path;
    // } else {
      return '/assets/images/placeholder.png';
    // }
  }
}
