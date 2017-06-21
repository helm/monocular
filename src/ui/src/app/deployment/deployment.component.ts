import { Router, ActivatedRoute, Params } from '@angular/router';
import { Component, OnInit, ElementRef, ViewChild } from '@angular/core';
import { DeploymentsService } from '../shared/services/deployments.service';
import { Deployment } from '../shared/models/deployment';
import { Chart } from '../shared/models/chart';
import { ConfigService } from '../shared/services/config.service';
import { DomSanitizer } from '@angular/platform-browser';
import { MdIconRegistry } from '@angular/material';
import RGBaster from '../../assets/js/RGBaster';

@Component({
  selector: 'app-deployment',
  templateUrl: './deployment.component.html',
  styleUrls: ['./deployment.component.scss'],
  viewProviders: [MdIconRegistry]
})
export class DeploymentComponent implements OnInit {
  backgroundColor: string;

  deployment: Deployment;
  resources = [];
  loading: boolean = true;
  name: String = '';

  constructor(
    private deploymentsService: DeploymentsService,
    private router: Router,
    private route: ActivatedRoute,
    private config: ConfigService,
    private mdIconRegistry: MdIconRegistry,
    private sanitizer: DomSanitizer
  ) {}

  ngOnInit() {
    const icons = [
      'layers',
      'schedule',
      'web-asset',
      'info-outline',
      'arrow-back'
    ];

    icons.forEach(icon => {
      this.mdIconRegistry.addSvgIcon(
        icon,
        this.sanitizer.bypassSecurityTrustResourceUrl(`/assets/icons/${icon}.svg`)
      );
    });
    // Do not show the page if the feature is not enabled
    if (!this.config.releasesEnabled) {
      return this.router.navigate(['/404']);
    }

    this.route.params.forEach((params: Params) => {
      this.name = params['deploymentName'];
      this.loadDeployment(params['deploymentName']);
    });
  }

  loadDeployment(deploymentName: string): void {
    this.deploymentsService
      .getDeployment(deploymentName)
      .finally(() => {
        this.loading = false;
      })
      .subscribe(deployment => {
        this.deployment = deployment;
        this.resources = this.deploymentsService.loadResources(deployment);
      });
  }

  deploymentDeleted(event) {
    if (event.state == 'deleted') {
      return this.router.navigate(['/deployments']);
    }
  }

  getIconUrl(): string {
    if (this.deployment.attributes.chartIcon && !this.backgroundColor) {
      RGBaster.colors(this.deployment.attributes.chartIcon, {
        success: payload => {
          this.backgroundColor = payload.best
            .replace('rgb', 'rgba')
            .replace(')', ', 0.1)');
        }
      });
    }
    return this.deployment.attributes.chartIcon
      ? this.deployment.attributes.chartIcon
      : '/assets/images/placeholder.png';
  }
}
