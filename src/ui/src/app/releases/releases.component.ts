import { Component, OnInit } from '@angular/core';
import { ReleasesService } from '../shared/services/releases.service';
import { Release } from '../shared/models/release';
import { MdIconRegistry } from '@angular/material';
import { DomSanitizer } from '@angular/platform-browser';
import { MdSnackBar } from '@angular/material';
import { DialogsService } from '../shared/services/dialogs.service';
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
    private mdIconRegistry: MdIconRegistry,
    private sanitizer: DomSanitizer,
    private dialogsService: DialogsService,
    public snackBar: MdSnackBar,
    private router: Router,
    private config: ConfigService
  ){
    mdIconRegistry
      .addSvgIcon('delete',
        sanitizer.bypassSecurityTrustResourceUrl('/assets/icons/delete.svg'));
  }

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

  deleteRelease(releaseName: string): void {
    this.dialogsService
      .confirm(`Do you want to delete "${releaseName}"?`, '')
      .subscribe(res => {
        if(res)
          this.performDelete(releaseName);
      })
  }

  performDelete(releaseName: string): void {
    this.snackBar.open("Deleting release", "close", {});
    this.releases =  this.releases.filter(item => item.id !== releaseName);
    this.releasesService.deleteRelease(releaseName).subscribe(
      release => {
        this.snackBar.open("Release deleted", "", { duration: 2500 });
      },
      error => {
        this.snackBar.open("Error deleting the release", "", { duration: 2500 });
      }
    );
  }
}
