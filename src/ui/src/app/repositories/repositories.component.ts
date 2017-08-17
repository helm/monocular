import { Component, OnInit } from '@angular/core';
import { ReposService } from '../shared/services/repos.service';
import { Repo } from '../shared/models/repo';
import { Router } from '@angular/router';
import { ConfigService } from '../shared/services/config.service';
import { DomSanitizer } from '@angular/platform-browser';
import { MdDialogRef, MdDialog, MdDialogConfig, MdIconRegistry, MdSnackBar } from '@angular/material';
import { DialogsService } from '../shared/services/dialogs.service';
import { RepositoryNewComponent } from '../repository-new/repository-new.component';

@Component({
  selector: 'app-repositories',
  templateUrl: './repositories.component.html',
  styleUrls: ['./repositories.component.scss'],
  viewProviders: [MdIconRegistry]
})
export class RepositoriesComponent implements OnInit {
  loading: boolean = true;
  repos: Repo[] = [];

  constructor(
    private reposService: ReposService,
    private router: Router,
    private config: ConfigService,
    private mdIconRegistry: MdIconRegistry,
    private sanitizer: DomSanitizer,
    private dialogsService: DialogsService,
    private dialog: MdDialog,
    public snackBar: MdSnackBar,
  ) {}

  ngOnInit() {
    this.mdIconRegistry.addSvgIcon(
      'delete',
      this.sanitizer.bypassSecurityTrustResourceUrl(
        '/assets/icons/delete.svg'
      )
    );
    // // Do not show the page if the feature is not enabled
    if (!this.config.releasesEnabled) {
      return this.router.navigate(['/404']);
    }
    this.loadRepos();
  }

  loadRepos(): void {
    this.reposService
      .getRepos()
      .finally(() => {
        this.loading = false;
      })
      .subscribe(repos => {
        this.repos = repos;
      });
  }

  goToRepoUrl(repo): string {
    return `/charts/${repo.attributes.name}`;
  }

  deleteRepo(repo: Repo) {
    this.dialogsService
    .confirm(
      `Remove ${repo.attributes.name} repository`,
      'You are going to remove this repository and all charts associated with it',
      'Remove it',
      'Cancel'
    )
    .subscribe(res => {
      if (res) {
        this.reposService.deleteRepo(repo.attributes.name)
          .subscribe(
            repo => {
              this.loadRepos();
            },
            error => {
              this.snackBar.open(
                `Error deleting the repository, please try later`,
                'close',
                { duration: 5000 }
              );
            }
          )
      }
    });
  }

  addRepo() {
    let dialogRef: MdDialogRef<RepositoryNewComponent>;
    dialogRef = this.dialog.open(RepositoryNewComponent);
    dialogRef.afterClosed().subscribe(res => this.loadRepos());
  }
}
