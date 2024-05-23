import React, {useMemo} from "react";
import {IResourceComponentsProps, BaseRecord, CrudFilters, useParsed, getDefaultFilter} from "@refinedev/core";
import {useTable, List, ShowButton, TagField, TextField, useSelect} from "@refinedev/antd";
import {Table, Space, Tag, Row, Col, Card, FormProps, Form, Input, Select, DatePicker, Button, Tooltip} from "antd";
import {SearchOutlined} from "@ant-design/icons";
import dayjs from "dayjs";

const {RangePicker} = DatePicker;

const Status: React.FC<{ value: string }> = ({value}) => {
    value = value || 'unknown'
    let color;
    switch (value) {
        case "ok":
            color = "green";
            break;
        case "failed":
            color = "red";
            break;
        case "ready":
            color = "blue";
            break;
        case "received":
            color = "blue";
            break;
        case "unknown":
            color = "yellow";
            break;
    }
    return <TagField value={value} color={color}/>
}

const Filter: React.FC<{ formProps: FormProps, filters: CrudFilters }> = ({formProps, filters}) => {
    const updatedAt = useMemo(() => {
        const start = getDefaultFilter("info.updatedAt", filters, "gte");
        const end = getDefaultFilter("info.updatedAt", filters, "lte");

        const startFrom = dayjs(start);
        const endAt = dayjs(end);

        if (start && end) {
            return [startFrom, endAt];
        }
        return undefined;
    }, [filters]);


    return (
        <Form
            layout="vertical"
            {...formProps}
            initialValues={
                {
                    ids: getDefaultFilter('id', filters, 'in'),
                    statuses: getDefaultFilter('info.status', filters, 'in'),
                    updatedAt,
                }
            }
        >
            <Form.Item label="Search" name="ids">
                <Select
                    allowClear
                    mode="tags"
                    placeholder={"IDs"}
                    options={[]}
                />
            </Form.Item>
            <Form.Item label="Status" name="statuses">
                <Select
                    mode="multiple"
                    allowClear
                    tagRender={({value}) => <Status value={value}/>}
                    // dropdownRender={option => <Status value={option.props.label}/>}
                    options={[
                        {
                            label: "ok",
                            value: "ok",
                        },
                        {
                            label: "failed",
                            value: "failed",
                        },
                        {
                            label: "ready",
                            value: "ready",
                        },
                        {
                            label: "received",
                            value: "received",
                        },
                        {
                            label: "finished",
                            value: "finished",
                        },
                        {
                            label: "unknown",
                            value: "unknown",
                        },
                    ]}
                    placeholder="Task Status"
                />
            </Form.Item>
            <Form.Item label="Updated At" name="updatedAt">
                <RangePicker format="DD/MM/YYYY hh:mm A" showTime={{use12Hours: false}}/>
            </Form.Item>
            <Form.Item>
                <Button htmlType="submit" type="primary">
                    Filter
                </Button>
            </Form.Item>
        </Form>
    );
};

export const TaskList: React.FC<IResourceComponentsProps> = () => {
    const {tableProps, filters, searchFormProps, setCurrent, current} = useTable<any, any, {
        ids: string[];
        statuses: string[];
        updatedAt: any
    }>({
        syncWithLocation: true,
        pagination: {
            mode: 'server'
        },
        filters: {
            defaultBehavior: 'replace',
        },
        onSearch(params) {
            const filters: CrudFilters = [];
            let {ids, statuses, updatedAt} = params;

            if (!updatedAt) {
                updatedAt = [undefined, undefined]
            }

            filters.push(
                {
                    field: "id",
                    operator: "in",
                    value: ids,
                },
                {
                    field: "info.status",
                    operator: "in",
                    value: statuses,
                },
                {
                    field: "info.updated_at",
                    operator: "gte",
                    value: updatedAt[0] ? updatedAt[0].toISOString() : undefined,
                },
                {
                    field: "info.updated_at",
                    operator: "lte",
                    value: updatedAt[1] ? updatedAt[1].toISOString() : undefined,
                },
                {
                    field: 'requestId',
                    operator: 'eq',
                    value: Math.round(Math.random() * 1e9).toString(16)
                }
            )

            console.log('filters:', filters)

            return filters;
        }
    });

    return (
        <Row gutter={[16, 16]}>
            <Col lg={6} xs={24}>
                <Card>
                    <Filter filters={filters} formProps={searchFormProps}/>
                </Card>
            </Col>
            <Col lg={18} xs={24}>
                <List>
                    <Table
                        {...tableProps}
                        rowKey="id"
                        pagination={{
                            ...tableProps.pagination,
                            position: ["bottomCenter"],
                            size: "small",
                            showTotal(total, range) {
                                return <span>{`${range[0]}-${range[1]} of ${total} items`}</span>
                            }
                        }}
                    >
                        <Table.Column
                            title="ID"
                            render={
                                (_, record: BaseRecord) => (
                                    // <Tooltip placement="bottomRight" title={record.id}>
                                    <TextField value={record.id} copyable/>
                                    // </Tooltip>
                                )
                            }
                        />

                        <Table.Column
                            title="Status"
                            render={(_, record: BaseRecord) => <Status value={record.info?.status}/>}
                        />

                        <Table.Column
                            title="Actions"
                            dataIndex="actions"
                            render={(_, record: BaseRecord) => (
                                <Space>
                                    <ShowButton
                                        hideText
                                        size="small"
                                        recordItemId={record.id}
                                    />
                                </Space>
                            )}
                        />
                    </Table>
                </List>
            </Col>
        </Row>
    );
};
