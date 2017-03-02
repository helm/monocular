import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core';
import { Repo } from '../shared/models/repo';

@Component({
  selector: 'app-charts-filters',
  templateUrl: './charts-filters.component.html',
  styleUrls: ['./charts-filters.component.scss']
})
export class ChartsFiltersComponent implements OnInit {
  @Input() currentRepo: Repo;
  @Input() repositories: Repo[];
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
      name: 'Created at',
      value: 'created'
    }
  ];

  allRepo: Repo
  // Order of the elements
  orderBy: string = this.orderElements[0].value;

  constructor() {}

  ngOnInit() {
    this.allRepo = new Repo(); this.allRepo.id = "all";
    this.currentRepo  = this.allRepo;
  }

  // Emit the changes of the filters
  onChangeFilter(type, value) {
    this.onChange.emit({ type, value });
  }

  firstUpper(str: string): string {
    return str.substring(0,1).toUpperCase() + str.substring(1).toLowerCase();
  }
}
