import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-panel',
  templateUrl: './panel.component.html',
  styleUrls: ['./panel.component.scss'],
  inputs: ['title']
})
export class PanelComponent {
  // Title of the panel
  public title:string = '';
}
