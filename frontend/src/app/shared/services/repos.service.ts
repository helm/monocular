import { Injectable } from '@angular/core';
import { RepoAttributes } from '../models/repo';
import { ConfigService } from './config.service';

import { Observable } from 'rxjs';
import 'rxjs/add/operator/switchMap';
import 'rxjs/add/operator/find';
import 'rxjs/add/operator/map';

import { Http, Response } from '@angular/http';

@Injectable()
export class ReposService {
  constructor(
    private http: Http,
  ) {
  }

  /**
   * Get all repos from the API
   *
   * @return {Observable} An observable that will an array with all repos
   */
  getRepos(): Observable<RepoAttributes[]> {
    return this.http.get(`/assets/js/repos.json`)
                  .map(this.extractData)
                  .catch(this.handleError);
  }

  private extractData(res: Response) {
    let body = res.json();
    return body.data || { };
  }

  private handleError (error: any) {
    let errMsg = (error.json().message) ? error.json().message :
      error.status ? `${error.status} - ${error.statusText}` : 'Server error';
    console.error(errMsg); // log to console instead
    return Observable.throw(errMsg);
  }
}
