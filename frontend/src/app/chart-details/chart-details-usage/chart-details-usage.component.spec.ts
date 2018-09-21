/* tslint:disable:no-unused-variable */

import { TestBed, async } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ChartDetailsUsageComponent } from './chart-details-usage.component';

describe('Component: ChartDetailsUsage', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [],
      declarations: [ChartDetailsUsageComponent],
      providers: [
      ],
      schemas: [NO_ERRORS_SCHEMA]
    }).compileComponents();
  });

  it('should create an instance', () => {
    let component = TestBed.createComponent(ChartDetailsUsageComponent);
    expect(component).toBeTruthy();
  });
});
