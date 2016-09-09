import { Injectable } from '@angular/core';
import { Chart } from './chart';
import { CHARTS } from './chart-mock';

// To get the Mocked Readme file
import { Observable } from 'rxjs';
import { Http, Response } from '@angular/http';

/* TODO, This is a mocked class. */
@Injectable()
export class ChartsService {

  constructor(private http: Http) { }

  getCharts(): Promise<Chart[]> {
    return Promise.resolve(CHARTS);
  }

  getChart(repo: String, chartName: String): Promise<Chart> {
    let found: Chart[]
    found = CHARTS.filter(chart => {
      return chart.attributes.repo == repo && chart.attributes.name == chartName
    })
    return Promise.resolve(found[0]);
  }

  searchCharts(query): Promise<Chart[]> {
    let found: Chart[]
    let re = new RegExp(query, 'i');
    // Mocked. The backend will be the one returning search results
    found = CHARTS.filter(chart => {
      return chart.attributes.name.match(re) || chart.attributes.description.match(re)
    })
    return Promise.resolve(found);
  }

  getMockedReadme(): Observable<Response> {
    let readmeUrl = '/assets/mock_readme.md'
    return this.http.get(readmeUrl)
  }
}
