import { Component, OnInit, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-charts-filters',
  templateUrl: './charts-filters.component.html',
  styleUrls: ['./charts-filters.component.scss']
})
export class ChartsFiltersComponent implements OnInit {
  // Order elements
  orderElements: {
    name: string,
    value: string
  }[] = [
    {
      name: 'Title',
      value: 'title'
    },
    {
      name: 'Repository',
      value: 'repository'
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
  repositoryType: string = this.repositoryElements[0].value;

  @Output() onChange = new EventEmitter();

  constructor() { }

  ngOnInit() {
  }

  // Emit the changes of the filters
  onChangeFilter(type, value) {
    this.onChange.emit({ type, value });
  }
}
