import { Component, OnInit } from '@angular/core';
import { ChartsService } from '../shared/services/charts.service';
import { ReposService } from '../shared/services/repos.service';
import { Chart } from '../shared/models/chart';
import { Repo, RepoAttributes } from '../shared/models/repo';
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

  // Order elements
  orders: {
    name: string,
    value: string
  }[] = [
    { name: 'Name', value: 'name' },
    { name: 'Created at', value: 'created' }
  ];
  orderBy: string = this.orders[0].value;

  // Repos
  repoName: string;
  allRepo: Repo;
  selectedRepository: Repo;
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
    this.allRepo = new Repo();
    this.allRepo.id = "all";
    this.allRepo.attributes = new RepoAttributes();
    this.allRepo.attributes.name = 'All';
    this.route.params.forEach((params: Params) => {
      this.repoName = params['repo'] ? params['repo'] : undefined;
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
      if(this.repoName) {
        this.selectedRepository = repos.filter(r => r.id == this.repoName)[0];
      } else {
        this.selectedRepository = this.allRepo;
      }
      this.repositories = repos;
      this.repositories.splice(0, 0, this.allRepo)
    });
  }

  goToRepo(repo: string) {
    this.router.navigate(['/charts', repo === 'all' ? '' : repo], {replaceUrl:true});
  }

  changeOrderBy(orderByValue: string) {
    this.orderBy = orderByValue;
    this.orderedCharts = this.orderCharts(this.orderedCharts);
  }

  searchChange(e) {
    let searchValue = e.target.value;
    if (!searchValue) {
      return this.orderedCharts = this.orderCharts(this.charts);
    }
    this.loading = true;
    this.chartsService.searchCharts(searchValue, this.repoName).subscribe(charts => {
      this.loading = false;
      this.orderedCharts = this.orderCharts(charts);
    });
  }

  // Sort charts
  orderCharts(charts): Chart[] {
    switch(this.orderBy) {
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
