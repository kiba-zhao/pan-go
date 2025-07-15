import {ComponentProps} from 'react';
import type {ViewStyle} from 'react-native';
import {ScrollView, StyleSheet, View} from 'react-native';
import {scale, verticalScale} from './SizeMatters';
import {withStyle} from './StyleSheet';
import {useTheme} from './Theme';

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

type ScreenLayoutProps = ComponentProps<typeof ScrollView>;
export const ScreenLayout = ({
  children,
  style,
  ...props
}: ScreenLayoutProps) => {
  const {sizes} = useTheme();
  return (
    <ScrollView
      {...props}
      style={StyleSheet.compose({paddingHorizontal: scale(sizes.base)}, style)}>
      <View style={{height: verticalScale(sizes.base * 1.5)}} />
      {children}
    </ScrollView>
  );
};

export const FlexLayout = withStyle(View, styles.flexLayout as ViewStyle);
export const FlexRowLayout = withStyle(View, styles.flexRowLayout as ViewStyle);
export const RowLayout = withStyle(View, styles.rowLayout as ViewStyle);
