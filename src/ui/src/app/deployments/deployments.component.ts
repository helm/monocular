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
  visibleDeployments: Deployment[] = [];
  namespaces: string[] = ['All'];
  orders: string[] = ['Name', 'Date', 'Status'];
  orderBy: string = 'Date';
  namespace: string = 'All';
  loading: boolean = true;

  constructor(
    private deploymentsService: DeploymentsService,
    private router: Router,
    private config: ConfigService,
    private mdIconRegistry: MdIconRegistry,
    private sanitizer: DomSanitizer
  ){
    mdIconRegistry
      .addSvgIcon('search',
        sanitizer.bypassSecurityTrustResourceUrl(`/assets/icons/search.svg`));
  }

  ngOnInit() {
    // Do not show the page if the feature is not enabled
    if(!this.config.releasesEnabled) {
      return this.router.navigate(['/404']);
    }
    this.loadDeployments();
  }

  loadDeployments(): void {
    this.deploymentsService.getDeployments()
    .finally(()=> {
      this.loading = false;
    }).subscribe(deployments => {
      this.deployments = deployments;
      this.filterDeployments();
      this.namespaces = this.exportNamespaces();
    })
  }

  exportNamespaces(): string[] {
    var list: string[] = ['All'];
    this.deployments.forEach(dp => {
      if (list.indexOf(dp.attributes.namespace) == -1) {
        list.push(dp.attributes.namespace);
      }
    })
    return list;
  }

  filterDeployments() {
    let filtered = this.deployments
    if (this.namespace !== 'All') {
      filtered = filtered.filter(deployment => {
        return deployment.attributes.namespace === this.namespace;
      })
    }
    filtered = filtered.sort((a, b) => {
      if (this.orderBy === 'Name') {
        return a.id <= b.id ? -1 : 1;
      } else if (this.orderBy === 'Status') {
        return a.attributes.status <= b.attributes.status ? -1 : 1;
      } else {
        return a.attributes.updated <= b.attributes.updated ? -1 : 1;
      }
    })
    this.visibleDeployments = filtered;
  }

  searchChange(e) {
    let newValue = e.target.value;
    if (!newValue) {
      return this.filterDeployments();
    }
    let searchTerm = newValue.toLowerCase();
    this.visibleDeployments = this.deployments.filter(deployment => {
      return deployment.id.indexOf(searchTerm) != -1;
    })
  }

  clickNamespace(ns) {
    this.namespace = ns;
    this.filterDeployments();
  }

  clickOrderBy(order) {
    this.orderBy = order;
    this.filterDeployments();
  }

}
