import IonIcon from '@react-native-vector-icons/ionicons';
import type {ComponentProps} from 'react';
import {Text} from 'react-native';
import type {ScreenProps, TabBarIconType} from '../App/App.navigation';
import {ScreenLayout} from '../Template/Screen.layout';
import {DeviceName} from './Device.constants';

export const DeviceScreenName = DeviceName;

const DeviceScreen = () => {
  return (
    <ScreenLayout>
      <Text>Device Screen</Text>
    </ScreenLayout>
  );
};
export default DeviceScreen;

type DeviceScreenIconProps = ComponentProps<TabBarIconType>;
const DeviceScreenIcon = ({focused, size, color}: DeviceScreenIconProps) => (
  <IonIcon
    name={focused ? 'radio-sharp' : 'radio-outline'}
    size={size}
    color={color}
  />
);

export const DeviceScreenOptions: ScreenProps['options'] = {
  title: DeviceName,
  tabBarIcon: DeviceScreenIcon,
};
