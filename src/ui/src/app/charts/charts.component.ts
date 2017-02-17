import { Component, OnInit } from '@angular/core';
import { ChartsService } from '../shared/services/charts.service';
import { Chart } from '../shared/models/chart';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { MetaService } from 'ng2-meta';

@Component({
  selector: 'app-charts',
  templateUrl: './charts.component.html',
  styleUrls: ['./charts.component.scss']
})
export class ChartsComponent implements OnInit {
  charts: Chart[] = [];
  orderedCharts: Chart[] = [];
  loading: boolean = true;
  currentRepo: string;

  constructor(
    private chartsService: ChartsService,
    private route: ActivatedRoute,
    private router: Router,
    private metaService: MetaService
  ) { }

  // Default filters
  filters = {
    orderBy: 'name'
  }

  ngOnInit() {
		this.loadCharts();
    this.updateMetaTags();
  }

  loadCharts(): void {
    this.route.params.forEach((params: Params) => {
      this.currentRepo = params["repo"]
  		this.chartsService.getCharts(this.currentRepo).subscribe(charts => {
        this.loading = false;
        this.charts = charts;
        this.orderedCharts = this.orderCharts(this.charts);
      });
    })
  }

  // Update a filter
  onChangeFilter(filter): void {
    if (this.filters[filter.type] !== filter.value) {
      // Repository change
      if (filter.type == "repositoryType") {
        return this.goToRepo(filter.value)
      }

      // in place filtering
      this.filters[filter.type] = filter.value;
      this.orderedCharts = this.orderCharts(this.charts);
    }
  }

  goToRepo(repo: string) {
    repo = repo === 'all' ? '' : repo;
    this.router.navigate(['/charts', repo]);
  }

  // Sort charts
  orderCharts(charts): Chart[] {
    switch(this.filters.orderBy) {
      case 'created': {
        return charts.sort(this.sortByCreated).reverse()
      }
      default: {
        return charts.sort((a, b) =>
          a.attributes.name.localeCompare(b.attributes.name));
      }
    }
  }

  private
  sortByCreated(a: Chart, b: Chart) {
      let aVersion = a.relationships.latestChartVersion.data
      let bVersion = b.relationships.latestChartVersion.data
      if(aVersion.created < bVersion.created){
          return -1
      } else if (aVersion.created > bVersion.created){
        return 1
      }
      return 0
  }

  updateMetaTags(): void {
    let title: string = `${this.currentRepo} repository charts`;
    this.metaService.setTitle(title);
    this.metaService.setTag('og:title', title);
  }
}
