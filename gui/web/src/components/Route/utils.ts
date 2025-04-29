export function generateCreatePath(path: string): string {
  return `${path}/create`;
}

export function generateEditPath<T extends any>(path: string, id?: T): string {
  return `${path}/${id || ":id"}`;
}

export function generateShowPath<T extends any>(path: string, id?: T): string {
  return `${path}/${id || ":id"}/show`;
}
