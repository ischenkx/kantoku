import React, {useContext, useState} from "react";
import {BaseRecord, IResourceComponentsProps, useList, useMany, useNavigation, useShow} from "@refinedev/core";
import {List, Show, ShowButton, TagField, TextField} from "@refinedev/antd";
import {Collapse, Descriptions, Divider, Space, Spin, Table, Tag, Typography} from "antd";

const {Panel} = Collapse
const {Title, Text} = Typography;
import ReactJson from 'react-json-view'

import {List as AntdList} from 'antd';
import {Link} from "react-router-dom";
import {ColorModeContext} from "../../contexts/color-mode";
import {Status as ResourceStatus} from "../resource/resourceList";
import Viewer from "../utils/objectViewer/DynamicViewer";
import {TaskStatus} from "./taskList";
import {CaretRightOutlined} from "@ant-design/icons";

function formatUnixTime(unixTime) {
    // Create a new Date object using the Unix time (multiply by 1000 to convert to milliseconds)
    const date = new Date(unixTime * 1000);

    // Extract the date components
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0'); // Months are 0-based, so add 1
    const day = String(date.getDate()).padStart(2, '0');

    // Extract the time components
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');

    // Format the date and time string
    const formattedDateTime = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;

    return formattedDateTime;
}


const ResourceIdList = ({header, data}) => {
    const columns = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            render: id => <Link to={showUrl('resources', id)}>{id}</Link>,
        },
        {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            render: status => <ResourceStatus value={status}/>
        },
    ];

    if (!data) {
        data = []
    }

    if (data.length === 0) {
        return (
            <Collapse>
                <Panel header={header} key="1">
                    <Table dataSource={[]} columns={columns}/>
                </Panel>
            </Collapse>
        );
    }

    const [activeKey, setActiveKey] = useState(null);
    const {showUrl} = useNavigation();

    const {data: resources, isLoading, error} = useMany({
        resource: 'resources',
        ids: data,
    })

    if (isLoading) {
        return <Spin/>
    }

    if (error) {
        return <p>Error!</p>
    }

    const handlePanelChange = (key) => {
        setActiveKey(activeKey === key ? null : key);
    };

    return (
        <Collapse onChange={handlePanelChange} activeKey={activeKey}>
            <Panel header={header} key="1">
                <Table dataSource={resources?.data} columns={columns}/>
            </Panel>
        </Collapse>
    );
};


