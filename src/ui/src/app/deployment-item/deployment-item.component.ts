import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { Deployment } from '../shared/models/deployment';
import { ConfigService } from '../shared/services/config.service';
import ColorThief from 'color-thief-browser'

@Component({
  selector: 'app-deployment-item',
  templateUrl: './deployment-item.component.html',
  styleUrls: ['./deployment-item.component.scss'],
  inputs: ['deployment']
})
export class DeploymentItemComponent implements OnInit {
  backgroundColor: string;

  // Chart to represent
  public deployment: Deployment;

  constructor(
    private config: ConfigService
  ) {}

  ngOnInit() { }

  /**
   * Display the icon of the application if it's provided. In the other case,
   * It will return a string for a placeholder.
   *
   * @return {string} The URL of the icon or a placeholder
   */
  getIconUrl(): string {
    if (this.deployment.attributes.chartIcon && !this.backgroundColor) {
      var imgObj = new Image();
      imgObj.crossOrigin = 'Anonymous';
      imgObj.src = this.deployment.attributes.chartIcon;
      imgObj.addEventListener('load', (e) => {
        const ct = new ColorThief();
        const palette = ct.getPalette(imgObj, 2);
        if (palette.length > 0) {
          const rgb = palette[0];
          this.backgroundColor = `rgba(${rgb[0]}, ${rgb[1]}, ${rgb[2]}, 0.1)`
        }
      })
    }
    return this.deployment.attributes.chartIcon ? this.deployment.attributes.chartIcon : '/assets/images/placeholder.png';
  }
}
