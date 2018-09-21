import { Component, OnInit } from '@angular/core';
import { ChartsService } from '../shared/services/charts.service';
import { Chart } from '../shared/models/chart';
import { RepoAttributes } from '../shared/models/repo';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { SeoService } from '../shared/services/seo.service';
import { ConfigService } from '../shared/services/config.service';
import { MatIconRegistry } from '@angular/material';
import { DomSanitizer } from '@angular/platform-browser';

@Component({
  selector: 'app-charts',
  templateUrl: './charts.component.html',
  styleUrls: ['./charts.component.scss'],
  viewProviders: [MatIconRegistry]
})
export class ChartsComponent implements OnInit {
  charts: Chart[] = [];
  orderedCharts: Chart[] = [];
  loading: boolean = true;
  searchTerm: string;
  searchTimeout: any;
  filtersOpen: boolean = false;

  // Default filters
  filters = [
    {
      title: 'Repository',
      onSelect: i => this.onSelectRepo(i),
      items: [{ title: 'All', value: 'all', selected: true }]
    },
    {
      title: 'Order By',
      onSelect: i => this.onSelectOrderBy(i),
      items: [
        { title: 'Name', value: 'name', selected: true },
        { title: 'Created At', value: 'created', selected: false }
      ]
    }
  ];

  // Order elements
  orderBy: string = 'name';

  // Repos
  repoName: string;

  constructor(
    private chartsService: ChartsService,
    private route: ActivatedRoute,
    private router: Router,
    private config: ConfigService,
    private seo: SeoService,
    private mdIconRegistry: MatIconRegistry,
    private sanitizer: DomSanitizer
  ) {}

  ngOnInit() {
    this.mdIconRegistry.addSvgIcon(
      'search',
      this.sanitizer.bypassSecurityTrustResourceUrl(`/assets/icons/search.svg`)
    );
    this.mdIconRegistry.addSvgIcon(
      'close',
      this.sanitizer.bypassSecurityTrustResourceUrl(`/assets/icons/close.svg`)
    );
    this.mdIconRegistry.addSvgIcon(
      'menu',
      this.sanitizer.bypassSecurityTrustResourceUrl(`/assets/icons/menu.svg`)
    );
    this.route.queryParams.forEach((params: Params) => {
      this.searchTerm = params['q'] ? params['q'] : undefined;
      if (this.searchTerm) {
        this.searchCharts();
      }
    });
    this.route.params.forEach((params: Params) => {
      this.repoName = params['repo'] ? params['repo'] : undefined;
      this.updateMetaTags();
      this.loadCharts();
    });
  }

  loadCharts(): void {
    this.chartsService.getCharts().subscribe(allCharts => {
      this.loading = false;
      this.charts = allCharts.filter(c => !this.repoName || c.attributes.repo.name === this.repoName);
      if (!this.searchTerm) {
        this.orderedCharts = this.orderCharts(this.charts);
      }
      this.setReposFromCharts(allCharts);
    });
  }
  
  // This takes a list of charts, extracts the unique set of repositories the
  // charts are from and sets the Repositories filter with that list. We also
  // add an 'all' repository filter at the top.
  setReposFromCharts(charts: Chart[]): void {
    let repoMap = new Map<string, RepoAttributes>();
    repoMap['all'] = { name: 'All' };
    repoMap = charts.reduce((repos, chart) => {
      repos[chart.attributes.repo.name] = chart.attributes.repo;
      return repos;
    }, repoMap);
    
    this.filters[0].items = Object.keys(repoMap).map(k => {
      const r = repoMap[k];
      return {
        title: r.name,
        value: k,
        selected: this.repoName ? k === this.repoName : k == 'all'
      }
    });
  }

  onSelectRepo(index) {
    this.repoName = this.filters[0].items[index].value;
    this.filters[0].items = this.filters[0].items.map(r => {
      r.selected = r.value == this.repoName;
      return r;
    });
    this.router.navigate(
      ['/charts', this.repoName === 'all' ? '' : this.repoName],
      { replaceUrl: true }
    );
  }

  onSelectOrderBy(index) {
    this.orderBy = this.filters[1].items[index].value;
    this.filters[1].items = this.filters[1].items.map(o => {
      o.selected = o.value == this.orderBy;
      return o;
    });
    this.orderedCharts = this.orderCharts(this.orderedCharts);
  }

  searchChange(e) {
    this.searchTerm = e.target.value;
    clearTimeout(this.searchTimeout);
    if (!this.searchTerm) {
      return (this.orderedCharts = this.orderCharts(this.charts));
    }
    this.searchTimeout = setTimeout(() => this.searchCharts(), 1000);
  }

  searchCharts() {
    if (!this.searchTerm) {
      return false;
    }
    this.loading = true;
    this.chartsService
      .searchCharts(this.searchTerm, this.repoName)
      .subscribe(charts => {
        this.loading = false;
        this.orderedCharts = this.orderCharts(charts);
      });
  }

  // Sort charts
  orderCharts(charts): Chart[] {
    switch (this.orderBy) {
      case 'created': {
        return charts.sort(this.sortByCreated).reverse();
      }
      default: {
        return charts.sort((a, b) =>
          a.attributes.name.localeCompare(b.attributes.name)
        );
      }
    }
  }

  sortByCreated(a: Chart, b: Chart) {
    let aVersion = a.relationships.latestChartVersion.data;
    let bVersion = b.relationships.latestChartVersion.data;
    if (aVersion.created < bVersion.created) {
      return -1;
    } else if (aVersion.created > bVersion.created) {
      return 1;
    }
    return 0;
  }

  updateMetaTags(): void {
    if (this.repoName) {
      this.seo.setMetaTags('repoCharts', {
        repo: this.capitalize(this.repoName)
      });
    } else {
      this.seo.setMetaTags('charts');
    }
  }

  capitalize(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
  }
}
