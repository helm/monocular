import { ChartAttributes } from "./chart"
export class ChartVersion {
  id: String;
  type: String;
  attributes: ChartVersionAttributes;
  relationships: ChartVersionRelationships;
}

export class ChartVersionAttributes {
  created: Date
  digest: String
  version: String
  urls: String[]
}

class ChartVersionRelationships {
  chart: ChartVersionChart;
}

class ChartVersionChart {
  data: ChartAttributes
  links: {
    self: String
  }
}
