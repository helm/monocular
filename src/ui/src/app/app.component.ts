import { Component } from '@angular/core';
import { MetaService } from 'ng2-meta';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  constructor(private metaService: MetaService) {}
}
