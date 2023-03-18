---
title: "[TypeScript] อ้างอิง type จาก property ของ type อื่นด้วยท่า indexed access type"
date: 2023-03-18T10:11:59+07:00
draft: false
---

บางค่าของ TypeScript เราต้องระบุ callback function ซึ่ง type ของ callback function ก็มี parameter ที่ต้องเหมือนกัน แทนที่เราจะต้องจำและใส่ให้ตรงเอง เราสามารถใช้ indexed access type ช่วยอ้างอิง type ของ property ที่เป็น callback function ได้เลย

<!--more-->

ตัวอย่างที่เจอในงานที่ผมใช้คือ `onChange` callback ของ Table component ของ antd library มี type แบบนี้อยู่

```typescript
onChange?: (pagination: TablePaginationConfig, filters: Record<string, FilterValue | null>, sorter: SorterResult<RecordType> | SorterResult<RecordType>[], extra: TableCurrentDataSource<RecordType>) => void;
```

จะเห็นว่ายาวมากเลย ถ้าต้องเขียน calback function เองก็ต้องก็อปไปใส่เอง ไหนจะมี type parameter `RecordType` อีกที่ต้องเปลี่ยนให้ตรงกับ type ของ data ที่เราจะเอาไป render ด้วย Table component

เราจะใช้ indexed access type ของ TypeScript เข้ามาช่วยได้ เนื่องจากว่า property ของ Table component ถูก define ไว้เป็น type แบบนี้

```typescript
export interface TableProps<RecordType>
  extends Omit<
    RcTableProps<RecordType>,
    | "transformColumns"
    | "internalHooks"
    | "internalRefs"
    | "data"
    | "columns"
    | "scroll"
    | "emptyText"
  > {
  dropdownPrefixCls?: string;
  dataSource?: RcTableProps<RecordType>["data"];
  columns?: ColumnsType<RecordType>;
  pagination?: false | TablePaginationConfig;
  loading?: boolean | SpinProps;
  size?: SizeType;
  bordered?: boolean;
  locale?: TableLocale;
  onChange?: (
    pagination: TablePaginationConfig,
    filters: Record<string, FilterValue | null>,
    sorter: SorterResult<RecordType> | SorterResult<RecordType>[],
    extra: TableCurrentDataSource<RecordType>
  ) => void;
  rowSelection?: TableRowSelection<RecordType>;
  getPopupContainer?: GetPopupContainer;
  scroll?: RcTableProps<RecordType>["scroll"] & {
    scrollToFirstRowOnChange?: boolean;
  };
  sortDirections?: SortOrder[];
  showSorterTooltip?: boolean | TooltipProps;
}
```

เราสามารถ indexed access type ของ property `onChange?` ของ type `TableProps<RecordType>` ได้ ทำให้เราเขียน callback function ได้แบบนี้แทน

```typescript
const onChange: TableProps<Customer>["onChange"] = (
  pagination,
  filter,
  sorter,
  extra
) => {
  // ...
};
```

เราแค่ใช้ TableProps ระบุ generic type parameter แล้วก็ใช้ `[]` operator ในการระบุ property ไหนที่เราต้องการอ้างอิง type ของ TableProps เท่านี้เอง เช่นเคสนี้เราก็แค่ใส่ `["onChange"]` เพื่อเอา type ของ property `onChange?` จาก TableProps ให้เรา โดยที่ไม่ต้องระบุ type ของ parameter เอง

ดูเพิ่มเติมเกี่ยวกับ indexed access type ได้ที่นี่ https://www.typescriptlang.org/docs/handbook/2/indexed-access-types.html
