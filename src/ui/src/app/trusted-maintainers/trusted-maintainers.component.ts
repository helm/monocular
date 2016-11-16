import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-trusted-maintainers',
  templateUrl: './trusted-maintainers.component.html',
  styleUrls: ['./trusted-maintainers.component.scss']
})
export class TrustedMaintainersComponent {
  // TODO: Fill this from the API?
  public maintainers:{ name: string, logo: string }[] = [
    {
      name: 'Deis',
      logo: '/assets/images/maintainers/deis.png'
    },
    {
      name: 'Bitnami',
      logo: '/assets/images/maintainers/bitnami.png'
    }
  ];
}
