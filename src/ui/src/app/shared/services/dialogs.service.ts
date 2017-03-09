import { Observable } from 'rxjs/Rx';
import { ConfirmDialog } from '../../confirm-dialog/confirm-dialog.component';
import { MdDialogRef, MdDialog, MdDialogConfig } from '@angular/material';
import { Injectable } from '@angular/core';

@Injectable()
export class DialogsService {

    constructor(private dialog: MdDialog) { }

    public confirm(title: string, message: string, ok="Continue", cancel="Cancel"): Observable<boolean> {

        let dialogRef: MdDialogRef<ConfirmDialog>;

        dialogRef = this.dialog.open(ConfirmDialog);
        dialogRef.componentInstance.title = title;
        dialogRef.componentInstance.message = message;
        dialogRef.componentInstance.ok = ok;
        dialogRef.componentInstance.cancel = cancel;

        return dialogRef.afterClosed();
    }
}
