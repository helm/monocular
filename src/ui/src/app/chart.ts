export class Chart {
  id: String
  type: String
  links: String[]
  attributes: ChartAttributes
}

class ChartAttributes {
  description: String
  name: String
  repo: String
  home: String
}
