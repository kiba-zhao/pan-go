import {Device} from '@pango/data';

export type EditorParam = Pick<Device, 'id'>;
export type CreatorParam = Pick<Device, 'name' | 'peerId'>;

export const DeviceCreatorScreenName = 'device.creator';
export const DeviceEditorScreenName = 'device.editor';
export const DeviceSearchScreenName = 'device.search';
