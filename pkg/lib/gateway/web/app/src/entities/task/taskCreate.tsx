import React, {useContext, useState} from 'react'
import {BaseRecord, HttpError, IResourceComponentsProps, useGo, useList} from '@refinedev/core'
import {Create, SaveButton, useForm} from '@refinedev/antd'
import {Form, notification, Spin, theme, TreeSelect, Typography} from 'antd'
import ReactJson, {InteractionProps} from 'react-json-view'
import {ColorModeContext} from '../../contexts/color-mode'
import DynamicForm, {Settings as TypeSettings} from '../utils/dynamicForm/DynamicForm'
import {Link} from 'react-router-dom'
import {buildSpecificationTree, convertTreeToAntdTree} from '../utils/specs/tree'
import {Specification, Type} from '../utils/specs/specification'

const {Text} = Typography

type JsonFormProps = {
    value?: number;
    onChange?: (obj: object) => void
    mode?: any
};

type CreationFormValues = {
    parameters: Record<number, any> | string[]
    specification: string
    info: any
}

const JsonForm: React.FC<JsonFormProps> = ({onChange, mode}) => {
    const [state, setState] = useState({})

    const change = (action: InteractionProps) => {
        setState(action.updated_src)
        if (typeof onChange === 'function') {
            onChange(action.updated_src)
        }

        return true
    }

    return <ReactJson
        src={state}
        name={false}
        onEdit={change}
        onAdd={change}
        onDelete={change}
        style={{background: 'transparent'}}
        theme={mode === 'light' ? 'summerfruit:inverted' : 'summerfruit'}
    />
}

const typeToSettings = (typ: Type): TypeSettings => {
    switch (typ.name) {
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
            // eslint-disable-next-line no-case-declarations
            const props: Record<string, any> = {}
            for (const key in typ.sub_types) {
                const value = typ.sub_types[key]
                props[key] = {
                    ...typeToSettings(value),
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
                items: typeToSettings(typ.sub_types.item),
            }
    }

    return {}
}

const generateParametersSettings = (spec: Specification) => {
    const {inputs} = spec.io

    const settings: TypeSettings = {
        type: 'object',
        label: 'Parameters',
    }

    const naming: Record<number, string> = {}, types: Record<number, Type> = {}

    for (const entry of (inputs.naming || [])) {
        naming[entry.index] = entry.name
    }

    for (const entry of (inputs.types || [])) {
        types[entry.index] = entry.type
    }

    settings.properties = {}
    for (const index in (inputs.naming || [])) {
        const name = naming[index]
        const _type = types[index]
        settings.properties[index] = {
            ...typeToSettings(_type),
            label: name,
        }
    }

    return settings
}

export const TaskCreate: React.FC<IResourceComponentsProps> = () => {
    const go = useGo()
    const [notifications, notificationsContextHolder] = notification.useNotification()
    const {mode} = useContext(ColorModeContext)
    const {token: themeToken} = theme.useToken()

    const {
        form,
        formProps,
        onFinish
    } = useForm<BaseRecord, HttpError, CreationFormValues>({
        resource: 'tasks',
        action: 'create',
        redirect: false,
        onMutationSuccess(data) {
            console.log(`created a new task: ${data.data.id}`)

            const url = go({
                to: {
                    resource: 'tasks', // resource name or identifier
                    action: 'show',
                    id: data.data.id || '',
                },
                type: 'path',
            })

            notifications.info({
                message: `Successfully created`,
                description: <>
                    <Text copyable>{data.data.id}</Text>
                    <br/>
                    <Link to={url || ''}>View</Link>
                </>,
                placement: 'bottomLeft',
                duration: 30,
            })
        },
        onMutationError(err) {
            console.log(`failed to create a new task: ${err}`)
        }
    })

    const {
        data: specifications,
        isLoading: areSpecificationsLoading,
        error: specificationsLoadingError
    } =
        useList<Specification>({
            resource: 'specifications',
        })

    const [currentSpecification, setCurrentSpecification]
        = React.useState<Specification | null>(null)

    if (areSpecificationsLoading) {
        return <Spin/>
    }

    if (specificationsLoadingError) {
        return <>failed to load specs: {specificationsLoadingError}</>
    }

    const specificationTree = buildSpecificationTree(specifications.data)
    const antdTree = convertTreeToAntdTree(specificationTree.children)

    const handleSubmit = async (values: CreationFormValues) => {
        const rawParams = (values.parameters || {}) as Record<number, any>

        const amountOfParameters = Object.keys(rawParams).length
        const parameters: string[] = []

        for (let i = 0; i < amountOfParameters; i++) {
            parameters.push(JSON.stringify(rawParams[i]))
        }

        await onFinish({
            specification: values.specification,
            parameters: parameters,
            info: values.info,
        })
    }

    return (
        <>
            {notificationsContextHolder}
            <Create
                title={'New Task'}
                footerButtons={({saveButtonProps}) => (
                    <SaveButton
                        {...saveButtonProps}
                        type='primary'
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
                <Form
                    {...formProps}
                    layout='vertical'
                    onFinish={handleSubmit}
                    onSubmitCapture={() => {
                        console.log('submitted')
                    }}>
                    <Form.Item
                        label='Type'
                        name='specification'
                        rules={[
                            {
                                required: true,
                            },
                        ]}
                    >
                        <TreeSelect
                            placeholder='Select task type'
                            treeData={antdTree}
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
                            (
                                (currentSpecification.io.inputs.naming || []).length > 0 ?
                                    (<DynamicForm
                                        path={['parameters']}
                                        settings={generateParametersSettings(currentSpecification)}
                                        themeToken={themeToken}
                                    />) :
                                    // <div>No parameters</div>
                                    null
                            )
                    }

                    <Form.Item
                        label='Meta Info'
                        name='info'
                    >
                        <JsonForm mode={mode}/>
                    </Form.Item>
                </Form>
            </Create>
        </>
    )
}
