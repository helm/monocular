export class Repo {
  id: string;
  type: string;
  attributes: RepoAttributes
}

export class RepoAttributes {
  name: string = '';
  url: string = '';
}