export const TaskShow: React.FC<IResourceComponentsProps> = () => {
    const {showUrl} = useNavigation();
    const {mode} = useContext(ColorModeContext);
    const [inputsList, setInputResourceList] = useState([])
    const [outputsList, setOutputResourceList] = useState([])
    const [allResourcesList, setAllResourcesList] = useState([])
    const [resourcesSynced, setResourcesSynced] = useState(false)

    const {queryResult} = useShow({
        successNotification(result, _1, _2) {
            setInputResourceList(result?.data.inputs)
            setOutputResourceList(result?.data.outputs)
            setAllResourcesList([
                ...result?.data.inputs,
                ...result?.data.outputs
            ])
        }
    });
    const {data, isLoading: isTaskLoading, error: recordLoadingError } = queryResult;
    const record = data?.data;

    const {
        data: specifications,
        isLoading: areSpecificationsLoading,
        error: specificationsLoadingError
    } =
        useList({
            resource: 'specifications',
        })

    const {
        data: resourcesList,
        isLoading: areResourcesLoading,
        error: resourcesLoadingError
    } = useMany({
        resource: 'resources',
        ids: allResourcesList,
        successNotification() {
            if (!isTaskLoading) {
                setResourcesSynced(true)
            }
        }
    })

    if (areSpecificationsLoading || areResourcesLoading || isTaskLoading || !resourcesSynced) {
        return <Spin/>
    }

    if (recordLoadingError) {
        return <div>failed to load a record: {recordLoadingError}</div>
    }

    if (specificationsLoadingError) {
        return <div>failed to load specs: {specificationsLoadingError}</div>
    }

    if (resourcesLoadingError) {
        return <div>failed to load resources: {resourcesLoadingError}</div>
    }

    const specification = specifications?.data?.find((spec) => spec.id == record?.info?.type)

    const allResources = {}
    for (const resource of resourcesList?.data) {
        allResources[resource.id] = resource
    }

    let parametersObject = null, resultsObject = null;

    if (specification) {
        console.log(allResources, specification, inputsList, outputsList, resourcesList,
            resourcesLoadingError)
        parametersObject = {}
        resultsObject = {}

        const inputsNames = {}, outputsNames = {}
        for (const namingObject of specification.io.inputs.naming) {
            inputsNames[namingObject.index] = namingObject.name
        }

        for (const namingObject of specification.io.outputs.naming) {
            outputsNames[namingObject.index] = namingObject.name
        }

        for (const index in inputsList) {
            const inputResourceID = inputsList[index];
            const resource = allResources[inputResourceID];
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
            const outputResourceID = outputsList[index];
            const resource = allResources[outputResourceID];
            const name = outputsNames[index]

            if (resource.status !== 'ready') {
                resultsObject[name] = 'N/A'
                continue
            }

            const data = JSON.parse(resource.value)
            resultsObject[name] = data
        }
    }

    const resourcesColumns = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            render: (_, record) => {
                let {id, label} = record
                // console.log('data:', id, record)
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
            render: status => <ResourceStatus value={status}/>
        },
    ];

    const resourceTableDataSource = [
        ...inputsList.map((id, index) => {
            let resource = allResources[id];
            return {
                id: resource.id,
                status: resource.status,
                label: `Input #${index + 1}`
            }
        }),
        ...outputsList.map((id, index) => {
            let resource = allResources[id];
            return {
                id: resource.id,
                status: resource.status,
                label: `Output #${index + 1}`
            }
        })
    ]

    const dependenciesColumns = [
        {
            title: 'Data',
            dataIndex: 'data',
            key: 'data',
            render: data => <Text code ellipsis>{data}</Text>,
        },
        {
            title: 'Name',
            dataIndex: 'name',
            key: 'name',
            render: name => <span>{name}</span>
        },
    ];

    const dependenciesTableDataSource = record?.info?.dependencies?.specs || [];

    console.log(record)

    return (
        <Show isLoading={isTaskLoading}>
            <Descriptions layout={'horizontal'} column={1}>
                <Descriptions.Item label="Id"><TextField copyable value={record?.id}/></Descriptions.Item>
                <Descriptions.Item label="Status">
                    <TaskStatus status={record?.info?.status} subStatus={record?.info?.sub_status}/>
                </Descriptions.Item>
                {specification && <Descriptions.Item label="Spec">{specification.id}</Descriptions.Item>}
                {record?.info?.updated_at && <Descriptions.Item
                    label="Updated At">{formatUnixTime(record?.info?.updated_at)}</Descriptions.Item>}

            </Descriptions>

            <Title level={5}>Info</Title>
            <ReactJson
                src={record?.info}
                name={false}
                theme={mode === 'light' ? 'summerfruit:inverted' : 'summerfruit'}
                style={{background: 'transparent'}}
                collapsed={true}
                collapseStringsAfterLength={80}
            />

            <br/>

            {/*<Divider/>*/}

            {parametersObject && <><Viewer data={parametersObject} label={'Parameters'}/> <br/></>}
            {resultsObject && <><Viewer data={resultsObject} label={'Results'}/><br/></>}



            <Collapse
                // bordered={false}
                expandIcon={({isActive}) => <CaretRightOutlined rotate={isActive ? 90 : 0}/>}
            >
                <Panel header={'Resources'} key="1">
                    <Table dataSource={resourceTableDataSource} columns={resourcesColumns}/>
                </Panel>
            </Collapse>

            <br/>

            <Collapse
                expandIcon={({isActive}) => <CaretRightOutlined rotate={isActive ? 90 : 0}/>}
            >
                <Panel header={'Dependencies'} key="1">
                    <Table dataSource={dependenciesTableDataSource} columns={dependenciesColumns}/>
                </Panel>
            </Collapse>
        </Show>
    );
};
