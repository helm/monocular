import { Router, ActivatedRoute, Params } from '@angular/router';
import { Component, OnInit } from '@angular/core';
import { ReleasesService } from '../shared/services/releases.service';
import { Release } from '../shared/models/release';
import { ConfigService } from '../shared/services/config.service';

@Component({
  selector: 'app-release',
  templateUrl: './release.component.html',
  styleUrls: ['./release.component.scss']
})
export class ReleaseComponent implements OnInit {
  release: Release;
  loading: boolean = true;

  constructor(
    private releasesService: ReleasesService,
    private router: Router,
    private route: ActivatedRoute,
    private config: ConfigService
  ){ }

  ngOnInit() {
    // Do not show the page if the feature is not enabled
    if(!this.config.releasesEnabled) {
      return this.router.navigate(['/404']);
    }

    this.route.params.forEach((params: Params) => {
      this.loadRelease(params["releaseName"]);
    });

  }

  loadRelease(releaseName: string): void {
    this.releasesService.getRelease(releaseName)
    .finally(()=> {
      this.loading = false;
    }).subscribe(release => {
      this.release = release;
    })
  }

  releaseDeleted(event) {
    if (event.state == "deleted") {
      return this.router.navigate(['/releases']);
    }
  }
}
