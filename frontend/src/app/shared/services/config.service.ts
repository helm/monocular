import { Injectable } from '@angular/core';

@Injectable()
export class ConfigService {
  // Configurable options
  // They can be overriden using assets/js/overrides.js
  backendHostname: string
  appName: string
  aboutUrl: string
  // EO configurable options

  constructor() {
    var overrides: any = {}
    if (window["monocular"]) {
      overrides = window["monocular"]["overrides"] || {};
    }

    this.backendHostname = overrides.backendHostname || "/api";
    this.appName = overrides.appName || "Monocular";
    this.aboutUrl = overrides.aboutUrl || "https://github.com/helm/monocular/blob/master/docs/about.md";
  }
}
