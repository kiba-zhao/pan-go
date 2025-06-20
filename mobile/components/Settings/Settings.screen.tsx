import IonIcon from '@react-native-vector-icons/ionicons';
import {ComponentProps} from 'react';
import {Text} from 'react-native';
import type {ScreenPropsOptions, TabBarIconType} from '../App/App.navigation';
import {ScreenLayout} from '../Template/Screen.layout';
import {SettingsName} from './Settings.constants';

export const SettingsScreenName = SettingsName;

const SettingsScreen = () => {
  return (
    <ScreenLayout>
      <Text>Settings Screen</Text>
    </ScreenLayout>
  );
};
export default SettingsScreen;

type SettingsScreenIconProps = ComponentProps<TabBarIconType>;
const SettingsScreenIcon = ({
  focused,
  size,
  color,
}: SettingsScreenIconProps) => (
  <IonIcon
    name={focused ? 'settings-sharp' : 'settings-outline'}
    size={size}
    color={color}
  />
);

export const SettingsScreenOptions: ScreenPropsOptions = {
  title: SettingsName,
  tabBarIcon: SettingsScreenIcon,
};
