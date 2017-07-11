import { Observable } from 'rxjs/Rx';
import { ConfirmDialogComponent } from '../../confirm-dialog/confirm-dialog.component';
import { MdDialogRef, MdDialog, MdDialogConfig } from '@angular/material';
import { Injectable } from '@angular/core';

@Injectable()
export class DialogsService {
  constructor(private dialog: MdDialog) {}

  public confirm(
    title: string,
    message: string,
    ok = 'Continue',
    cancel = 'Cancel',
    actionButtonClass = 'primary'
  ): Observable<boolean> {
    let dialogRef: MdDialogRef<ConfirmDialogComponent>;

    dialogRef = this.dialog.open(ConfirmDialogComponent);
    dialogRef.componentInstance.title = title;
    dialogRef.componentInstance.message = message;
    dialogRef.componentInstance.actionButtonClass = actionButtonClass;
    dialogRef.componentInstance.ok = ok;
    dialogRef.componentInstance.cancel = cancel;

    return dialogRef.afterClosed();
  }
}
