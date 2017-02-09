import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { Router, NavigationExtras } from '@angular/router';

@Component({
  selector: 'app-header-bar',
  templateUrl: './header-bar.component.html',
  styleUrls: ['./header-bar.component.scss'],
  encapsulation: ViewEncapsulation.None,
  inputs: ['showSearch']
})
export class HeaderBarComponent implements OnInit {
  // Show search form by default
  public showSearch:boolean = true;

  constructor(private router: Router) { }
  ngOnInit() { }

  searchCharts(input: HTMLInputElement): void {
    // Empty query
    if(input.value === ''){
      this.router.navigate(['/']);
    } else {
      let navigationExtras: NavigationExtras = {
        queryParams: { 'q': input.value }
      };
      this.router.navigate(["/charts/search"], navigationExtras);
    }
  }
}
