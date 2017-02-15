import { ChartVersionAttributes } from "./chart-version"
import { Repo } from "./repo"
import { Maintainer } from "./maintainer"
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
  maintainers: Maintainer[];
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
