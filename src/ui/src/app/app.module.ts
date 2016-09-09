import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { routing, appRoutingProviders } from './app.routing';

/* Components */
import { AppComponent } from './app.component';
import { ChartIndexComponent } from './chart-index/chart-index.component';
import { ChartListComponent } from './chart-list/chart-list.component';
import { ChartItemComponent } from './chart-list/chart-item/chart-item.component';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import { ChartDetailsComponent } from './chart-details/chart-details.component';

/* Services */
import { ChartsService } from './charts.service';
import { ChartSearchInputComponent } from './page-header-bar/chart-search-input/chart-search-input.component';
import { ChartSearchComponent } from './chart-search/chart-search.component';
import { PageHeaderBarComponent } from './page-header-bar/page-header-bar.component';
import { ChartDetailsHeaderComponent } from './chart-details/chart-details-header/chart-details-header.component';
import { ChartDetailsUsageComponent } from './chart-details/chart-details-usage/chart-details-usage.component';
import { ChartDetailsReadmeComponent } from './chart-details/chart-details-readme/chart-details-readme.component'

require('bootstrap-loader');

@NgModule({
  declarations: [
    AppComponent,
    ChartIndexComponent,
    ChartListComponent,
    ChartItemComponent,
    PageNotFoundComponent,
    ChartDetailsComponent,
    ChartSearchInputComponent,
    ChartSearchComponent,
    PageHeaderBarComponent,
    ChartDetailsHeaderComponent,
    ChartDetailsUsageComponent,
    ChartDetailsReadmeComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
		routing
  ],
  providers: [
    appRoutingProviders,
    ChartsService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
