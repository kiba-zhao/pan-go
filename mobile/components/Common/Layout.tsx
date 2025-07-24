import type {ViewStyle} from 'react-native';
import {StyleSheet, View} from 'react-native';
import {withStyle} from './StyleSheet';

const styles = StyleSheet.create({
  flexLayout: {
    flex: 1,
  },
  flexRowLayout: {
    flex: 1,
    flexDirection: 'row',
  },
  rowLayout: {
    flexDirection: 'row',
  },
});

export const FlexLayout = withStyle(View, styles.flexLayout as ViewStyle);
export const FlexRowLayout = withStyle(View, styles.flexRowLayout as ViewStyle);
export const RowLayout = withStyle(View, styles.rowLayout as ViewStyle);
