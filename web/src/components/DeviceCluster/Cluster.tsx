import { ArrowLeftRight, SearchIcon, default as DevicesIcon } from "./Icon";
import { withExtraState, ExtraType } from "./ExtraBase";
import {
  type PropsWithChildren,
  useMemo,
  useRef,
  useEffect,
  type ComponentProps,
} from "react";

import { cn } from "@/lib/utils";
import { useTranslation } from "@/components/App/I18Next";
import { useAppDispatch } from "@/components/App/Context";
import {
  List,
  ListItem,
  ListItemText,
  ListItemButton,
  ListItemVariant,
} from "@/components/App/List";

import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Item,
  ItemContent,
  ItemDescription,
  ItemTitle,
} from "@/components/ui/item";
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  ChartLegend,
  type ChartConfig,
} from "@/components/ui/chart";
import { Label, Pie, PieChart } from "recharts";

export const ClusterAction = () => (
  <>
    <ClusterNewAction />
    <Separator
      orientation="vertical"
      className="mr-2 data-[orientation=vertical]:h-4"
    />
    <ClusterSwitchAction />
  </>
);

export const ClusterInfoCard = () => (
  <Card>
    <CardHeader>
      <CardTitle>Zhao</CardTitle>
      <CardDescription>我自己设备</CardDescription>
      <CardAction>
        <ClusterRemoveAction />
      </CardAction>
    </CardHeader>
    <CardContent>
      <ClusterChart />
    </CardContent>
    <CardContent className="flex flex-col gap-2">
      <ClusterInfoField title="健康数量">6</ClusterInfoField>
      <ClusterInfoField title="同步序号">32</ClusterInfoField>
      <ClusterInfoField title="更新时间">3h ago</ClusterInfoField>
    </CardContent>
    <CardFooter className="flex-wrap items-center gap-3">
      <ClusterEditAction />
      <ClusterPassphraseEditAction />
    </CardFooter>
  </Card>
);

const chartData = [
  { state: "online", total: 6, fill: "var(--color-online)" },
  { state: "offline", total: 2, fill: "var(--color-offline)" },
];
const chartConfig = {
  online: {
    label: "在线",
    color: "var(--chart-1)",
  },
  offline: {
    label: "离线",
    color: "var(--chart-5)",
  },
} satisfies ChartConfig;
const ClusterChart = () => {
  const totalDevices = useMemo(() => {
    return chartData.reduce((acc, curr) => acc + curr.total, 0);
  }, []);
  return (
    <ChartContainer
      config={chartConfig}
      className="mx-auto aspect-square max-w-3xs"
    >
      <PieChart>
        <Pie
          data={chartData}
          dataKey="total"
          nameKey="state"
          innerRadius={60}
          strokeWidth={5}
        >
          <Label
            content={({ viewBox }) => {
              if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                return (
                  <text
                    x={viewBox.cx}
                    y={viewBox.cy}
                    textAnchor="middle"
                    dominantBaseline="middle"
                  >
                    <tspan
                      x={viewBox.cx}
                      y={viewBox.cy}
                      className="fill-foreground text-3xl font-bold"
                    >
                      {totalDevices.toLocaleString()}
                    </tspan>
                    <tspan
                      x={viewBox.cx}
                      y={(viewBox.cy || 0) + 24}
                      className="fill-muted-foreground"
                    >
                      设备总数
                    </tspan>
                  </text>
                );
              }
            }}
          />
        </Pie>
        <ChartTooltip cursor={false} content={<ChartTooltipContent />} />
        <ChartLegend
          iconSize={8}
          verticalAlign="top"
          formatter={renderClusterChartLegendContent}
          className="-translate-y-2 flex-wrap gap-2 *:basis-1/4 *:justify-center"
        />
      </PieChart>
    </ChartContainer>
  );
};

function renderClusterChartLegendContent(
  name: string,
  entry: any,
  idx: number,
) {
  const payload = entry.payload as (typeof chartData)[0];
  return (
    <span className="h-4">
      {chartConfig[name as keyof typeof chartConfig].label} {payload.total}
    </span>
  );
}

