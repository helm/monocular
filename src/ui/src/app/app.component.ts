import { Angulartics2GoogleAnalytics } from 'angulartics2';
import { Component } from '@angular/core';
// import { MetaService } from 'ng2-meta';
import { MenuService } from './shared/services/menu.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  providers: [MenuService]
})
export class AppComponent {
  // Show the global menu
  public showMenu: boolean = false;

  // constructor(private metaService: MetaService) {}
  constructor(
    angulartics2GoogleAnalytics: Angulartics2GoogleAnalytics,
    private menuService: MenuService
  ) {
    console.log('Initialize');
    menuService.menuOpen$.subscribe(() => {
      console.log('received!');
      this.showMenu = !this.showMenu;
    });
  }
}
