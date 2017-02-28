import { Component, OnInit, Input, ViewEncapsulation, ChangeDetectionStrategy } from '@angular/core';
import { Chart } from '../../shared/models/chart';
import { DomSanitizer } from '@angular/platform-browser';
import { MdIconRegistry } from '@angular/material';
import { MdSnackBar } from '@angular/material';

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

  constructor(
    mdIconRegistry: MdIconRegistry,
    sanitizer: DomSanitizer,
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

  // Deletes /stable prefix not needed for stable repos
  get cmdChartId(): string {
    return this.chart.id.replace("stable/", "")
  }

  get showRepoInstructions(): boolean {
    return this.chart.attributes.repo.name != 'stable'
  }

  get repoAddInstructions(): string {
    return `helm repo add incubator ${this.chart.attributes.repo.registryURL}`;
  }

  get installInstructions(): string {
    return `helm install ${this.cmdChartId} --version ${this.currentVersion}`;
  }
}
