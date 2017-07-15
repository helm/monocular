/* tslint:disable:no-unused-variable */

import { TestBed, async } from '@angular/core/testing';
import { MaterialModule } from '@angular/material';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ChartSearchComponent } from './chart-search.component';
import { ChartListComponent } from '../chart-list/chart-list.component';
import { ChartItemComponent } from '../chart-item/chart-item.component';
import { LoaderComponent } from '../loader/loader.component';
import { PanelComponent } from '../panel/panel.component';
import { HeaderBarComponent } from '../header-bar/header-bar.component';
import { ChartsService } from '../shared/services/charts.service';
import { ActivatedRoute, Router } from '@angular/router';
import { SeoService } from '../shared/services/seo.service';
import { MenuService } from '../shared/services/menu.service';
import { ConfigService } from '../shared/services/config.service';

describe('Component: ChartSearch', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [MaterialModule],
      declarations: [
        ChartSearchComponent,
        ChartListComponent,
        ChartItemComponent,
        LoaderComponent,
        PanelComponent,
        HeaderBarComponent
      ],
      providers: [
        ChartsService,
        ConfigService,
        MenuService,
        { provide: SeoService },
        { provide: Router },
        { provide: ActivatedRoute }
      ],
      schemas: [NO_ERRORS_SCHEMA]
    }).compileComponents();
  });

  it('should create an instance', () => {
    let component = TestBed.createComponent(ChartSearchComponent);
    expect(component).toBeTruthy();
  });
});
