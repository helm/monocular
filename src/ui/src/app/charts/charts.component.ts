import { Component, OnInit } from '@angular/core';
import { ChartsService } from '../shared/services/charts.service';
import { ReposService } from '../shared/services/repos.service';
import { Chart } from '../shared/models/chart';
import { Repo } from '../shared/models/repo';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { SeoService } from '../shared/services/seo.service';
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
  // Repos
  repoName: string;
  currentRepo: Repo;
  repositories: Repo[];

  constructor(
    private chartsService: ChartsService,
    private reposService: ReposService,
    private route: ActivatedRoute,
    private router: Router,
    private config: ConfigService,
    private seo: SeoService
  ) { }

  // Default filters
  filters = {
    orderBy: 'name'
  }

  ngOnInit() {
    this.route.params.forEach((params: Params) => {
      this.repoName = params['repo'];
      this.updateMetaTags();
      this.loadRepos();
      this.loadCharts();
    });
  }

  loadCharts(): void {
    this.chartsService.getCharts(this.repoName).subscribe(charts => {
      this.loading = false;
      this.charts = charts;
      this.orderedCharts = this.orderCharts(this.charts);
    });
  }

  loadRepos(): void {
    this.reposService.getRepos().subscribe(repos => {
      this.repositories = repos;
      if(this.repoName) {
        this.currentRepo = repos.filter(r => r.id == this.repoName)[0];
      }
    });
  }

  // Update a filter
  onChangeFilter(filter): void {
    if (this.filters[filter.type] !== filter.value) {
      // Repository change
      if (filter.type == "repositoryType") {
        return this.goToRepo(filter.value.id)
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
    if (this.repoName) {
      this.seo.setMetaTags('repoCharts', { repo: this.capitalize(this.repoName) });
    } else {
      this.seo.setMetaTags('charts');
    }
  }

  capitalize(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
  }
}
