export class Release {
  id: string;
  type: string;
  attributes: ReleaseAttributes
}

class ReleaseAttributes {
  chartName: string;
  chartVersion: string;
  name: string;
  namespace: string;
  status: string;
  updated: Date;
}
