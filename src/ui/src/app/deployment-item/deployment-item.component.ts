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
  }

  /**
   * Display the icon of the application if it's provided. In the other case,
   * It will return a string for a placeholder.
   *
   * @return {string} The URL of the icon or a placeholder
   */
  getIconUrl(): string {
    return this.deployment.attributes.chartIcon ? this.deployment.attributes.chartIcon : '/assets/images/placeholder.png';
  }
}
