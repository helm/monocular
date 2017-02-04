import { Routes, RouterModule } from '@angular/router';
import { ModuleWithProviders }  from '@angular/core';

import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import { ChartIndexComponent } from './chart-index/chart-index.component';
import { ChartDetailsComponent } from './chart-details/chart-details.component';
import { ChartSearchComponent } from './chart-search/chart-search.component';
import { ChartsComponent } from './charts/charts.component';

const appRoutes: Routes = [
  { path: '', component: ChartIndexComponent },
  { path: 'charts', component: ChartsComponent },
  { path: 'charts/search', component: ChartSearchComponent },
  { path: 'charts/:repo/:chartName', component: ChartDetailsComponent },
  { path: 'charts/:repo/:chartName/:version', component: ChartDetailsComponent },
  { path: '**', component: PageNotFoundComponent }
];

export const appRoutingProviders: any[] = [

];

export const routing: ModuleWithProviders = RouterModule.forRoot(appRoutes);
