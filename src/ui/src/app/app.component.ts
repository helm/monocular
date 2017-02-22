import { Angulartics2GoogleAnalytics } from 'angulartics2';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { MenuService } from './shared/services/menu.service';
import { ChartsService } from './shared/services/charts.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  providers: [MenuService, ChartsService]
})
export class AppComponent {
  // Show the global menu
  public showMenu: boolean = false;

  constructor(
    angulartics2GoogleAnalytics: Angulartics2GoogleAnalytics,
    private menuService: MenuService,
    private router: Router
  ) {
    menuService.menuOpen$.subscribe(show => {
      this.showMenu = show;
    });

    // Hide menu when user changes the route
    router.events.subscribe(() => {
      menuService.hideMenu();
    });
  }
}
