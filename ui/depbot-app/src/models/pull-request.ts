
type Author = {
  login: string
}

export type PullRequest = {
  title: string
  number: number
  url: string
  state: string
  createdAt: string
  author: Author
}

export class PullRequestState {
  static OPEN = 'open';
  static MERGED = 'merged';
  static CLOSED = 'closed';
}

