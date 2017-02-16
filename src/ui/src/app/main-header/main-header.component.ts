import { Component, OnInit, Input } from '@angular/core';
import { Router, NavigationExtras } from '@angular/router';

@Component({
  selector: 'app-main-header',
  templateUrl: './main-header.component.html',
  styleUrls: ['./main-header.component.scss']
})
export class MainHeaderComponent implements OnInit {
  @Input() totalChartsNumber: number
  // Store the router
  constructor(private router: Router) { }
  ngOnInit() { }

  searchCharts(input: HTMLInputElement): void {
    // Empty query
    if(input.value === ""){
      this.router.navigate(["/"]);
      return;
    }

    let navigationExtras: NavigationExtras = {
      queryParams: { 'q': input.value }
    };
    this.router.navigate(["/charts/search"], navigationExtras)
  }
}
