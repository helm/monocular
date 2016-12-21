import { Component, OnInit } from '@angular/core';
import { ChartsService } from '../charts.service';
import { Chart } from '../chart';

import { ActivatedRoute } from '@angular/router';
import { Observable }         from 'rxjs/Observable';
import 'rxjs/add/operator/map';

@Component({
  selector: 'app-chart-search',
  templateUrl: './chart-search.component.html',
  styleUrls: ['./chart-search.component.scss']
})
export class ChartSearchComponent implements OnInit {
  query: String;
	charts: Chart[] = []

  constructor(
    private route: ActivatedRoute,
    private chartsService: ChartsService
  ) { }

  ngOnInit() {
    this.route
      .queryParams
      .forEach(params => {
        let q: String = params['q']
        this.query = q
        this.searchCharts(q)
      })
  }

  searchCharts(q: String): void {
		this.chartsService.searchCharts(q).subscribe(charts => this.charts = charts);
  }

  resultMessage(): String {
    if (this.charts.length > 0) {
      return this.charts.length + " results found for \"" + this.query + "\"";
    } else {
      return "\"" + this.query + "\" did not return any results";
    }
  }
}
