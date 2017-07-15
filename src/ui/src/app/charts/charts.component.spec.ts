/* tslint:disable:no-unused-variable */

import { TestBed, async } from '@angular/core/testing';
import { MaterialModule } from '@angular/material';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ChartsComponent } from './charts.component';
import { ChartListComponent } from '../chart-list/chart-list.component';
import { ChartsFiltersComponent } from '../charts-filters/charts-filters.component';
import { ChartItemComponent } from '../chart-item/chart-item.component';
import { LoaderComponent } from '../loader/loader.component';
import { PanelComponent } from '../panel/panel.component';
import { HeaderBarComponent } from '../header-bar/header-bar.component';
import { ChartsService } from '../shared/services/charts.service';
import { ReposService } from '../shared/services/repos.service';
import { ActivatedRoute, Router } from '@angular/router';
import { SeoService } from '../shared/services/seo.service';
import { MenuService } from '../shared/services/menu.service';
import { ConfigService } from '../shared/services/config.service';

describe('Component: Charts', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [MaterialModule],
      declarations: [
        ChartsComponent,
        ChartListComponent,
        ChartsFiltersComponent,
        ChartItemComponent,
        LoaderComponent,
        PanelComponent,
        HeaderBarComponent
      ],
      providers: [
        ChartsService,
        ConfigService,
        ReposService,
        MenuService,
        { provide: SeoService },
        { provide: ActivatedRoute },
        { provide: Router }
      ],
      schemas: [NO_ERRORS_SCHEMA]
    }).compileComponents();
  });

  it('should create an instance', () => {
    let component = TestBed.createComponent(ChartsComponent);
    expect(component).toBeTruthy();
  });
});
