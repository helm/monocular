import { MdDialogRef } from '@angular/material';
import { Component } from '@angular/core';

@Component({
    selector: 'confirm-dialog',
    template: `
        <p>{{ title }}</p>
        <p>{{ message }}</p>
        <button type="button" md-raised-button color="primary"
            (click)="dialogRef.close(true)">{{ ok }}</button>
        <button type="button" md-button md-raised-button
            (click)="dialogRef.close()">{{ cancel }}</button>
    `,
})
export class ConfirmDialog {

    public title: string;
    public message: string;
    public ok: string;
    public cancel: string;

    constructor(public dialogRef: MdDialogRef<ConfirmDialog>) {

    }
}
