import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-charts-filters',
  templateUrl: './charts-filters.component.html',
  styleUrls: ['./charts-filters.component.scss']
})
export class ChartsFiltersComponent implements OnInit {
  @Input() currentRepo: string
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
  // Repository Types
  repositoryElements: {
    name: string,
    value: string
  }[] = [
    {
      name: 'All',
      value: 'all'
    },
    {
      name: 'Stable',
      value: 'stable'
    },
    {
      name: 'Incubator',
      value: 'incubator'
    }
  ]

  // Order of the elements
  orderBy: string = this.orderElements[0].value;
  repositoryType: string

  constructor() {}

  ngOnInit() {
    this.repositoryType = this.currentRepo || this.repositoryElements[0].value;
  }

  // Emit the changes of the filters
  onChangeFilter(type, value) {
    this.onChange.emit({ type, value });
  }
}
