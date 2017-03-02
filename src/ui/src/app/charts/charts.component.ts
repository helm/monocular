import { Component, OnInit } from '@angular/core';
import { ChartsService } from '../shared/services/charts.service';
import { Chart } from '../shared/models/chart';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { MetaService } from 'ng2-meta';
import { ConfigService } from '../shared/services/config.service';

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
  repositories: string[];

  constructor(
    private chartsService: ChartsService,
    private route: ActivatedRoute,
    private router: Router,
    private config: ConfigService,
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
        this.repositories = this.getAvailableRepos(charts);
      });
    })
  }

  // Get The list of repositories our charts are placed in
  // TODO: This should be retrieved via an API call
  getAvailableRepos(charts: Chart[]): string[] {
    console.warn("checking repos")
    var unique = {};
    var repos = [];
    charts.forEach(function (chart) {
      if( typeof(unique[chart.attributes.repo.name]) == "undefined"){
        repos.push(chart.attributes.repo.name);
      }
      unique[chart.attributes.repo.name] = 0;
    })
    console.warn(repos)
    return repos;
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
    let title: string = `${this.currentRepo || "stable"} repository charts`;
    this.metaService.setTitle(title, ` | ${this.config.appName}`);
    this.metaService.setTag('og:title', title);
  }
}
