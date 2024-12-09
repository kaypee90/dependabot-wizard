
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
