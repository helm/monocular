import { Component, Output, EventEmitter } from '@angular/core';
import { ReleasesService } from '../shared/services/releases.service';
import { Release } from '../shared/models/release';
import { MdIconRegistry } from '@angular/material';
import { DomSanitizer } from '@angular/platform-browser';
import { MdSnackBar } from '@angular/material';
import { DialogsService } from '../shared/services/dialogs.service';

@Component({
  selector: 'app-release-item',
  templateUrl: './release-item.component.html',
  styleUrls: ['./release-item.component.scss'],
  inputs: ['release', 'extended']
})
export class ReleaseItemComponent {
  @Output() onDelete = new EventEmitter();
  deleting: boolean

  constructor(
    private releasesService: ReleasesService,
    private mdIconRegistry: MdIconRegistry,
    private sanitizer: DomSanitizer,
    private dialogsService: DialogsService,
    public snackBar: MdSnackBar
  ){
    mdIconRegistry
      .addSvgIcon('delete',
        sanitizer.bypassSecurityTrustResourceUrl('/assets/icons/delete.svg'));
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
    this.deleting = true;
    this.onDelete.emit({ name: releaseName, state: "deleting" });
    this.snackBar.open("Deleting release", "close", {});
    this.releasesService.deleteRelease(releaseName)
    .finally(() => {
      this.deleting = false
    }).subscribe(
      release => {
        this.onDelete.emit({ name: releaseName, state: "deleted" });
        this.snackBar.open("Release deleted", "", { duration: 2500 });
      },
      error => {
        this.snackBar.open("Error deleting the release", "", { duration: 2500 });
      }
    );
  }
}
