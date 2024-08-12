import React, {ReactElement, useMemo} from 'react'
import {
    BaseRecord,
    CrudFilters,
    getDefaultFilter,
    HttpError,
    IResourceComponentsProps,
    useGo,
    useInvalidate
} from '@refinedev/core'
import {List, RefreshButton, ShowButton, TagField, useTable} from '@refinedev/antd'
import {Button, Card, Col, DatePicker, Form, FormProps, notification, Row, Select, Space, Table, Typography} from 'antd'
import dayjs from 'dayjs'
import './taskList.css'
import {API} from '../../providers/common'
import {Link} from 'react-router-dom'
import {RedoOutlined} from '@ant-design/icons'
import {RefineButtonClassNames, RefineButtonTestIds} from '@refinedev/ui-types'

const {Text} = Typography
const {RangePicker} = DatePicker

export const TaskStatus: React.FC<{ status: string, subStatus: string }> = ({status, subStatus}) => {
    status = status || 'unknown'
    let color, value
    switch (status) {
        case 'finished':
            switch (subStatus) {
                case 'ok':
                    color = 'green'
                    value = 'ok'
                    break
                case 'failed':
                    color = 'red'
                    value = 'failed'
                    break
                default:
                    color = undefined
                    value = 'finished'
                    break
            }
            break

        case 'ready':
            color = 'blue'
            value = 'ready'
            break
        case 'received':
            color = 'blue'
            value = 'received'
            break
        case 'unknown':
            color = 'yellow'
            value = 'unknown'
            break
    }
    return <TagField value={value} color={color}/>
}

