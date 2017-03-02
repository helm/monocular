import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-charts-filters',
  templateUrl: './charts-filters.component.html',
  styleUrls: ['./charts-filters.component.scss']
})
export class ChartsFiltersComponent implements OnInit {
  @Input() currentRepo: string
  @Input() repositories: string[]
  @Output() onChange = new EventEmitter();

  // Order elements
  orderElements: {
    name: string,
    value: string
  }[] = [
    {
      name: 'Name',
      value: 'name'
    },
    {
      name: 'Creation date',
      value: 'created'
    }
  ];
  // Order of the elements
  orderBy: string = this.orderElements[0].value;
  repositoryType: string

  constructor() {}

  ngOnInit() {
    this.repositoryType = this.currentRepo || this.repositoryElements[0];
  }

  // Emit the changes of the filters
  onChangeFilter(type, value) {
    this.onChange.emit({ type, value });
  }

  get repositoryElements(): string[] {
    return ["all"].concat(this.repositories);
  }
}
