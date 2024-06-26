import React, {useMemo} from "react";
import {IResourceComponentsProps, BaseRecord, CrudFilters, useParsed, getDefaultFilter} from "@refinedev/core";
import {useTable, List, ShowButton, TagField, TextField, useSelect} from "@refinedev/antd";
import {Table, Space, Tag, Row, Col, Card, FormProps, Form, Input, Select, DatePicker, Button, Tooltip} from "antd";
import {SearchOutlined} from "@ant-design/icons";
import dayjs from "dayjs";

const {RangePicker} = DatePicker;

export const TaskStatus: React.FC<{ status: string, subStatus: string }> = ({status, subStatus}) => {
    status = status || 'unknown'
    let color, value;
    switch (status) {
        case "finished":
            switch (subStatus) {
                case "ok":
                    color = "green";
                    value = "ok";
                    break;
                case "failed":
                    color = "red";
                    value = "failed";
                    break;
                default:
                    color = undefined;
                    value = "finished";
                    break
            }
            break

        case "ready":
            color = "blue";
            value = "ready";
            break;
        case "received":
            color = "blue";
            value = "received";
            break;
        case "unknown":
            color = "yellow";
            value = "unknown";
            break;
    }
    return <TagField value={value} color={color}/>
}

const Filter: React.FC<{ formProps: FormProps, filters: CrudFilters }> = ({formProps, filters}) => {
    const updatedAt = useMemo(() => {
        const start = getDefaultFilter("info.updated_at", filters, "gte");
        const end = getDefaultFilter("info.updated_at", filters, "lte");

        const startFrom = start ? dayjs(start) : undefined;
        const endAt = end ? dayjs(end) : undefined;

        if (start || end) {
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
                    tagRender={({value}) => <TaskStatus status={value} subStatus={''}/>}
                    // dropdownRender={option => <Status value={option.props.label}/>}
                    options={[
                        // {
                        //     label: "ok",
                        //     value: "ok",
                        // },
                        // {
                        //     label: "failed",
                        //     value: "failed",
                        // },
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
                <RangePicker format="DD/MM/YYYY hh:mm A" showTime={{use12Hours: false}} allowEmpty={[true, true]}/>
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
                            title="Specification"
                            render={(_, record: BaseRecord) =>
                                <span>
                                    {record.info?.type ? record.info?.type : '-' }
                                </span>
                            }
                        />


                        <Table.Column
                            title="Status"
                            render={(_, record: BaseRecord) => <TaskStatus
                                status={record.info?.status}
                                subStatus={record.info?.sub_status}/>
                            }
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
                                    >
                                        Show
                                    </ShowButton>
                                </Space>
                            )}
                        />
                    </Table>
                </List>
            </Col>
        </Row>
    );
};