import type {ComponentType, PropsWithChildren, ReactNode} from 'react';
import {StyleSheet, View} from 'react-native';
import {scale, verticalScale} from 'react-native-size-matters';
import {useTheme} from './Theme';

const styles = StyleSheet.create({
  listItem: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'flex-start',
  },
  listItemText: {
    flexGrow: 1,
  },
});

type ListProps<T extends any, Props extends PropsWithChildren<{}>> = Props & {
  data: T[];
};

export function withList<
  T extends any,
  BaseProps extends PropsWithChildren<{}>,
>(
  BaseComponent: ComponentType<BaseProps>,
  renderItem: (item: T, index: number, items: T[]) => ReactNode,
): ComponentType<ListProps<T, BaseProps>> {
  return ({data, children, ...props}: ListProps<T, BaseProps>) => {
    if (data.length <= 0)
      return (
        <BaseComponent {...(props as unknown as BaseProps)}>
          {children}
        </BaseComponent>
      );
    return (
      <BaseComponent {...(props as unknown as BaseProps)}>
        {data.map(renderItem)}
      </BaseComponent>
    );
  };
}

type ListItemProps = PropsWithChildren<{icon?: ReactNode; text: ReactNode}>;
export const ListItem = ({children, icon, text}: ListItemProps) => {
  const {sizes} = useTheme();
  return (
    <View
      style={StyleSheet.compose(styles.listItem, {
        gap: scale(sizes.base),
        paddingVertical: verticalScale(sizes.base),
        paddingHorizontal: scale(sizes.base),
      })}>
      {icon && <View>{icon}</View>}
      <View style={styles.listItemText}>{text}</View>
      {children}
    </View>
  );
};
