import {
  type ComponentType,
  type PropsWithChildren,
  type ReactNode,
} from 'react';
import type {ViewProps} from 'react-native';
import {View} from 'react-native';

type ListBaseProps<
  T extends unknown,
  Props extends PropsWithChildren<{}>,
> = Props & {
  data: T[];
};

export function withList<
  T extends unknown,
  BaseProps extends PropsWithChildren<{}>,
>(
  BaseComponent: ComponentType<BaseProps>,
  renderItem: (item: T, index: number, items: T[]) => ReactNode,
): ComponentType<ListBaseProps<T, BaseProps>> {
  return ({data, children, ...props}: ListBaseProps<T, BaseProps>) => (
    <BaseComponent {...(props as unknown as BaseProps)}>
      {data.length > 0 ? data.map(renderItem) : children}
    </BaseComponent>
  );
}

type ListProps<T extends unknown> = ListBaseProps<T, ViewProps> & {
  renderItem: (item: T, index: number, items: T[]) => ReactNode;
};
const List = <T extends unknown>({
  data,
  renderItem,
  children,
  ...props
}: ListProps<T>) => {
  return (
    <View {...props}>{data.length > 0 ? data.map(renderItem) : children}</View>
  );
};

export default List;
