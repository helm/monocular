import { Component } from '@angular/core';
import { MONOCULAR_VERSION } from '../../version';

@Component({
  selector: 'app-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss']
})
export class FooterComponent {
  monocularVersion: string = MONOCULAR_VERSION;
}
