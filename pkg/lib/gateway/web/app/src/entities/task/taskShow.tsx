import React, {useContext, useState} from "react";
import {BaseRecord, IResourceComponentsProps, useMany, useNavigation, useShow} from "@refinedev/core";
import {List, Show, ShowButton, TagField, TextField} from "@refinedev/antd";
import {Collapse, Space, Spin, Table, Typography} from "antd";

const {Panel} = Collapse
const {Title} = Typography;
import ReactJson from 'react-json-view'

import {List as AntdList} from 'antd';
import {Link} from "react-router-dom";
import {ColorModeContext} from "../../contexts/color-mode";
import {Status as ResourceStatus} from "../resource/resourceList";

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
                    <Table dataSource={[]} columns={columns} />
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
        return <Spin />
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
                <Table dataSource={resources?.data} columns={columns} />
            </Panel>
        </Collapse>
    );
};

export const TaskShow: React.FC<IResourceComponentsProps> = () => {
    const {queryResult} = useShow();
    const {data, isLoading} = queryResult;
    const record = data?.data;

    const { mode } = useContext(ColorModeContext);

    return (
        <Show isLoading={isLoading}>
            <Title level={5}>ID</Title>
            <TextField copyable value={record?.id}/>
            <br/>
            <br/>

            <Title level={5}>Info</Title>
            <ReactJson
                src={record?.info}
                name={null}
                theme={mode === 'light' ? 'summerfruit:inverted' : 'summerfruit' }
                collapseStringsAfterLength={80}
            />
            <br/>

            <ResourceIdList header={<Title level={5}>Inputs</Title>} data={record?.inputs}/>
            <br/>

            <ResourceIdList header={<Title level={5}>Outputs</Title>} data={record?.outputs}/>
        </Show>
    );
};
