/* tslint:disable:no-unused-variable */

import { TestBed, async } from '@angular/core/testing';
import { MdDialogRef, MaterialModule } from '@angular/material';
import { ConfirmDialogComponent } from './confirm-dialog.component';

describe('Component: ConfirmDialog', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [MaterialModule],
      declarations: [ConfirmDialogComponent],
      providers: [{ provide: MdDialogRef }]
    }).compileComponents();
  });

  it('should create an instance', () => {
    let component = TestBed.createComponent(ConfirmDialogComponent);
    expect(component).toBeTruthy();
  });
});
