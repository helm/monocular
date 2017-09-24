import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { Router, NavigationExtras } from '@angular/router';
import { ConfigService } from '../shared/services/config.service';
import { MenuService } from '../shared/services/menu.service';
import { DomSanitizer } from '@angular/platform-browser';
import { MdIconRegistry } from '@angular/material';
import { MdSnackBar } from '@angular/material';
import { CookieService } from 'ngx-cookie';

@Component({
  selector: 'app-header-bar',
  templateUrl: './header-bar.component.html',
  styleUrls: ['./header-bar.component.scss'],
  encapsulation: ViewEncapsulation.None,
  viewProviders: [MdIconRegistry],
  inputs: ['showSearch', 'transparent']
})
export class HeaderBarComponent implements OnInit {
  // public loggedIn
  public loggedIn: boolean;
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
    public config: ConfigService,
    private menuService: MenuService,
    private mdIconRegistry: MdIconRegistry,
    private sanitizer: DomSanitizer,
    private cookieService: CookieService,
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
    this.mdIconRegistry.addSvgIcon(
      'search',
      this.sanitizer.bypassSecurityTrustResourceUrl('/assets/icons/search.svg')
    );
    this.appName = this.config.appName;

    let userClaims = this.cookieService.get("ka_claims")
    if (userClaims) {
      this.loggedIn = true;
    }
  }

  searchCharts(input: HTMLInputElement): void {
    // Empty query
    if (input.value === '') {
      this.router.navigate(['/charts']);
    } else {
      let navigationExtras: NavigationExtras = {
        queryParams: { q: input.value }
      };
      this.router.navigate(['/charts'], navigationExtras);
    }
  }

  openMenu() {
    // Open the menu
    this.openedMenu = !this.openedMenu;
    this.menuService.toggleMenu();
  }
}
