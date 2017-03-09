import { Injectable } from '@angular/core';
import { Release } from '../models/release';
import { ConfigService } from './config.service';

import { Observable } from 'rxjs';
import 'rxjs/add/operator/switchMap';
import 'rxjs/add/operator/find';
import 'rxjs/add/operator/map';

import { Http, Response } from '@angular/http';

/* TODO, This is a mocked class. */
@Injectable()
export class ReleasesService {
  hostname: string;

  constructor(
    private http: Http,
    private config: ConfigService
  ) {
    this.hostname = config.backendHostname;
  }

  /**
   * Get all charts from the API
   *
   * @return {Observable} An observable that will an array with all Charts
   */
  getReleases(): Observable<Release[]> {
      return this.http.get(`${this.hostname}/v1/releases`)
                    .map((response) => {
                      return this.extractData(response, [])
                    }).catch(this.handleError);
  }

  installRelease(chartID: string, version: string): Observable<Release> {
      var params = { "chartId": chartID, "chartVersion": version }
      return this.http.post(`${this.hostname}/v1/releases`, params)
                    .map(this.extractData)
                    .catch(this.handleError);
  }

  deleteRelease(releaseName: string): Observable<Release> {
    return this.http.delete(`${this.hostname}/v1/releases/${releaseName}`)
                    .map(this.extractData)
                    .catch(this.handleError);
  }

  private extractData(res: Response, fallback = {}) {
    let body = res.json();
    return body.data || fallback;
  }

  private handleError (error: any) {
    let errMsg = (error.message) ? error.message :
      error.status ? `${error.status} - ${error.statusText}` : 'Server error';
    console.error(errMsg); // log to console instead
    return Observable.throw(errMsg);
  }
}
