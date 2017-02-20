import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { Router, NavigationExtras } from '@angular/router';
import { ConfigService } from '../shared/services/config.service';
import { MenuService } from '../shared/services/menu.service';

@Component({
  selector: 'app-header-bar',
  templateUrl: './header-bar.component.html',
  styleUrls: ['./header-bar.component.scss'],
  encapsulation: ViewEncapsulation.None,
  inputs: ['showSearch', 'transparent']
})
export class HeaderBarComponent implements OnInit {
  // Show search form by default
  public showSearch: boolean = true;
  // Set the background as transparent
  public transparent: boolean = false;

  appName: string
  constructor(
    private router: Router,
    private config: ConfigService,
    private menuService: MenuService
  ) {
    this.appName = config.appName
  }
  ngOnInit() { }

  searchCharts(input: HTMLInputElement): void {
    // Empty query
    if(input.value === ''){
      this.router.navigate(['/']);
    } else {
      let navigationExtras: NavigationExtras = {
        queryParams: { 'q': input.value }
      };
      this.router.navigate(['/charts/search'], navigationExtras);
    }
  }

  openMenu() {
    // Open the menu
    console.log('test');
    this.menuService.toggleMenu();
  }
}
