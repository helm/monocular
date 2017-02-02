export class Chart {
  id: String;
  type: String;
  links: String[];
  attributes: ChartAttributes;
  relationships: ChartRelationships;
}

class ChartAttributes {
  description: String;
  name: String;
  icon: String;
  repo: String;
  home: String;
  sources: String[];
}

class ChartRelationships {
  latestChartVersion: ChartVersion;
}

class ChartVersion {
  data: {
    created: String,
    digest: String,
    version: String,
    urls: String[]
  }
  links: {
    self: String
  }
}
