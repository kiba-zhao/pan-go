import {PropsWithChildren} from 'react';
import {StyleSheet, View} from 'react-native';

type ScreenLayoutProps = PropsWithChildren<{}>;
export const ScreenLayout = ({children}: ScreenLayoutProps) => (
  <View>{children}</View>
);

const styles = StyleSheet.create({});
