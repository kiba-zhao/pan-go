import type {ComponentType, Context} from 'react';
import {useContext} from 'react';

export function withTransform<
  Props extends {},
  BaseProps extends {},
  Options extends {},
>(
  BaseComponent: ComponentType<BaseProps>,
  transform: (props: Props, options: Options) => BaseProps,
  options: Options,
): ComponentType<Props> {
  return (props: Props) => {
    const transformProps = transform(props, options);
    return <BaseComponent {...transformProps} />;
  };
}

export function withContext<
  Props extends {},
  BaseProps extends {},
  State extends unknown,
  Enhance extends State,
>(
  BaseComponent: ComponentType<BaseProps>,
  context: Context<State>,
  createProps: (state: Enhance, props: Props) => BaseProps,
): ComponentType<Props> {
  return (props: Props) => {
    const state = useContext(context) as Enhance;
    return <BaseComponent {...createProps(state, props)} />;
  };
}
