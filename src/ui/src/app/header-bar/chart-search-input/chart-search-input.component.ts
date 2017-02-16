import { Component, OnInit } from '@angular/core';
import { Router, NavigationExtras } from '@angular/router';

@Component({
  selector: 'app-chart-search-input',
  templateUrl: './chart-search-input.component.html',
  styleUrls: ['./chart-search-input.component.scss']
})
export class ChartSearchInputComponent implements OnInit {
  constructor(private router: Router) { }
  ngOnInit() {
  }

  searchCharts(input: HTMLInputElement): void {
    // Empty query
    if(input.value === ''){
      this.router.navigate(['/'])
    } else {
      let navigationExtras: NavigationExtras = {
        queryParams: { 'q': input.value }
      };
      input.value = '';
      this.router.navigate(['/charts/search'], navigationExtras)
    }
  }

}
