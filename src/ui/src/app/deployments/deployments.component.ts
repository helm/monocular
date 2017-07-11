import { Component, OnInit } from '@angular/core';
import { DeploymentsService } from '../shared/services/deployments.service';
import { Deployment } from '../shared/models/deployment';
import { Router } from '@angular/router';
import { ConfigService } from '../shared/services/config.service';
import { DomSanitizer } from '@angular/platform-browser';
import { MdIconRegistry } from '@angular/material';

@Component({
  selector: 'app-deployments',
  templateUrl: './deployments.component.html',
  styleUrls: ['./deployments.component.scss'],
  viewProviders: [MdIconRegistry]
})
export class DeploymentsComponent implements OnInit {
  deployments: Deployment[] = [];
  loading: boolean = true;

  constructor(
    private deploymentsService: DeploymentsService,
    private router: Router,
    private config: ConfigService,
    private mdIconRegistry: MdIconRegistry,
    private sanitizer: DomSanitizer
  ) {
    mdIconRegistry.addSvgIcon(
      'layers',
      sanitizer.bypassSecurityTrustResourceUrl(`/assets/icons/layers.svg`)
    );
  }

  ngOnInit() {
    // Do not show the page if the feature is not enabled
    if (!this.config.releasesEnabled) {
      return this.router.navigate(['/404']);
    }
    this.loadDeployments();
  }

  loadDeployments(): void {
    this.deploymentsService
      .getDeployments()
      .finally(() => {
        this.loading = false;
      })
      .subscribe(deployments => {
        this.deployments = deployments;
      });
  }
}
