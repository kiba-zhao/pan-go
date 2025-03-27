

export const extractPathSeq = (path: string): string => {
  if (path.indexOf("/") >= 0) {
    return "/";
  }
  if (path.indexOf("\\") > 0) {
    return "\\";
  }
  const ROOT_DRIVE_REGEX = /^[a-z|A-Z]+(\:+)(\S*)$/gi;
  if (ROOT_DRIVE_REGEX.test(path)) {
    return "\\";
  }
  throw new Error("Unknown path separator");
};


export const extractRoot = (path: string): string => {
  if (path.startsWith("/")) {
    return "/";
  }
  if (path.startsWith("\\\\")) {
    return "\\\\";
  }

  const ROOT_DRIVE_REGEX = /^[a-z|A-Z]+(\:+)(\S*)$/gi;

  if (ROOT_DRIVE_REGEX.test(path)) {
    const idx = path.indexOf(":");
    return path.slice(0, idx+1);
  }
  throw new Error("Unknown path root");
};

export type DirnameOpts = { seq?: string; root?: string };
export const dirname = (
  path: string,
  { seq, root }: DirnameOpts = {}
): string => {
  const _root = root || extractRoot(path);
  if (_root == path) {
    return path;
  }
  const _sep = seq || extractPathSeq(path);
  const lastIdx = path.lastIndexOf(_sep);
  if (lastIdx == 0) return _root;
  if (lastIdx < 0) throw new Error("Invalid path");
  let idx = lastIdx - 1;
  while (idx >= _root.length) {
    if (path[idx] != _sep) {
      break;
    }
    idx--;
  }
  return path.slice(0, idx + 1);
};

export const generateParents = (path: string, reverse = false): string[] => {
  const parents: string[] = [];
  const root = extractRoot(path);
  const append = reverse ? parents.unshift : parents.push;
  let filepath = path;
  while (filepath != root) {
    filepath = dirname(filepath, { root });
    append.call(parents, filepath);
  }
  return parents.length <= 0 ? [root] : parents;
};

export type BasenameOpts = { seq?: string; root?: string; dirname?: string };
export const basename = (path: string, opts: BasenameOpts = {}): string => {
  const _seq = opts.seq || extractPathSeq(path);
  const _dirname = opts.dirname || dirname(path, { seq: _seq, ...opts });
  let idx = _dirname.length;
  while (idx < path.length) {
    if (path[idx] !== _seq) break;
    idx++;
  }
  return path.slice(idx);
};
