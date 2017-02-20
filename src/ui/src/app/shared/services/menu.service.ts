import { Injectable } from '@angular/core';
import { Subject }    from 'rxjs/Subject';

@Injectable()
export class MenuService {
  // Emitter
  private menuOpenSource = new Subject<string>();
  private open: boolean = false;
  // Observable boolean streams
  public menuOpen$ = this.menuOpenSource.asObservable();


  // Service message commands
  toggleMenu() {
    this.open = !this.open;
    console.log(`Emit: ${this.open}`);
    this.menuOpenSource.next('test');
  }
}
