import type { CSSProperties, ReactNode } from "react";
import {
  Children,
  createContext,
  createElement,
  isValidElement,
  useContext,
  useId,
} from "react";
import AutoSizer from "react-virtualized-auto-sizer";
import { FixedSizeList } from "react-window";

export type ListItemData<T extends any> = {
  item: T;
  style?: CSSProperties;
};

const ListItemsContext = createContext<ListItemData<any>>(null!);
export const useListItems = <T extends any>() =>
  useContext<ListItemData<T>>(ListItemsContext);

export type ListItemsProps<T extends any> = {
  items: T[];
  isFetching: boolean;
  children: ReactNode;
  itemSize: number;
};
export const ListItems = <T extends any>({
  items,
  isFetching,
  children,
  itemSize,
}: ListItemsProps<T>) => {
  const id = useId();
  return (
    <AutoSizer hidden={isFetching}>
      {({ height, width }) => (
        <FixedSizeList
          height={height}
          width={width}
          itemSize={itemSize}
          itemCount={items.length}
        >
          {({ index, style }) => (
            <ListItemsContext.Provider
              value={{ item: items[index], style }}
              key={`${id}-${index}`}
            >
              {Children.map(children, (child) =>
                child && isValidElement(child)
                  ? createElement(child.type, child.props)
                  : child
              )}
            </ListItemsContext.Provider>
          )}
        </FixedSizeList>
      )}
    </AutoSizer>
  );
};
