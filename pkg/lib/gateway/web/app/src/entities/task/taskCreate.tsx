import React, {useContext, useState} from "react";
import {IResourceComponentsProps} from "@refinedev/core";
import {AutoSaveIndicator, Create, useForm} from "@refinedev/antd";
import {Form, Input, InputNumber, Select, List, Button, Space, Collapse} from "antd";
import ReactJson from "react-json-view";
import {ColorModeContext} from "../../contexts/color-mode";
import {CloseOutlined} from "@ant-design/icons";
import {TreeSelect} from 'antd';

const treeData = [
    {
        value: 'http',
        title: 'http',
        selectable: false,
        children: [
            {
                value: 'http.Do',
                title: 'Do',
            },
            {
                value: 'http.Get',
                title: 'Get',
            },
            {
                value: 'http.Post',
                title: 'Post',
            },
        ]
    },
];

export const TaskCreate: React.FC<IResourceComponentsProps> = () => {
    const {autoSaveProps, formProps, saveButtonProps, queryResult} = useForm({
        resource: 'tasks',
        action: 'create',
        autoSave: {
            enabled: true,
            debounce: 300, // Debounce interval to trigger the auto save, defaults to 1000
            invalidateOnUnmount: true,
        },
    });

    const {mode} = useContext(ColorModeContext);

    const [dependencies, setDependencies] = useState<{}[]>([])
    const [selectedDependency, setSelectedDependency] = useState([]);

    const taskTypeToParams = {
        task: () => ({id: ''}),
        time: () => ({time: Date.now()}),
        resource: () => ({id: ''}),
    }

    return (
        <Create saveButtonProps={saveButtonProps}>
            <Form {...formProps} layout="vertical">
                <Form.Item
                    label="Type"
                    name="type"
                    rules={[
                        {
                            required: true,
                        },
                    ]}
                >
                    <TreeSelect
                        placeholder="e.g. 'Multiply'"
                        treeData={treeData}
                        showSearch
                        allowClear
                        treeExpandAction={'click'}
                        treeLine
                        onSelect={(...args) => console.log(args)}
                    />
                </Form.Item>

                <Form.Item
                    label="Inputs"
                    name="inputs"
                >
                    <Select mode='tags' allowClear></Select>
                </Form.Item>

                <Form.Item
                    label="Outputs"
                    name="outputs"
                >
                    <Select mode='tags' allowClear></Select>
                </Form.Item>

                <Form.Item
                    label="Info"
                    name="info"
                >
                    <ReactJson
                        src={{}}

                        name={false}

                        onEdit={() => true}
                        onAdd={() => true}
                        onDelete={() => true}

                        theme={mode === 'light' ? 'summerfruit:inverted' : 'summerfruit'}
                    ></ReactJson>
                </Form.Item>

                <Form.Item
                    label="Dependencies"
                    name="dependencies"
                >
                    <Select
                        mode="tags"
                        value={selectedDependency}
                        onSelect={(value, option) => {
                            let paramsGenerator = taskTypeToParams[value]
                            if (!paramsGenerator) paramsGenerator = () => ({})
                            setDependencies([{type: value, params: paramsGenerator()}, ...dependencies])
                            setSelectedDependency([])
                        }}
                        options={
                            [
                                {label: 'task', value: 'task'},
                                {label: 'time', value: 'time'},
                                {label: 'resource', value: 'resource'},
                            ]
                        }
                    />
                    <List
                        dataSource={dependencies}
                        itemLayout="horizontal"
                        // itemLayout="horizontal"
                        size="large"
                        renderItem={(item) => (
                            <List.Item
                                itemLayout="horizontal"
                                actions={[
                                    <Button
                                        danger
                                        onClick={() => {
                                            const updatedDependencies = dependencies.filter((dep) => dep != item);
                                            setDependencies(updatedDependencies);
                                        }}
                                    >Remove</Button>
                                ]}
                            >
                                <ReactJson
                                    onEdit={() => true}
                                    onAdd={() => true}
                                    onDelete={() => true}

                                    src={item}
                                    name={item.type}
                                    collapsed={true}
                                    theme={mode === 'light' ? 'summerfruit:inverted' : 'summerfruit'}/>
                            </List.Item>
                        )}
                    />

                </Form.Item>
                <AutoSaveIndicator {...autoSaveProps} />

            </Form>
        </Create>
    );
};
