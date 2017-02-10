import { Routes, RouterModule } from '@angular/router';
import { ModuleWithProviders }  from '@angular/core';

import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import { ChartIndexComponent } from './chart-index/chart-index.component';
import { ChartDetailsComponent } from './chart-details/chart-details.component';
import { ChartSearchComponent } from './chart-search/chart-search.component';
import { ChartsComponent } from './charts/charts.component';

const appRoutes: Routes = [
  {
    path: '',
    component: ChartIndexComponent,
    data: {
      meta: {
        title: 'Monocular'
      }
    }
  },
  {
    path: 'charts',
    component: ChartsComponent,
    data: {
      meta: {
        title: 'All charts'
      }
    }
  },
  {
    path: 'charts/search',
    component: ChartSearchComponent
  },
  {
    path: 'charts/:repo',
    component: ChartsComponent,
    data: {
      meta: {
        title: 'All charts by repository'
      }
    }
  },
  {
    path: 'charts/:repo/:chartName',
    component: ChartDetailsComponent
  },
  {
    path: 'charts/:repo/:chartName/:version',
    component: ChartDetailsComponent
  },
  {
    path: '**',
    component: PageNotFoundComponent,
    data: {
      meta: {
        title: 'Not Found'
      }
    }
  }
];

export const appRoutingProviders: any[] = [

];

export const routing: ModuleWithProviders = RouterModule.forRoot(appRoutes);
