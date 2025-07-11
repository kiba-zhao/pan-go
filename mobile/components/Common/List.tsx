import {
  type ComponentType,
  type PropsWithChildren,
  type ReactNode,
} from 'react';

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
