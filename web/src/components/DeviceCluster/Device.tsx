import {
  default as DevicesIcon,
  SearchIcon,
  ListIcon,
  TableIcon,
  LockIcon,
} from "./Icon";
import { withExtraState, ExtraType } from "./ExtraBase";

import { type ComponentProps } from "react";
import { cn } from "@/lib/utils";
import { useAppDispatch } from "@/components/App/Context";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";
import {
  Item,
  ItemContent,
  ItemDescription,
  ItemTitle,
} from "@/components/ui/item";
import { Marker, MarkerContent } from "@/components/ui/marker";
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
  FieldContent,
  FieldTitle,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";

export const DeviceAction = () => (
  <>
    <DeviceNewAction />
  </>
);

export const DeviceCard = () => (
  <Card>
    <CardHeader className="flex items-center justify-between gap-3">
      <DeviceTableSearchFilter />
      <div className="flex items-center gap-1">
        <DeviceDisplayControl />
        <Separator
          orientation="vertical"
          className="mx-2 data-[orientation=vertical]:h-4"
        />
        <DeviceRemoveAction />
      </div>
    </CardHeader>
    <CardContent className="flex flex-col gap-2">
      <DeviceTable />
      <div className="flex items-center justify-end py-4">
        <DeviceTableSelection />
        <DeviceTablePagination />
      </div>
    </CardContent>
  </Card>
);

export enum DeviceLevel {
  Owner = 1,
  Member = 2,
  Guest = 3,
}
export const DeviceLevelTempTexts = {
  [DeviceLevel.Owner]: "管理员",
  [DeviceLevel.Member]: "成员",
  [DeviceLevel.Guest]: "访客",
};
export const devices = [
  {
    id: 1,
    name: "desktop",
    memo: "-",
    version: "1.0.0",
    type: "pc",
    level: 1,
    height: 1080,
    apps: [],
    tags: [],
    duration: 24 * 60 * 60,
  },
  {
    id: 2,
    name: "OnePlus ACE 5",
    memo: "-",
    version: "1.0.0",
    type: "mobile",
    level: 1,
    height: 200,
    apps: [],
    tags: [],
    duration: 60 * 60,
  },
  {
    id: 3,
    name: "xps13",
    memo: "-",
    version: "1.0.0",
    type: "pc",
    level: 1,
    height: 300,
    apps: [],
    tags: [],
    duration: 0,
  },
  {
    id: 4,
    name: "raspberrypi",
    memo: "-",
    version: "1.0.0",
    type: "pc",
    level: 2,
    height: 1200,
    apps: [],
    tags: [],
    duration: 124 * 60 * 60,
  },
  {
    id: 5,
    name: "OPPO Reno 5",
    memo: "-",
    version: "1.0.0",
    type: "mobile",
    level: 2,
    height: 300,
    apps: [],
    tags: [],
    duration: 0,
  },
];
const DeviceTable = () => (
  <Table>
    <TableHeader>
      <TableRow>
        <TableHead className="w-8">
          <Checkbox />
        </TableHead>
        <TableHead>设备名</TableHead>
        <TableHead>设备等级</TableHead>
        <TableHead>应用版本</TableHead>
        <TableHead>同步序号</TableHead>
        <TableHead>在线时长</TableHead>
        <TableHead className="text-center">功能</TableHead>
      </TableRow>
    </TableHeader>
    <TableBody>
      {devices.map((device) => (
        <TableRow
          key={device.id}
          className={device.duration > 0 ? "" : "text-muted-foreground"}
        >
          <TableCell className="text-center">
            {device.id !== 2 ? (
              <Checkbox />
            ) : (
              <LockIcon size="16" className="text-muted-foreground" />
            )}
          </TableCell>
          <TableCell className="font-medium">{device.name}</TableCell>
          <TableCell>
            {DeviceLevelTempTexts[device.level as DeviceLevel]}
          </TableCell>
          <TableCell>{device.version}</TableCell>
          <TableCell>{device.height}</TableCell>
          <TableCell>{device.duration}</TableCell>
          <TableCell className="text-center">
            <DeviceEditAction disabled={device.id === 2} />
          </TableCell>
        </TableRow>
      ))}
    </TableBody>
  </Table>
);

const DeviceTableSelection = () => (
  <div className="flex-1 text-sm text-muted-foreground">
    0 of 5 row(s) selected.
  </div>
);

function DeviceTablePagination() {
  return (
    <Pagination className="w-fit">
      <PaginationContent>
        <PaginationItem>
          <PaginationPrevious href="#" />
        </PaginationItem>
        <PaginationItem>
          <PaginationLink href="#">1</PaginationLink>
        </PaginationItem>
        <PaginationItem>
          <PaginationLink href="#" isActive>
            2
          </PaginationLink>
        </PaginationItem>
        <PaginationItem>
          <PaginationLink href="#">3</PaginationLink>
        </PaginationItem>
        <PaginationItem>
          <PaginationEllipsis />
        </PaginationItem>
        <PaginationItem>
          <PaginationNext href="#" />
        </PaginationItem>
      </PaginationContent>
    </Pagination>
  );
}

