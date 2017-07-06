import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { Deployment } from '../shared/models/deployment';
import { ConfigService } from '../shared/services/config.service';
import RGBaster from '../../assets/js/RGBaster';

@Component({
  selector: 'app-deployment-item',
  templateUrl: './deployment-item.component.html',
  styleUrls: ['./deployment-item.component.scss'],
  inputs: ['deployment']
})
export class DeploymentItemComponent implements OnInit {
  private themeColor: string;
  private iconUrl: string;

  // Chart to represent
  public deployment: Deployment;


  ngOnInit() {
    this.iconUrl = this.getIconUrl();
  }

  getIconUrl(): string {
    if (this.deployment.attributes.chartIcon) {
      RGBaster.colors(this.deployment.attributes.chartIcon, {
        success: payload => {
          this.themeColor = payload.best
            .replace('rgb', 'rgba')
            .replace(')', ', 0.1)');
        }
      });
      return this.deployment.attributes.chartIcon
    }
    return '/assets/images/placeholder.png';
  }
}
