import { ChartVersionAttributes } from "./chart-version"
export class Chart {
  id: String;
  type: String;
  links: String[];
  attributes: ChartAttributes;
  relationships: ChartRelationships;
}

export class ChartAttributes {
  description: String;
  name: String;
  icon: String;
  repo: String;
  home: String;
  sources: String[];
}

class ChartRelationships {
  latestChartVersion: ChartVersionRelationship;
}

class ChartVersionRelationship {
  data: ChartVersionAttributes
  links: {
    self: String
  }
}
