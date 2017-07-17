import { Component, OnInit } from '@angular/core';
import { ChartsService } from '../shared/services/charts.service';
import { MetaService } from 'ng2-meta';
import { Chart } from '../shared/models/chart';
import { ConfigService } from '../shared/services/config.service';
import { SeoService } from '../shared/services/seo.service';

import { ActivatedRoute } from '@angular/router';
import { Observable }         from 'rxjs/Observable';
import 'rxjs/add/operator/map';

@Component({
  selector: 'app-chart-search',
  templateUrl: './chart-search.component.html',
  styleUrls: ['./chart-search.component.scss']
})
export class ChartSearchComponent implements OnInit {
  query: string;
  loading: boolean = true;
	charts: Chart[] = [];

  constructor(
    private route: ActivatedRoute,
    private chartsService: ChartsService,
    private config: ConfigService,
    private seo: SeoService
  ) { }

  ngOnInit() {
    this.route
      .queryParams
      .forEach(params => {
        let q: string = params['q']
        this.query = q;
        this.searchCharts(q);
      });

    // Update meta tags
    this.seo.setMetaTags('search', { search: this.query });
  }

  searchCharts(q: string): void {
		this.chartsService.searchCharts(q).subscribe(charts => {
      this.loading = false;
      this.charts = charts;
    });
  }

  resultMessage(): string {
    if (this.charts.length > 0) {
      return `${this.charts.length} results found for "${this.query}"`;
    } else {
      return `"${this.query}" did not return any results`;
    }
  }
}
