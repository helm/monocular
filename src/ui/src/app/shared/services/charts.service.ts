import { Injectable } from '@angular/core';
import { Chart } from '../models/chart';
import { CONFIG } from '../../config';

// To get the Mocked Readme file
import { Observable } from 'rxjs';
import 'rxjs/add/operator/switchMap';
import 'rxjs/add/operator/find';
import 'rxjs/add/operator/map';

import { Http, Response } from '@angular/http';

/* TODO, This is a mocked class. */
@Injectable()
export class ChartsService {
  hostname: String = CONFIG.backendHostname

  constructor(private http: Http) { }

  /**
   * Get all charts from the API
   *
   * @return {Observable} An observable that will an array with all Charts
   */
  getCharts(): Observable<Chart[]> {
    return this.http.get(`${this.hostname}/v1/charts`)
                  .map(this.extractData)
                  .catch(this.handleError);
  }

  /**
   * Get a chart using the API
   *
   * @param {String} repo Repository name
   * @param {String} chartName Chart name
   * @return {Observable} An observable that will a chart instance
   */
  getChart(repo: String, chartName: String): Observable<Chart> {
    // Transform Observable<Chart[]> into Observable<Chart>[]
    return this.http.get(`${this.hostname}/v1/charts/${repo}/${chartName}`)
                  .map(this.extractData)
                  .catch(this.handleError);
  }

  /* TODO, use backend search API endpoint */
  searchCharts(query): Observable<Chart[]> {
    let re = new RegExp(query, 'i');
    return this.getCharts().map(charts => {
      return charts.filter(chart => {
        return chart.attributes.name.match(re) || chart.attributes.description.match(re)
      })
    })
  }

  /* TODO, get remote README */
  getMockedReadme(): Observable<Response> {
    let readmeUrl = '/assets/mock_readme.md'
    return this.http.get(readmeUrl)
  }

  private extractData(res: Response) {
    let body = res.json();
    return body.data || { };
  }

  private handleError (error: any) {
    let errMsg = (error.message) ? error.message :
      error.status ? `${error.status} - ${error.statusText}` : 'Server error';
    console.error(errMsg); // log to console instead
    return Observable.throw(errMsg);
  }

}
