import {useMemo, type ReactNode} from 'react';
import {StyleSheet, View} from 'react-native';
import {Divider} from './Divider';
import {withLayout} from './Layout';
import {scale, verticalScale} from './SizeMatters';
import {useTheme} from './Theme';

const styles = StyleSheet.create({
  sectionItem: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  sectionItemDivider: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
  },
  sectionItemContainer: {
    flex: 1,
  },
  sectionBox: {
    boxShadow: [
      {
        color: 'rgba(0,0,0, .2)',
        offsetX: 0,
        offsetY: 2,
        blurRadius: 1,
        spreadDistance: -1,
      },
      {
        color: 'rgba(0,0,0, .14)',
        offsetX: 0,
        offsetY: 1,
        blurRadius: 1,
        spreadDistance: 0,
      },
      {
        color: 'rgba(0,0,0, .12)',
        offsetX: 0,
        offsetY: 1,
        blurRadius: 3,
        spreadDistance: 0,
      },
    ],
  },
});

export const Section = withLayout(View, {
  padding: [0, 0, 1, 0],
});

type SectionHeaderProps = {
  icon?: ReactNode;
  children?: ReactNode;
  extra?: ReactNode;
};
export const SectionHeader = ({icon, children, extra}: SectionHeaderProps) => {
  const {sizes} = useTheme();
  return (
    <View
      style={StyleSheet.compose(
        {
          paddingBottom: verticalScale(sizes.base),
          paddingHorizontal: scale(sizes.base * 2),
        },
        styles.sectionItem,
      )}>
      {icon}
      <View style={styles.sectionItemContainer}>{children}</View>
      {extra}
    </View>
  );
};

export const SectionBox = withLayout(View, {
  padding: [0, 2],
  radius: 1,
  style: styles.sectionBox,
  color: 'surface',
});

type SectionItemProps = {
  icon?: ReactNode;
  children: NonNullable<ReactNode>;
  extra?: ReactNode;
  variant?: 'default' | 'divider';
};
export const SectionItem = ({
  variant = 'default',
  icon,
  children,
  extra,
}: SectionItemProps) => {
  const {sizes} = useTheme();

  const themeStyle = useMemo(
    () => ({
      paddingVertical: verticalScale(sizes.base * 1.5),
      gap: scale(sizes.base),
    }),
    [sizes],
  );

  if (variant === 'divider') {
    if (icon === void 0)
      return (
        <Divider
          style={StyleSheet.compose(
            {
              paddingVertical: themeStyle.paddingVertical,
            },
            styles.sectionItem,
          )}>
          <View style={styles.sectionItemContainer}>{children}</View>
          {extra}
        </Divider>
      );
    return (
      <View
        style={StyleSheet.compose(
          {
            gap: themeStyle.gap,
          },
          styles.sectionItem,
        )}>
        {icon}
        <Divider
          style={StyleSheet.compose(
            {
              paddingVertical: themeStyle.paddingVertical,
            },
            styles.sectionItemDivider,
          )}>
          <View style={styles.sectionItemContainer}>{children}</View>
          {extra}
        </Divider>
      </View>
    );
  }

  return (
    <View style={StyleSheet.compose(themeStyle, styles.sectionItem)}>
      {icon}
      <View style={styles.sectionItemContainer}>{children}</View>
      {extra}
    </View>
  );
};
