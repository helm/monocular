/* tslint:disable:no-unused-variable */

import { TestBed, async } from '@angular/core/testing';
import { MainHeaderComponent } from './main-header.component';
import { Router } from '@angular/router';
import { MaterialModule } from '@angular/material';
import { ConfigService } from '../shared/services/config.service';
import { MenuService } from '../shared/services/menu.service';
import { HeaderBarComponent } from '../header-bar/header-bar.component';

describe('Component: MainHeader', () => {
  beforeEach(
    async(() => {
      TestBed.configureTestingModule({
        declarations: [MainHeaderComponent, HeaderBarComponent],
        imports: [MaterialModule],
        providers: [
          { provide: Router },
          { provide: ConfigService, useValue: { appName: 'app-name' } },
          { provide: MenuService }
        ]
      }).compileComponents();
    })
  );

  it('should create an instance', () => {
    let component = TestBed.createComponent(MainHeaderComponent);
    expect(component).toBeTruthy();
  });
});
