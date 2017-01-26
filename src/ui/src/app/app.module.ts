import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { routing, appRoutingProviders } from './app.routing';

/* Material library */
import { MaterialModule } from '@angular/material';

/* Pipes */
import { TruncatePipe } from './pipes/truncate.pipe';

/* Services */
import { ChartsService } from './charts.service';

/* Components */
import { AppComponent } from './app.component';
import { ChartIndexComponent } from './chart-index/chart-index.component';
import { ChartListComponent } from './chart-list/chart-list.component';
import { ChartItemComponent } from './chart-item/chart-item.component';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import { ChartDetailsComponent } from './chart-details/chart-details.component';
import { ChartSearchInputComponent } from './header-bar/chart-search-input/chart-search-input.component';
import { ChartSearchComponent } from './chart-search/chart-search.component';
import { HeaderBarComponent } from './header-bar/header-bar.component';
import { ChartDetailsUsageComponent } from './chart-details/chart-details-usage/chart-details-usage.component';
import { ChartDetailsReadmeComponent } from './chart-details/chart-details-readme/chart-details-readme.component';
import { PanelComponent } from './panel/panel.component';
import { MainHeaderComponent } from './main-header/main-header.component';
import { TrustedMaintainersComponent } from './trusted-maintainers/trusted-maintainers.component';
import { FooterComponent } from './footer/footer.component';
import { FooterListComponent } from './footer-list/footer-list.component';
import { ChartDetailsInfoComponent } from './chart-details/chart-details-info/chart-details-info.component';
import { ChartListItemComponent } from './chart-list/chart-list-item/chart-list-item.component';
import { ChartsComponent } from './charts/charts.component';
import { ChartsFiltersComponent } from './charts-filters/charts-filters.component';

require('hammerjs');

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
    HeaderBarComponent,
    ChartDetailsUsageComponent,
    ChartDetailsReadmeComponent,
    PanelComponent,
    MainHeaderComponent,
    TruncatePipe,
    TrustedMaintainersComponent,
    FooterComponent,
    FooterListComponent,
    ChartDetailsInfoComponent,
    ChartListItemComponent,
    ChartsComponent,
    ChartsFiltersComponent
  ],
  imports: [
    MaterialModule.forRoot(),
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