const DeviceTableSearchFilter = () => (
  <InputGroup className="max-w-sm">
    <InputGroupInput placeholder="筛选设备" />
    <InputGroupAddon>
      <SearchIcon />
    </InputGroupAddon>
  </InputGroup>
);

const DeviceDisplayControl = () => (
  <ToggleGroup
    variant="outline"
    className="gap-0 outline-none"
    defaultValue={["table"]}
  >
    <ToggleGroupItem value="table" className="rounded-r-none border-r-0">
      <TableIcon />
      表格
    </ToggleGroupItem>
    <ToggleGroupItem value="list" className="rounded-l-none">
      <ListIcon />
      列表
    </ToggleGroupItem>
  </ToggleGroup>
);

const DeviceRemoveAction = () => {
  const dispatch = useAppDispatch();
  const handleClick = () => {
    dispatch?.(
      withExtraState({
        type: ExtraType.DevicesRemove,
        clusterId: 1,
        deviceIds: [1, 2, 3],
      }),
    );
  };
  return (
    <Button variant="destructive" onClick={handleClick}>
      移除设备
    </Button>
  );
};

const DeviceNewAction = () => (
  <Button variant="ghost" size="sm" disabled>
    添加新设备
  </Button>
);

const DeviceEditAction = (
  props: Omit<ComponentProps<typeof Button>, "children">,
) => {
  const dispatch = useAppDispatch();
  const handleClick = () => {
    dispatch?.(
      withExtraState({
        type: ExtraType.DeviceEdit,
        clusterId: 1,
        deviceId: 1,
      }),
    );
  };
  return (
    <Button {...props} variant="ghost" size="sm" onClick={handleClick}>
      设置
    </Button>
  );
};

export const DeviceList = () => {
  return (
    <>
      {devices.map((device) => (
        <Item variant="muted" key={device.id}>
          <ItemContent>
            <ItemTitle>{device.name}</ItemTitle>
            <ItemDescription>{device.memo}</ItemDescription>
          </ItemContent>
          <ItemContent className="flex-none text-center">
            <ItemDescription
              className={cn(
                "flex items-center gap-1",
                device.level === DeviceLevel.Owner && "text-destructive",
              )}
            >
              {DeviceLevelTempTexts[device.level as DeviceLevel]}
            </ItemDescription>
          </ItemContent>
        </Item>
      ))}
    </>
  );
};

export const DeviceFormInput = () => {
  return (
    <FieldGroup>
      <DeviceNameInput />
      <DeviceMemoInput />
      <Marker variant="separator">
        <MarkerContent>设备权限</MarkerContent>
      </Marker>
      <DeviceLevelInput />
    </FieldGroup>
  );
};

const DeviceNameInput = () => {
  return (
    <Field>
      <FieldLabel htmlFor="device-name">备注名</FieldLabel>
      <Input id="device-name" type="text" placeholder="e.g. desktop" />
      <FieldDescription>好名称可以让你更容易找到设备</FieldDescription>
    </Field>
  );
};

const DeviceMemoInput = () => {
  return (
    <Field>
      <FieldLabel htmlFor="device-memo">备忘</FieldLabel>
      <Textarea id="device-memo" placeholder="e.g. 家里台式机" rows={5} />
    </Field>
  );
};

const DeviceLevelInput = () => {
  return (
    <RadioGroup defaultValue="plus">
      <FieldLabel htmlFor="plus-plan">
        <Field orientation="horizontal">
          <FieldContent>
            <FieldTitle>所有</FieldTitle>
            <FieldDescription>
              <p>允许对设备组进行修改,设置或移除。</p>
              <p>允许对设备组内的其他设备进行修改,设置或移除。</p>
              <p>包含用户权限。</p>
            </FieldDescription>
          </FieldContent>
          <RadioGroupItem value="plus" id="plus-plan" />
        </Field>
      </FieldLabel>
      <FieldLabel htmlFor="pro-plan">
        <Field orientation="horizontal">
          <FieldContent>
            <FieldTitle>用户</FieldTitle>
            <FieldDescription>
              <p>允许使用设备组内其他设备上的资源和功能</p>
              <p>允许提供资源和功能给设备组内其他设备使用</p>
            </FieldDescription>
          </FieldContent>
          <RadioGroupItem value="pro" id="pro-plan" />
        </Field>
      </FieldLabel>
      <FieldLabel htmlFor="enterprise-plan">
        <Field orientation="horizontal">
          <FieldContent>
            <FieldTitle>访客</FieldTitle>
            <FieldDescription>
              <p>通过安全验证后,才允许使用设备组内其他设备上的部分资源和功能</p>
            </FieldDescription>
          </FieldContent>
          <RadioGroupItem value="enterprise" id="enterprise-plan" />
        </Field>
      </FieldLabel>
    </RadioGroup>
  );
};
