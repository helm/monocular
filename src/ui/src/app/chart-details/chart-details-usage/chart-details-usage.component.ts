import { Component, OnInit, Input, ViewEncapsulation, ChangeDetectionStrategy } from '@angular/core';
import { Chart } from '../../shared/models/chart';
import { Release } from '../../shared/models/release';
import { DomSanitizer } from '@angular/platform-browser';
import { MdIconRegistry } from '@angular/material';
import { MdSnackBar } from '@angular/material';
import { ConfigService } from '../../shared/services/config.service';
import { DialogsService } from '../../shared/services/dialogs.service';

import { ReleasesService } from '../../shared/services/releases.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-chart-details-usage',
  templateUrl: './chart-details-usage.component.html',
  styleUrls: ['./chart-details-usage.component.scss'],
  viewProviders: [MdIconRegistry],
  encapsulation: ViewEncapsulation.None
})
export class ChartDetailsUsageComponent implements OnInit {
  @Input() chart: Chart
  @Input() currentVersion: string
  installing: boolean

  constructor(
    mdIconRegistry: MdIconRegistry,
    sanitizer: DomSanitizer,
    private config: ConfigService,
    private releasesService: ReleasesService,
    private router: Router,
    private dialogsService: DialogsService,
    public snackBar: MdSnackBar
  ) {
    mdIconRegistry
      .addSvgIcon('content-copy',
        sanitizer.bypassSecurityTrustResourceUrl('/assets/icons/content-copy.svg'));
  }

  ngOnInit() {}

  // Show an snack bar to confirm the user that the code has been copied
  showSnackBar(): void {
    this.snackBar.open('Copied to the clipboard', '', {
      duration: 1500,
    });
  }

  get showRepoInstructions(): boolean {
    return this.chart.attributes.repo.name != 'stable'
  }

  get repoAddInstructions(): string {
    return `helm repo add ${this.chart.attributes.repo.name} ${this.chart.attributes.repo.URL}`;
  }

  get installInstructions(): string {
    return `helm install ${this.chart.id} --version ${this.currentVersion}`;
  }

  installRelease(chartID: string, version: string): void {
    this.dialogsService
      .confirm(`You will deploy ${chartID} v${version}`, '')
      .subscribe(res => {
        if (res)
          this.performInstallation(chartID, version);
      });

  }

  performInstallation(chartID: string, version: string): void {
    this.installing = true;

    this.releasesService.installRelease(chartID, version)
    .finally(() => {
      this.installing = false
    }).subscribe(
      release => {
        this.installOK(release)
      },
      error => {
        this.snackBar.open(`Error installing the application, please try later"`, 'close', { duration: 5000 });
      }
    );
  }

  installOK(release: Release) :void {
      let message = this.snackBar.open('Installation completed', 'view more', {
      });

      message.onAction().subscribe(() => {
        this.router.navigate(['/releases']);
      });

  }
}
