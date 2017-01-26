import { Component, OnInit } from '@angular/core';
import { ChartsService } from '../charts.service';
import { Chart } from '../chart';

@Component({
  selector: 'app-charts',
  templateUrl: './charts.component.html',
  styleUrls: ['./charts.component.scss']
})
export class ChartsComponent implements OnInit {
  charts: Chart[] = [];
  orderedCharts: Chart[] = [];
  constructor(private chartsService: ChartsService) { }

  // Default filters
  filters = {
    orderBy: 'title',
    repositoryType: 'all'
  }

  ngOnInit() {
		this.loadCharts();
  }

  loadCharts(): void {
		this.chartsService.getCharts().subscribe(charts => {
      this.charts = charts;
      this.orderedCharts = this.orderCharts(this.filterCharts(charts));
    });
  }

  // Update a filter
  onChangeFilter(filter): void {
    if (this.filters[filter.type] !== filter.value) {
      this.filters[filter.type] = filter.value;
      this.orderedCharts = this.orderCharts(this.filterCharts(this.charts));
    }
  }

  // Filter and sort the charts
  filterCharts(charts): Chart[] {
    let filteredCharts = charts;
    // Check if we need to apply a filter
    if (this.filters.repositoryType !== 'all') {
      filteredCharts = charts.filter(c =>
        c.attributes.repo === this.filters.repositoryType);
    }

    return filteredCharts;
  }

  // Sort charts
  orderCharts(charts): Chart[] {
    switch(this.filters.orderBy) {
      case 'repository': {
        return charts.sort((a, b) =>
          a.attributes.repo.localeCompare(b.attributes.repo));
      }
      default: {
        return charts.sort((a, b) =>
          a.attributes.name.localeCompare(b.attributes.name));
      }
    }
  }
}
