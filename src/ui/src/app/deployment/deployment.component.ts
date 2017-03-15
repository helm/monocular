import { Router, ActivatedRoute, Params } from '@angular/router';
import { Component, OnInit } from '@angular/core';
import { DeploymentsService } from '../shared/services/deployments.service';
import { Deployment } from '../shared/models/deployment';
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
  ) {
    const icons = ['layers', 'schedule', 'web-asset', 'info-outline', 'arrow-back'];

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
      this.name = params['deploymentName'];
      this.loadDeployment(params['deploymentName']);
    });

  }

  /**
   * Prepare the resources for displaying in the UI.
   *
   * TODO: In the future, the backend will provide this information
   */
  loadResources(deployment: Deployment): any {
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

  loadDeployment(deploymentName: string): void {
    this.deploymentsService.getDeployment(deploymentName)
    .finally(()=> {
      this.loading = false;
    }).subscribe(deployment => {
      this.deployment = deployment;
      this.resources = this.loadResources(deployment);
    })
  }

  deploymentDeleted(event) {
    if (event.state == "deleted") {
      return this.router.navigate(['/deployments']);
    }
  }
}
