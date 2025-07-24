import {Device} from '@pango/data';

export type EditorParam = Pick<Device, 'id'>;
export type CreatorParam = Pick<Device, 'name' | 'peerId'>;
