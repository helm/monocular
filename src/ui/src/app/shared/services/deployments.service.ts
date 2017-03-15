import { Injectable } from '@angular/core';
import { Deployment } from '../models/deployment';
import { ConfigService } from './config.service';

import { Observable } from 'rxjs';
import 'rxjs/add/operator/switchMap';
import 'rxjs/add/operator/find';
import 'rxjs/add/operator/map';

import { Http, Response } from '@angular/http';

/* TODO, This is a mocked class. */
@Injectable()
export class DeploymentsService {
  hostname: string;

  constructor(
    private http: Http,
    private config: ConfigService
  ) {
    this.hostname = config.backendHostname;
  }

  getDeployments(): Observable<Deployment[]> {
      return this.http.get(`${this.hostname}/v1/releases`)
                    .map((response) => {
                      return this.extractData(response, [])
                    }).catch(this.handleError);
  }

  getDeployment(deploymentName: string): Observable<Deployment> {
      return this.http.get(`${this.hostname}/v1/releases/${deploymentName}`)
                    .map((response) => {
                      return this.extractData(response, [])
                    }).catch(this.handleError);
  }

  installDeployment(chartID: string, version: string): Observable<Deployment> {
      var params = { "chartId": chartID, "chartVersion": version }
      return this.http.post(`${this.hostname}/v1/releases`, params)
                    .map(this.extractData)
                    .catch(this.handleError);
  }

  deleteDeployment(deploymentName: string): Observable<Deployment> {
    return this.http.delete(`${this.hostname}/v1/releases/${deploymentName}`)
                    .map(this.extractData)
                    .catch(this.handleError);
  }

  private extractData(res: Response, fallback = {}) {
    let body = res.json();
    return body.data || fallback;
  }

  private handleError (error: any) {
    error = error.json();
    let errMsg = (error.message) ? error.message :
      error.status ? `${error.status} - ${error.statusText}` : 'Server error';
    console.error(errMsg); // log to console instead
    return Observable.throw(errMsg);
  }
}
