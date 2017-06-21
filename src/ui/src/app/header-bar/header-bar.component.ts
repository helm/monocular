import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { Router, NavigationExtras } from '@angular/router';
import { ConfigService } from '../shared/services/config.service';
import { MenuService } from '../shared/services/menu.service';
import { DomSanitizer } from '@angular/platform-browser';
import { MdIconRegistry } from '@angular/material';
import { MdSnackBar } from '@angular/material';

@Component({
  selector: 'app-header-bar',
  templateUrl: './header-bar.component.html',
  styleUrls: ['./header-bar.component.scss'],
  encapsulation: ViewEncapsulation.None,
  viewProviders: [MdIconRegistry],
  inputs: ['showSearch', 'transparent']
})
export class HeaderBarComponent implements OnInit {
  // Show search form by default
  public showSearch: boolean = true;
  // Set the background as transparent
  public transparent: boolean = false;
  // Check if  the menu is opened
  public openedMenu: boolean = false;
  // Config

  appName: string;
  constructor(
    private router: Router,
    private config: ConfigService,
    private menuService: MenuService,
    private mdIconRegistry: MdIconRegistry,
    private sanitizer: DomSanitizer
  ) {}

  ngOnInit() {
    // Set the icon
    this.mdIconRegistry.addSvgIcon(
      'menu',
      this.sanitizer.bypassSecurityTrustResourceUrl('/assets/icons/menu.svg')
    );
    this.mdIconRegistry.addSvgIcon(
      'close',
      this.sanitizer.bypassSecurityTrustResourceUrl('/assets/icons/close.svg')
    );
    this.appName = this.config.appName;
  }

  searchCharts(input: HTMLInputElement): void {
    // Empty query
    if (input.value === '') {
      this.router.navigate(['/']);
    } else {
      let navigationExtras: NavigationExtras = {
        queryParams: { q: input.value }
      };
      this.router.navigate(['/charts/search'], navigationExtras);
    }
  }

  openMenu() {
    // Open the menu
    this.openedMenu = !this.openedMenu;
    this.menuService.toggleMenu();
  }
}
