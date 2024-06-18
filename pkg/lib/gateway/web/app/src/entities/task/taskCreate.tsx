import React, {useContext, useState} from "react";
import {
    BaseRecord,
    IResourceComponentsProps,
    useGetToPath,
    useGo,
    useList,
    useNavigation,
    useResource
} from "@refinedev/core";
import {Create, SaveButton, useForm} from "@refinedev/antd";
import {Form, notification, Spin, Typography} from "antd";
import ReactJson from "react-json-view";
import {ColorModeContext} from "../../contexts/color-mode";
import {TreeSelect} from 'antd';
import DynamicForm from "../utils/dynamicForm/DynamicForm";
import {Link} from "react-router-dom";

const {Text} = Typography;

function buildSpecificationTree(specifications) {
    const newTree = (path: string, spec: any) => ({value: {path, spec}, children: {}})
    let tree = newTree('', null)

    for (const spec of specifications) {
        const path = spec.id.split('.').filter(part => !!part)
        let currentNode = tree
        let currentPath = ''
        for (const part of path) {
            if (currentPath.length > 0) currentPath += '.'
            currentPath += part

            if (!currentNode.children[part]) currentNode.children[part] = newTree(currentPath, null)

            currentNode = currentNode.children[part]
        }

        currentNode.value.spec = spec
    }


    const convertTreeToAntdTree = (tree) => {
        return Object.keys(tree).map(key => {
            const subTree = tree[key]

            return {
                value: subTree.value.path,
                __spec: subTree.value.spec,
                title: key,
                selectable: Object.keys(subTree.children).length === 0,
                children: convertTreeToAntdTree(subTree.children)
            }
        })
    }

    return convertTreeToAntdTree(tree.children)
}

type JsonFormProps = {
    value?: number;
    onChange?: (obj: object) => void
    mode?: any
};

const JsonForm: React.FC<JsonFormProps> = ({value, onChange, mode}) => {
    const [state, setState] = useState({})

    const change = (action) => {
        setState(action.updated_src)
        if (typeof onChange === 'function') {
            onChange(action.updated_src);
        }

        return true
    }

    return <ReactJson src={state}
                      name={false}

                      onEdit={change}
                      onAdd={change}
                      onDelete={change}
                      style={{background: 'transparent'}}
                      theme={mode === 'light' ? 'summerfruit:inverted' : 'summerfruit'}
    />
}

const typeToSchema = (_type) => {
    switch (_type.name) {
        case 'string':
            return {
                type: 'string',
                required: true,
            }
        case 'number':
            return {
                type: 'number',
                required: true,
            }
        // case 'boolean':
        //     return {
        //         'type': 'string',
        //     }
        case 'struct':
            let props = {}
            for (const key in _type.sub_types) {
                const value = _type.sub_types[key]
                props[key] = {
                    ...typeToSchema(value),
                    label: key,
                }
            }

            return {
                type: 'object',
                properties: props,
            }
        case 'array':
            return {
                type: 'array',
                items: typeToSchema(_type.sub_types.item),
            }
    }
}

const generateParametersSchema = (spec) => {
    const {inputs} = spec.io

    const schema = {
        type: 'object',
        label: 'Parameters',
        properties: {}
    }

    let naming = {}, types = {}

    for (const entry of inputs.naming) {
        naming[entry.index] = entry.name
    }

    for (const entry of inputs.types) {
        types[entry.index] = entry.type
    }

    for (const index in inputs.naming) {
        const name = naming[index]
        const _type = types[index]
        schema.properties[index] = {
            ...typeToSchema(_type),
            label: name,
        }
    }

    return schema
}

export const TaskCreate: React.FC<IResourceComponentsProps> = () => {
    const getToPath = useGetToPath();
    const go = useGo();

    const {select} = useResource();
    const [api, contextHolder] = notification.useNotification();

    const {
        form,
        formProps,
        saveButtonProps,
        onFinish
    } = useForm({
        resource: 'tasks',
        action: 'create',
        redirect: false,
        onMutationSuccess(data) {

            console.log(select('tasks'))

            const url = go({
                to: {
                    resource: "tasks", // resource name or identifier
                    action: "show",
                    id: data.data.id,
                },

                // query: {
                //     filters: [
                //         {
                //             field: "id",
                //             operator: "in",
                //             value: [data.data.id],
                //         },
                //     ],
                // },
                type: "path",
            })
            console.log('created:', data, url)

            api.info({
                message: `Successfully created`,
                description: <>
                    <Text copyable>{data.data.id}</Text>
                    <br/>
                    <Link to={url}>View</Link>
                </>,
                placement: 'bottomLeft',
                duration: 30,
            });
        },
        onMutationError(err) {
            console.log('failed to create:', err)
        }
    });

    const {mode} = useContext(ColorModeContext);

    const {
        data: specifications,
        isLoading: areSpecificationsLoading,
        error: specificationsLoadingError
    } =
        useList({
            resource: 'specifications',
        })

    const [currentSpecification, setCurrentSpecification] = React.useState(null)

    if (areSpecificationsLoading) {
        return <Spin/>
    }

    if (specificationsLoadingError) {
        return <div>failed to load specs: {specificationsLoadingError}</div>
    }

    const specificationTree = buildSpecificationTree(specifications.data)
    const handleSubmit = async (values) => {
        let amountOfParameters = Object.keys(values.params || {}).length
        let parameters = []

        for (let i = 0; i < amountOfParameters; i++) {
            parameters.push(JSON.stringify(values.params[i]))
        }

        await onFinish({
            specification: values.specification,
            parameters: parameters,
            info: values.info,
        });
    };

    return (
        <>
            {contextHolder}
            <Create
                title={'New Task'}
                footerButtons={({saveButtonProps}) => (
                    <SaveButton
                        {...saveButtonProps}
                        type="primary"
                        style={{marginRight: 8}}
                        icon={null}
                        onClick={
                            () => {
                                form.submit()
                            }
                        }
                    >
                        Submit
                    </SaveButton>
                )}
            >
                <Form {...formProps} layout="vertical" onFinish={handleSubmit} onSubmitCapture={
                    () => {
                        console.log('Submitted!')
                    }

                }>
                    <Form.Item
                        label="Type"
                        name="specification"
                        rules={[
                            {
                                required: true,
                            },
                        ]}
                    >
                        <TreeSelect
                            placeholder="Select task type"
                            treeData={specificationTree}
                            showSearch
                            allowClear
                            treeExpandAction={'click'}
                            treeLine
                            onSelect={(_, data) => setCurrentSpecification(data.__spec)}
                        />
                    </Form.Item>

                    {
                        currentSpecification === null ?
                            <span>
                            Please select specification to fill parameters
                            <br/>
                            <br/>
                        </span> :
                            <DynamicForm
                                path={['params']}
                                schema={generateParametersSchema(currentSpecification)}
                            />
                    }


                    <Form.Item
                        label="Meta Info"
                        name="info"
                    >
                        <JsonForm mode={mode}/>
                    </Form.Item>
                </Form>
            </Create>
        </>

    );
};
