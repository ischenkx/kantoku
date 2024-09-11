import React, {useContext, useState} from 'react'
import {
    IResourceComponentsProps,
    OpenNotificationParams,
    useGo,
    useList,
    useMany,
    useNavigation,
    useShow
} from '@refinedev/core'
import {Show, TextField} from '@refinedev/antd'
import {Button, Collapse, Descriptions, Space, Spin, Steps, Table, Tag, Typography} from 'antd'
import ReactJson from 'react-json-view'
import {Link} from 'react-router-dom'
import {ColorModeContext} from '../../contexts/color-mode'
import {Status as ResourceStatus} from '../resource/resourceList'
import Viewer from '../utils/objectViewer/DynamicViewer'
import {TaskStatus} from './taskList'
import {CaretRightOutlined} from '@ant-design/icons'
import {Resource, Task} from '../../api/generated'
import {APIWrapper} from '../../providers/common'

const {Panel} = Collapse
const {Title, Text} = Typography

function formatUnixTime(unixTime: number): string {
    const date = new Date(unixTime * 1000)

    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0') // Months are 0-based, so add 1
    const day = String(date.getDate()).padStart(2, '0')

    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    const seconds = String(date.getSeconds()).padStart(2, '0')

    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

const ResourceView: React.FC<{ inputs: Resource[], outputs: Resource[] }> =
    ({
         inputs,
         outputs
     }) => {
        inputs ||= []
        outputs ||= []

        const {showUrl} = useNavigation()

        type rowT = {
            id: string
            status: string
            label: string
        }

        const dataSource: rowT[] = [
            ...inputs.map((resource, index) => {
                return {
                    id: resource.id,
                    status: resource.status,
                    label: `Input #${index + 1}`
                }
            }),
            ...outputs.map((resource, index) => {
                return {
                    id: resource.id,
                    status: resource.status,
                    label: `Output #${index + 1}`
                }
            })
        ]

        const columns = [
            {
                title: 'ID',
                dataIndex: 'id',
                key: 'id',
                render: (_: string, {id, label}: rowT) => {
                    return <>
                        <Link to={showUrl('resources', id)}>{id}</Link>
                        <Tag style={{marginLeft: 16}}>{label}</Tag>
                    </>
                },
            },
            {
                title: 'Status',
                dataIndex: 'status',
                key: 'status',
                render: (status: string) => <ResourceStatus value={status}/>
            },
        ]

        return <Table pagination={{hideOnSinglePage: true}} dataSource={dataSource} columns={columns}/>
    }

const TaskRestartsView: React.FC<{ task: Task }> =
    ({task}) => {
        const [shouldFetch, setShouldFetch] = useState<boolean>(true)
        const [taskRestarts, setTaskRestarts] = useState<Task[] | null>(null)
        const go = useGo()

        if (shouldFetch) {
            APIWrapper.getTaskRestarts(task)
                .catch(err => {
                    console.log('failed to get task restarts:', err)
                })
                .then(tasks => {
                    if (!tasks) return
                    setShouldFetch(false)
                    setTaskRestarts(tasks)
                })

            return null
        }

        if (!taskRestarts) return null

        console.log('task restarts:', taskRestarts)

        taskRestarts.sort((a: Task, b: Task): number => {
            const aInfo = a.info as Record<string, any>,
                bInfo = b.info as Record<string, any>
            const at = aInfo['updated_at'] as number || 0,
                bt = bInfo['updated_at'] as number || 0
            return at - bt
        })

        const steps = taskRestarts.map(t => {
            const taskUrl = go({
                to: {
                    resource: 'tasks',
                    action: 'show',
                    id: t.id,
                },
                type: 'path',
            })

            const info = t.info as Record<string, any>

            return {
                title: task.id === t.id ? <>{t.id}</> : <Link to={taskUrl || ''}>{t.id}</Link>,
                description: <TaskStatus status={info['status'] || ''} subStatus={info['sub_status'] || ''}/>
            }
        })

        return (
            <Steps
                progressDot
                current={steps.length}
                direction='vertical'
                items={steps}
            />
        )
    }


export const TaskShow: React.FC<IResourceComponentsProps> = () => {
    const {showUrl} = useNavigation()
    const {mode} = useContext(ColorModeContext)
    const [inputsList, setInputResourceList] = useState<string[]>([])
    const [outputsList, setOutputResourceList] = useState<string[]>([])
    const [allResourcesList, setAllResourcesList] = useState<string[]>([])

    const {queryResult} =
        useShow<Task>(
            {
                successNotification(result): false | OpenNotificationParams {
                    const inputs: string[] = result?.data?.inputs || []
                    const outputs: string[] = result?.data?.outputs || []

                    setInputResourceList(inputs)
                    setOutputResourceList(outputs)
                    setAllResourcesList([...inputs, ...outputs])

                    return false
                }
            })
    const {data: taskResponse, isLoading: isTaskLoading, error: recordLoadingError} = queryResult
    const task = taskResponse?.data

    const {
        data: specifications,
        isLoading: areSpecificationsLoading,
        error: specificationsLoadingError
    } =
        useList({
            resource: 'specifications',
        })

    const {
        data: resourcesResponse,
        isLoading: areResourcesLoading,
        error: resourcesLoadingError
    } = useMany<Resource>({
        resource: 'resources',
        ids: allResourcesList,
        queryOptions: {
            enabled: allResourcesList.length > 0
        }
    })

    const go = useGo()

    if (areSpecificationsLoading || areResourcesLoading || isTaskLoading || !task) {
        return <Spin/>
    }

    if (recordLoadingError) {
        return <>failed to load a record: {recordLoadingError}</>
    }

    if (specificationsLoadingError) {
        return <>failed to load specs: {specificationsLoadingError}</>
    }

    if (resourcesLoadingError) {
        return <>failed to load resources: {resourcesLoadingError}</>
    }

    const taskInfo = (task?.info || {}) as Record<string, any>

    const resourcesList: Resource[] = resourcesResponse?.data || []

    const specification = specifications?.data?.find((spec) => spec.id == taskInfo.type)

    const allResources: Record<string, Resource> = {}
    for (const resource of resourcesList) {
        allResources[resource.id] = resource
    }

    const inputResources: Resource[] = task?.inputs.map(id => allResources[id]) || []
    const outputResources: Resource[] = task?.outputs.map(id => allResources[id]) || []

    let parametersObject: Record<string, any> | null = null,
        resultsObject: Record<string, any> | null = null

    if (specification) {
        console.log('specification:', specification)

        parametersObject = {}
        resultsObject = {}

        const inputsNames: Record<number, string> = {},
            outputsNames: Record<number, string> = {}
        for (const namingObject of (specification.io.inputs.naming || [])) {
            inputsNames[namingObject.index] = namingObject.name
        }

        for (const namingObject of (specification.io.outputs.naming || [])) {
            outputsNames[namingObject.index] = namingObject.name
        }

        for (const index in inputsList) {
            const inputResourceID = inputsList[index]
            const resource = allResources[inputResourceID]
            const name = inputsNames[index]

            if (resource.status !== 'ready') {
                parametersObject[name] = 'N/A'
                continue
            }

            let data = null
            try {
                data = JSON.parse(resource.value)
            } catch (e) {
                console.log('failed to parse json:', e)
                data = resource.value
            }

            parametersObject[name] = data
        }

        for (const index in outputsList) {
            const outputResourceID = outputsList[index]
            const resource = allResources[outputResourceID]
            const name = outputsNames[index]

            if (resource.status !== 'ready') {
                resultsObject[name] = 'N/A'
                continue
            }

            const data = JSON.parse(resource.value)
            resultsObject[name] = data
        }
    }

    const dependenciesColumns = [
        {
            title: 'Data',
            dataIndex: 'data',
            key: 'data',
            render: (data: string) => <Text code ellipsis>{data}</Text>,
        },
        {
            title: 'Name',
            dataIndex: 'name',
            key: 'name',
            render: (name: string) => <span>{name}</span>
        },
    ]

    const dependenciesTableDataSource = taskInfo.dependencies?.specs || []

    return (
        <Show isLoading={isTaskLoading}>
            <Descriptions layout={'horizontal'} column={1}>
                <Descriptions.Item label='Id'><TextField copyable value={task?.id}/></Descriptions.Item>
                <Descriptions.Item label='Status'>
                    <TaskStatus status={taskInfo.status} subStatus={taskInfo.sub_status}/>
                </Descriptions.Item>
                {specification && <Descriptions.Item label='Spec'>{specification.id}</Descriptions.Item>}
                {taskInfo.updated_at && <Descriptions.Item
                    label='Updated At'>{formatUnixTime(taskInfo.updated_at || 0)}</Descriptions.Item>}

                <Space>
                    {taskInfo['context_id'] && <Button onClick={() => {
                        return go({
                            to: {
                                resource: "flow",
                                action: "list",
                            },
                            query: {
                                'context_id': taskInfo['context_id']
                            },
                            type: "replace",
                        })
                    }}>Go to context</Button>}
                </Space>
            </Descriptions>

            <Title level={5}>Info</Title>
            <ReactJson
                src={task?.info || {}}
                name={false}
                theme={mode === 'light' ? 'summerfruit:inverted' : 'summerfruit'}
                style={{background: 'transparent'}}
                collapsed={true}
                collapseStringsAfterLength={80}
            />

            <br/>

            {/*<Divider/>*/}

            {parametersObject &&
                (
                    inputResources.length > 0 ?
                        <><Viewer data={parametersObject} label={'Parameters'}/> <br/></>
                        :
                        null
                )
            }
            {resultsObject && (
                outputResources.length > 0 ?
                    <><Viewer data={resultsObject} label={'Results'}/><br/></>
                    :
                    null
            )
            }

            <Collapse
                expandIcon={({isActive}) => <CaretRightOutlined rotate={isActive ? 90 : 0}/>}
            >
                <Panel header={'Resources'} key='2'>
                    <ResourceView inputs={inputResources} outputs={outputResources}/>
                </Panel>


                {dependenciesTableDataSource.length > 0 &&
                    <Panel header={'Dependencies'} key='3'>
                        <Table dataSource={dependenciesTableDataSource} columns={dependenciesColumns}/>
                    </Panel>
                }

                {taskInfo['restart_root'] &&
                    <Panel header={'Restarts'} key='1'>
                        <TaskRestartsView task={task}/>
                    </Panel>
                }

            </Collapse>
        </Show>
    )
}

