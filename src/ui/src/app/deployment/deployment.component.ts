import { Router, ActivatedRoute, Params } from '@angular/router';
import { Component, OnInit } from '@angular/core';
import { ReleasesService } from '../shared/services/releases.service';
import { Release } from '../shared/models/release';
import { Chart } from '../shared/models/chart';
import { ConfigService } from '../shared/services/config.service';
import { DomSanitizer } from '@angular/platform-browser';
import { MdIconRegistry } from '@angular/material';

@Component({
  selector: 'app-deployment',
  templateUrl: './deployment.component.html',
  styleUrls: ['./deployment.component.scss'],
  viewProviders: [MdIconRegistry]
})
export class DeploymentComponent implements OnInit {
  deployment: Release;
  resources = [];
  loading: boolean = true;

  constructor(
    private releasesService: ReleasesService,
    private router: Router,
    private route: ActivatedRoute,
    private config: ConfigService,
    private mdIconRegistry: MdIconRegistry,
    private sanitizer: DomSanitizer
  ) {
    const icons = ['layers', 'schedule', 'web-asset', 'info-outline'];

    icons.forEach(icon => {
      mdIconRegistry
        .addSvgIcon(icon,
          sanitizer.bypassSecurityTrustResourceUrl(`/assets/icons/${icon}.svg`));
    });
  }

  ngOnInit() {
    // Do not show the page if the feature is not enabled
    if(!this.config.releasesEnabled) {
      return this.router.navigate(['/404']);
    }

    this.route.params.forEach((params: Params) => {
      this.loadRelease(params['deploymentName']);
    });

  }

  /**
   * Prepare the resources for displaying in the UI.
   *
   * TODO: In the future, the backend will provide this information
   */
  loadResources(deployment: Release): any {
    let elements = deployment.attributes.resources.split('=='),
      resources = [];

    // Remove first element
    elements.shift();

    // Regex
    let nameRegex = /^> [\w\d\s\/]+\/(\w+)+/;

    elements.forEach(el => {
      let lines = el.split("\n");

      // Name
      let name = nameRegex.exec(lines.shift())[1];
      let headers = lines.shift().split(/\s+/);
      let services = [];

      // Remaining lines
      lines.forEach(line => {
        if (line !== ''){
          let values = line.split(/\s+/);
          let service = {};

          values.forEach((value, i) => {
            service[headers[i]] = value;
          });

          // Add to the array
          services.push(service);
        }
      });

      // Build the resource
      resources.push({ name, services });
    });

    return resources;
  }

  loadRelease(deploymentName: string): void {
    this.releasesService.getRelease(deploymentName)
    .finally(()=> {
      this.loading = false;
    }).subscribe(deployment => {
      this.deployment = deployment;
      this.resources = this.loadResources(deployment);
    })
  }

  releaseDeleted(event) {
    if (event.state == "deleted") {
      return this.router.navigate(['/deployments']);
    }
  }
}