const Filter: React.FC<{ formProps: FormProps, filters: CrudFilters }> = ({formProps, filters}) => {
    const updatedAt = useMemo(() => {
        const start = getDefaultFilter('info.updated_at', filters, 'gte')
        const end = getDefaultFilter('info.updated_at', filters, 'lte')

        const startFrom = start ? dayjs(start) : undefined
        const endAt = end ? dayjs(end) : undefined

        if (start || end) {
            return [startFrom, endAt]
        }
        return undefined
    }, [filters])


    return (
        <Form
            layout='vertical'
            {...formProps}
            initialValues={
                {
                    ids: getDefaultFilter('id', filters, 'in'),
                    statuses: getDefaultFilter('info.status', filters, 'in'),
                    updatedAt,
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
            <Form.Item label='Status' name='statuses'>
                <Select
                    mode='multiple'
                    allowClear
                    tagRender={({value}) => <TaskStatus status={value} subStatus={''}/>}
                    options={[
                        {
                            label: 'ready',
                            value: 'ready',
                        },
                        {
                            label: 'received',
                            value: 'received',
                        },
                        {
                            label: 'finished',
                            value: 'finished',
                        },
                        {
                            label: 'unknown',
                            value: 'unknown',
                        },
                    ]}
                    placeholder='Task Status'
                />
            </Form.Item>
            <Form.Item label='Updated At' name='updatedAt'>
                <RangePicker format='DD/MM/YYYY hh:mm A' showTime={{use12Hours: false}} allowEmpty={[true, true]}/>
            </Form.Item>
            <Form.Item>
                <Button htmlType='submit' type='primary'>
                    Filter
                </Button>
            </Form.Item>
        </Form>
    )
}

export const TaskList: React.FC<IResourceComponentsProps> = () => {
    const go = useGo()
    const [notifications, notificationsContextHolder] = notification.useNotification()
    const invalidate = useInvalidate()

    const {tableProps, filters, searchFormProps, tableQueryResult} = useTable<BaseRecord, HttpError, {
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
            const filters: CrudFilters = []
            const {ids, statuses} = params
            let {updatedAt} = params

            if (!updatedAt) {
                updatedAt = [undefined, undefined]
            }

            filters.push(
                {
                    field: 'id',
                    operator: 'in',
                    value: ids,
                },
                {
                    field: 'info.status',
                    operator: 'in',
                    value: statuses,
                },
                {
                    field: 'info.updated_at',
                    operator: 'gte',
                    value: updatedAt[0] ? updatedAt[0].toISOString() : undefined,
                },
                {
                    field: 'info.updated_at',
                    operator: 'lte',
                    value: updatedAt[1] ? updatedAt[1].toISOString() : undefined,
                },
                {
                    field: 'requestId',
                    operator: 'eq',
                    value: Math.round(Math.random() * 1e9).toString(16)
                }
            )

            return filters
        }
    })

    console.log(tableProps)
    tableProps.scroll = undefined

    const handleRestart = async (id: string, ...args: any[]) => {
        console.log('restarting:', id)

        try {
            const response = await API.tasksRestartPost({id: id})

            if (!response) {
                console.log('no response')
                return
            }
            const newTaskURL = go({
                to: {
                    resource: 'tasks',
                    action: 'show',
                    id: response.data.id,
                },
                type: 'path',
            })

            notifications.info({
                message: `Successfully restarted ${id}`,
                description: <>
                    <Text copyable>{response.data.id}</Text>
                    <br/>
                    <Link to={newTaskURL || ''}>View</Link>
                </>,
                placement: 'bottomLeft',
                duration: 30,
            })

            invalidate({resource: 'tasks', invalidates: ['list']})
                .then(() => console.log('invalidated tasks'))
                .catch(err => console.log('failed to invalidate tasks:', err))

        } catch (err1) {
            console.log('failed to restart:', err1)

            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-ignore
            const {message} = err1
            notifications.error({
                message: `Failed to restart ${id}`,
                description: <>
                    Error: {message}
                </>,
                placement: 'bottomLeft',
                duration: 10,
            })

            throw err1
        }
    }

    const columns = [
        {
            title: 'ID',
            width: '15%',
            ellipsis: true,
            render(_: any, record: BaseRecord): string | ReactElement {
                return <span>
                    <Text copyable={{text: record.id as string}}></Text>
                    <Text style={{
                        width: '100%',
                        boxSizing: 'border-box',
                        marginLeft: '4px'
                    }}>
                        {record.id}
                    </Text>
                </span>
            }
        },
        {
            title: 'Status',
            width: '15%',
            ellipsis: true,
            render(_: any, record: BaseRecord): string | ReactElement {
                return <TaskStatus status={record.info?.status}
                                   subStatus={record.info?.sub_status}/>
            }
        },

        {
            title: 'Specification',
            width: '25%',
            ellipsis: true,
            render(_: any, record: BaseRecord): string | ReactElement {
                return <Text
                    // copyable={true}
                    // ellipsis={true}
                    style={{
                        width: '100%',
                        boxSizing: 'border-box',
                    }}>
                    {record.info?.type ? record.info?.type : '-'}
                </Text>
            }
        },
        {
            title: 'Actions',
            width: '25%',
            render(_: any, record: BaseRecord): string | ReactElement {
                return <Space>
                    <ShowButton
                        hideText
                        size='small'
                        recordItemId={record.id}
                    >
                        Show
                    </ShowButton>

                    <Button
                        icon={<RedoOutlined/>}
                        size='small'
                        disabled={!(record.info?.status === 'finished'
                            && record.info?.sub_status === 'failed'
                            && !record.info.restarted)}
                        onClick={(event) => {
                            const tar = event.currentTarget as unknown as { disabled: boolean }
                            tar.disabled = true
                            handleRestart(record.id as string).catch((err) => {
                                tar.disabled = false
                            })
                        }}/>
                </Space>
            }
        }
    ]


    return (
        <>
            {notificationsContextHolder}
            <Row gutter={[16, 16]}>
                <Col lg={5} xs={24}>
                    <Card>
                        <Filter filters={filters} formProps={searchFormProps}/>
                    </Card>
                </Col>
                <Col lg={19} xs={24}>
                    <List
                        headerButtons={({defaultButtons}) => (
                            <>
                                {defaultButtons}
                                <RefreshButton onClick={() => tableQueryResult.refetch().then(result => {
                                    console.log('refetched:', result)
                                })}/>
                            </>
                        )}
                    >
                        <Table
                            {...tableProps}
                            rowKey='id'
                            bordered
                            size='middle'
                            style={{width: '100%'}}
                            pagination={{
                                ...tableProps.pagination,
                                position: ['bottomCenter'],
                                size: 'small',
                                showTotal(total, range) {
                                    return <span>{`${range[0]}-${range[1]} of ${total} items`}</span>
                                }
                            }}
                            columns={columns}
                            sticky

                        />
                    </List>
                </Col>
            </Row>
        </>

    )
}