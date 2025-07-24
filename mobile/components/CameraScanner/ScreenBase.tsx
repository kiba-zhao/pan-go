import type {ReactNode} from 'react';
import {use} from 'react';
import type {CameraDevice} from 'react-native-vision-camera';
import {useCameraDevice, useCameraPermission} from 'react-native-vision-camera';
import {ScreenLayout} from '../Common/ScreenBase';
import Text from '../Common/Text';

export type CameraScannerScreenProps = {
  children: (device: CameraDevice) => ReactNode;
};
export const CameraScannerScreen = ({children}: CameraScannerScreenProps) => {
  const device = useCameraDevice('back');
  const {hasPermission, requestPermission} = useCameraPermission();

  if (!hasPermission) {
    const hasPermission_ = use(requestPermission());
    if (!hasPermission_) return <PermissionScreen />;
  }
  if (!device) return <NoCameraDeviceErrorScreen />;
  return children(device);
};

export const NoCameraDeviceErrorScreen = () => {
  return (
    <ScreenLayout>
      <Text color="textPrimary">No Camera Screen</Text>
    </ScreenLayout>
  );
};

export const PermissionScreen = () => {
  return (
    <ScreenLayout>
      <Text color="textPrimary">Permission Screen</Text>
    </ScreenLayout>
  );
};
