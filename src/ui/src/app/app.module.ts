import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { Angulartics2Module, Angulartics2GoogleAnalytics } from 'angulartics2';
import { ClipboardModule } from 'ngx-clipboard';
import { MetaModule, MetaConfig } from 'ng2-meta';
import { routing, appRoutingProviders } from './app.routing';

/* Material library */
import { MaterialModule } from '@angular/material';

/* Pipes */
import { TruncatePipe } from './shared/pipes/truncate.pipe';

/* Services */
import { ChartsService } from './shared/services/charts.service';
import { ConfigService } from './shared/services/config.service';

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
import { ChartDetailsVersionsComponent } from './chart-details/chart-details-versions/chart-details-versions.component';
import { ChartsComponent } from './charts/charts.component';
import { ChartsFiltersComponent } from './charts-filters/charts-filters.component';
import { LoaderComponent } from './loader/loader.component';

require('hammerjs');

const metaConfig: MetaConfig = {
  //Append a title suffix such as a site name to all titles
  useTitleSuffix: true,
  defaults: {
    title: 'Monocular',
    titleSuffix: ' | Monocular',
    description: 'Find the Helm chart you need or share your own Kubernetes apps'
  }
};

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
    ChartDetailsVersionsComponent,
    ChartDetailsReadmeComponent,
    PanelComponent,
    MainHeaderComponent,
    TruncatePipe,
    TrustedMaintainersComponent,
    FooterComponent,
    FooterListComponent,
    ChartDetailsInfoComponent,
    ChartsComponent,
    ChartsFiltersComponent,
    LoaderComponent
  ],
  imports: [
    MaterialModule.forRoot(),
    BrowserModule,
    FormsModule,
    HttpModule,
		routing,
    Angulartics2Module.forRoot([Angulartics2GoogleAnalytics]),
    ClipboardModule,
    MetaModule.forRoot(metaConfig)
  ],
  providers: [
    appRoutingProviders,
    ChartsService,
    ConfigService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
