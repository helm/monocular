import { Component, OnInit } from '@angular/core';
import { ReleasesService } from '../shared/services/releases.service';
import { Release } from '../shared/models/release';
import { Router } from '@angular/router';
import { ConfigService } from '../shared/services/config.service';

@Component({
  selector: 'app-releases',
  templateUrl: './releases.component.html',
  styleUrls: ['./releases.component.scss']
})
export class ReleasesComponent implements OnInit {
  releases: Release[] = [];
  loading: boolean = true;

  constructor(
    private releasesService: ReleasesService,
    private router: Router,
    private config: ConfigService
  ){ }

  ngOnInit() {
    // Do not show the page if the feature is not enabled
    if(!this.config.releasesEnabled) {
      return this.router.navigate(['/404']);
    }
    this.loadReleases();
  }

  loadReleases(): void {
    this.releasesService.getReleases()
    .finally(()=> {
      this.loading = false;
    }).subscribe(releases => {
      this.releases = releases;
    })
  }

  releaseDeleted(event): void {
    // Optimist update
    if (event.state == "deleting") {
      this.releases =  this.releases.filter(item => item.id !== event.name);
    }
  }
}
