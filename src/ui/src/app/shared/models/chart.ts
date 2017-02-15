import { ChartVersionAttributes } from "./chart-version"
import { Repo } from "./repo"
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
  repo: Repo;
  home: string;
  sources: string[];
  keywords: string[];
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