type ClusterInfoFieldProps = PropsWithChildren<{
  title: string;
}>;
const ClusterInfoField = ({ title, children }: ClusterInfoFieldProps) => (
  <div className="flex w-full items-center justify-between text-sm">
    <h6 className="text-muted-foreground">{title}</h6>
    {children}
  </div>
);

const ClusterRemoveAction = () => {
  const dispatch = useAppDispatch();
  const handleClick = () => {
    dispatch?.(withExtraState({ type: ExtraType.ClusterRemove, clusterId: 1 }));
  };
  return (
    <Button variant="destructive" onClick={handleClick}>
      移除
    </Button>
  );
};

const ClusterNewAction = () => {
  const dispatch = useAppDispatch();
  const handleClick = () => {
    dispatch?.(withExtraState({ type: ExtraType.ClusterAdd }));
  };
  return (
    <Button variant="ghost" size="sm" onClick={handleClick}>
      新增
    </Button>
  );
};

const ClusterSwitchAction = () => {
  const dispatch = useAppDispatch();
  const handleClick = () => {
    dispatch?.(withExtraState({ type: ExtraType.ClusterSwitch }));
  };
  return (
    <Button variant="outline" size="sm" onClick={handleClick}>
      <ArrowLeftRight data-icon="inline-start" />
      切换
    </Button>
  );
};

const ClusterEditAction = () => {
  const dispatch = useAppDispatch();
  const handleClick = () => {
    dispatch?.(withExtraState({ type: ExtraType.ClusterEdit }));
  };
  return (
    <Button className="w-full" onClick={handleClick}>
      设置
    </Button>
  );
};

const ClusterPassphraseEditAction = () => {
  const dispatch = useAppDispatch();
  const handleClick = () => {
    dispatch?.(withExtraState({ type: ExtraType.ClusterPassphraseEdit }));
  };
  return (
    <Button variant="outline" className="w-full" onClick={handleClick}>
      管理口令
    </Button>
  );
};

export const ClusterSearchFilter = ({ onEsc }: { onEsc?: () => void }) => {
  const { t } = useTranslation();
  const ref = useRef<HTMLInputElement>(null);
  useEffect(() => {
    ref.current?.focus();
  }, []);

  const handleEsc = () => {
    onEsc?.();
  };

  return (
    <InputGroup className="border-0! bg-transparent! ring-0!">
      <InputGroupInput placeholder="搜索设备组" ref={ref} />
      <InputGroupAddon>
        <SearchIcon />
      </InputGroupAddon>
      <InputGroupAddon align="inline-end">
        <Button variant="outline" size="xs" onClick={handleEsc}>
          Esc
        </Button>
      </InputGroupAddon>
    </InputGroup>
  );
};

const clusters = [
  {
    id: 1,
    name: "Zhao",
    memo: "Zhao的设备组",
  },
  {
    id: 2,
    name: "朵小皮",
  },
  {
    id: 3,
    name: "路人甲",
    memo: "路边遇见的路人兄弟",
  },
  {
    id: 4,
    name: "其他",
  },
];

type ClusterListProps = {
  onSelect?: (clusterId: number) => void;
  selected?: number;
};
export const ClusterList = ({ selected, onSelect }: ClusterListProps) => {
  return (
    <List className="text-sm">
      {clusters.map((cluster) => (
        <ListItem
          hover={ListItemVariant.Primary}
          active={
            selected && selected === cluster.id
              ? ListItemVariant.Muted
              : undefined
          }
          key={cluster.id}
          className={cn(
            "h-13 odd:border-border odd:border-y-1 last:border-b-1",
            !(selected && selected === cluster.id) &&
              "odd:hover:border-transparent!",
          )}
        >
          <ListItemButton
            className="py-0 h-full"
            onClick={() => onSelect?.(cluster.id)}
            // to={
            //   namespace === ClusterName && Number(clusterId) === cluster.id
            //     ? "#"

            // }
            // state={locationState}
          >
            <ListItemText>
              {cluster.name}
              <p className="text-xs text-muted-foreground pt-1">
                {cluster.memo}
              </p>
            </ListItemText>
          </ListItemButton>
        </ListItem>
      ))}
    </List>
  );
};

