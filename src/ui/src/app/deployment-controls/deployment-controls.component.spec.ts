/* tslint:disable:no-unused-variable */

import { TestBed, async } from '@angular/core/testing';
import { DeploymentControlsComponent } from './deployment-controls.component';
import { DeploymentsService } from '../shared/services/deployments.service';
import { DialogsService } from '../shared/services/dialogs.service';
import { ConfigService } from '../shared/services/config.service';
import { MaterialModule } from '@angular/material';

describe('Component: DeploymentControls', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [MaterialModule],
      declarations: [DeploymentControlsComponent],
      providers: [DeploymentsService, DialogsService, ConfigService]
    }).compileComponents();
  });

  it('should create an instance', () => {
    let component = TestBed.createComponent(DeploymentControlsComponent);
    expect(component).toBeTruthy();
  });
});
