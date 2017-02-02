import { ChartVersionAttributes } from "./chart-version"
export class Chart {
  id: string;
  type: string;
  links: string[];
  attributes: ChartAttributes;
  relationships: ChartRelationships;
}


export class ChartAttributes {
  description: string;
  name: string;
  icon: string;
  repo: string;
  home: string;
  sources: string[];
}

class ChartRelationships {
  latestChartVersion: ChartVersionRelationship;
}

class ChartVersionRelationship {
  data: ChartVersionAttributes
  links: {
    self: string
  }
}