export const ClusterInfo = ({ variant }: { variant?: "disused" }) => {
  return (
    <Item variant={variant ? "muted" : "default"}>
      <ItemContent>
        <ItemTitle>Zhao</ItemTitle>
        <ItemDescription>Zhao的设备组</ItemDescription>
      </ItemContent>
      <ItemContent className="flex-none text-center">
        <div className="flex">
          <span
            className={cn(
              "flex size-2 rounded-full bg-destructive self-center mr-1",
              !variant && "hidden",
            )}
          />
          <DevicesIcon className="size-6" />
          <span className="text-xs">32</span>
        </div>
      </ItemContent>
    </Item>
  );
};

export const ClusterInfoInputFields = () => {
  return (
    <>
      <Field>
        <FieldLabel htmlFor="cluster-name">备注名</FieldLabel>
        <Input id="cluster-name" type="text" placeholder="e.g. Zhao" />
        <FieldDescription>好名称可以让你更容易快速识别</FieldDescription>
      </Field>
      <Field>
        <FieldLabel htmlFor="cluster-memo">备忘</FieldLabel>
        <Textarea id="cluster-memo" placeholder="e.g. Zhao的设备组" rows={5} />
      </Field>
    </>
  );
};

export const ClusterPassphraseForm = ({
  ...props
}: Omit<ComponentProps<"form">, "children">) => {
  return (
    <form {...props}>
      <FieldGroup>
        <Field>
          <FieldLabel htmlFor="password">旧口令</FieldLabel>
          <Input id="password" type="password" placeholder="••••••••" />
          <FieldDescription>请输入旧口令</FieldDescription>
        </Field>
        <Field>
          <FieldLabel htmlFor="new-password">新口令</FieldLabel>
          <Input id="new-password" type="password" placeholder="••••••••" />
          <FieldDescription>
            新口令至少8位,包含字母、数字和特殊字符中的至少两种.
          </FieldDescription>
        </Field>
        <Field>
          <FieldLabel htmlFor="confirm-password">确认新口令</FieldLabel>
          <Input id="confirm-password" type="password" placeholder="••••••••" />
          <FieldDescription>
            请重复输入新口令，确保两次输入一致
          </FieldDescription>
        </Field>
      </FieldGroup>
    </form>
  );
};

export enum PassphraseVariant {
  Default = "default",
  Confirm = "confirm",
  Old = "old",
  New = "new",
  NewConfirm = "new-confirm",
}
const PassphraseVariants = {
  [PassphraseVariant.Default]: {
    label: "安全口令",
    description: "请输入安全口令",
  },
  [PassphraseVariant.Confirm]: {
    label: "确认口令",
    description: "请重复输入安全口令，确保两次输入一致",
  },
  [PassphraseVariant.Old]: {
    label: "旧口令",
    description: "请输入旧口令",
  },
  [PassphraseVariant.New]: {
    label: "新口令",
    description: "新口令至少8位,包含字母、数字和特殊字符中的至少两种.",
  },
  [PassphraseVariant.NewConfirm]: {
    label: "确认新口令",
    description: "请重复输入新口令，确保两次输入一致",
  },
};
type ClusterPassphraseInputFieldProps = {
  variant?: PassphraseVariant;
} & Omit<ComponentProps<typeof Input>, "type">;
export const ClusterPassphraseInputField = ({
  variant = PassphraseVariant.Default,
}: ClusterPassphraseInputFieldProps) => {
  return (
    <Field>
      <FieldLabel htmlFor="password">
        {PassphraseVariants[variant].label}
      </FieldLabel>
      <Input id="password" type="password" placeholder="••••••••" />
      <FieldDescription>
        {PassphraseVariants[variant].description}
      </FieldDescription>
    </Field>
  );
};
