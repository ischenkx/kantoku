import React from 'react'
import {IResourceComponentsProps, BaseRecord, CrudFilters, getDefaultFilter} from '@refinedev/core'
import {useTable, List, ShowButton, TagField, TextField, CreateButton} from '@refinedev/antd'
import {Table, Space, FormProps, Form, Select, Button, Col, Card, Row} from 'antd'


export const Status: React.FC<{ value: string }> = ({value}) => {
    const colorMapping: Record<string, string> = {
        'ready': 'green',
        'allocated': 'blue',
        'does_not_exist': 'yellow',
    }
    value = value || 'unknown'
    let color: string = colorMapping[value]


    switch (value) {
        case 'ready':
            color = 'green'
            break
        case 'allocated':
            color = 'blue'
            break
        case 'does_not_exist':
            color = 'yellow'
            break
    }
    return <TagField value={value} color={color}/>
}

const Filter: React.FC<{ formProps: FormProps, filters: CrudFilters }> = ({formProps, filters}) => {
    return (
        <Form
            {...formProps}
            layout='vertical'
            initialValues={
                {
                    ids: (() => {
                        let filter = getDefaultFilter('id', filters, 'in')
                        if (!Array.isArray(filter) && (typeof filter === 'object')) filter = Object.values(filter)

                        return filter
                    })(),
                }
            }
        >
            <Form.Item label='Search' name='ids'>
                <Select
                    allowClear
                    mode='tags'
                    placeholder={'IDs'}
                    options={[]}
                />
            </Form.Item>
            <Form.Item>
                <Button htmlType='submit' type='primary'>
                    Filter
                </Button>
            </Form.Item>
        </Form>
    )
}

const AllocateButton: React.FC = () => (
    <CreateButton title='Allocate' resource='resources'>Allocate</CreateButton>
)

export const ResourceList: React.FC<IResourceComponentsProps> = () => {
    const {tableProps, searchFormProps, filters} = useTable<any, any, { ids: string[] }>({
        syncWithLocation: true,
        pagination: {
            mode: 'client'
        },
        onSearch: (params) => {
            const filters: CrudFilters = []
            const {ids} = params

            filters.push(
                {
                    field: 'id',
                    operator: 'in',
                    value: ids,
                },
            )

            return filters
        }
    })

    return (
        <Row gutter={[16, 16]}>
            <Col lg={6} xs={24}>
                <Card>
                    <Filter filters={filters} formProps={searchFormProps}/>
                </Card>
            </Col>
            <Col lg={18} xs={24}>
                <List headerButtons={<AllocateButton/>}>
                    <Table
                        {...tableProps}
                        rowKey={'id'}
                        pagination={{
                            ...tableProps.pagination,
                            position: ['bottomCenter'],
                            size: 'small',
                        }}
                    >
                        <Table.Column
                            title='ID'
                            sorter={{multiple: 1}}
                            dataIndex='id'
                            render={(_, record: BaseRecord) => <TextField value={record.id} copyable/>}
                        />
                        <Table.Column
                            title='Status'
                            sorter={{multiple: 2}}
                            dataIndex='status'
                            render={(_, record: BaseRecord) => <Status value={record.status}/>}
                        />
                        <Table.Column
                            title='Actions'
                            dataIndex='actions'
                            render={(_, record: BaseRecord) => (
                                <Space>
                                    <ShowButton
                                        hideText
                                        size='small'
                                        recordItemId={record.id}
                                    />
                                </Space>
                            )}
                        />
                    </Table>
                </List>
            </Col>
        </Row>
    )
}
