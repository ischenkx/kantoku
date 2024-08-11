import React from 'react'
import {Button, Collapse, Form, GlobalToken, Input, InputNumber, Select, Space, theme} from 'antd'
import {CaretRightOutlined, MinusCircleOutlined, PlusOutlined} from '@ant-design/icons'

const {Panel} = Collapse
const {Option} = Select

const getPanelStyles = (token: GlobalToken) => ({
    marginBottom: 24,
    background: token.colorFillAlter,
    borderRadius: token.borderRadiusLG,
    border: 'none',
})

export type Settings = {
    type?: string
    required?: boolean
    disableLabel?: boolean
    label?: string
    variants?: any[]
    initialValue?: any
    // Settings for array items
    items?: Settings
    // Settings for object items
    properties?: Record<string, Settings>
}

type CommonFormParams = {
    settings: Settings,
    path: string[],
    themeToken: GlobalToken
}

const StringForm: React.FC<CommonFormParams> = ({settings, path}) => {
    return <Form.Item
        required={settings.required || false}
        label={settings.disableLabel ? undefined : settings.label}
        name={path}
        initialValue={settings.initialValue || ''}
    >
        {
            settings.variants ? (
                <Select>
                    {settings.variants.map((variant) => (
                        <Option key={variant} value={variant}>{variant}</Option>
                    ))}
                </Select>
            ) : (
                <Input/>
            )
        }
    </Form.Item>
}

const NumberForm: React.FC<CommonFormParams> = ({settings, path}) => {
    return (
        <Form.Item
            required={settings.required || false}
            name={path}
            label={settings.disableLabel ? undefined : settings.label}
            initialValue={settings.initialValue || 0}
        >
            <InputNumber/>
        </Form.Item>
    )
}

const ArrayForm: React.FC<CommonFormParams> = ({settings, path, themeToken}) => {
    return (
        <Form.Item label={settings.disableLabel ? undefined : settings.label}>
            <Form.List name={path}>
                {(fields, {add, remove}) => (
                    <div>
                        {fields.map((field) => {
                            return (
                                <Space key={field.key} style={{display: 'flex', marginBottom: 8}} align='baseline'>
                                    <_DynamicForm
                                        settings={{...settings.items, label: `# ${field.name}`}}
                                        path={[field.name?.toString() || '']}
                                        themeToken={themeToken}
                                    />
                                    <MinusCircleOutlined onClick={() => remove(field.name)}/>
                                </Space>
                            )
                        })}
                        <Button type='dashed' onClick={() => add()} block icon={<PlusOutlined/>}>
                            Add {settings.label}
                        </Button>
                    </div>
                )}
            </Form.List>
        </Form.Item>
    )
}

const ObjectForm: React.FC<CommonFormParams> = ({settings, path, themeToken}) => {
    const properties = settings.properties || {}

    return (
        <Collapse bordered={false}
            // defaultActiveKey={'1'}
                  expandIcon={({isActive}) => <CaretRightOutlined rotate={isActive ? 90 : 0}/>}>
            <Panel
                header={settings.disableLabel ? undefined : settings.label}
                key={'1'}
                style={getPanelStyles(themeToken)}
            >
                {
                    Object.entries(properties).map(([subKey, subSettings]) => {
                        return <_DynamicForm
                            settings={subSettings}
                            path={[...path, subKey]}
                            themeToken={themeToken}
                        />
                    })
                }
            </Panel>
        </Collapse>
    )
}

// Helper function to render form items based on type
const _DynamicForm: React.FC<CommonFormParams> = ({settings, path, themeToken}) => {
    switch (settings.type) {
        case 'string':
            return <StringForm settings={settings} path={path} themeToken={themeToken}/>

        case 'number':
            return <NumberForm settings={settings} path={path} themeToken={themeToken}/>

        case 'array':
            return <ArrayForm settings={settings} path={path} themeToken={themeToken}/>

        case 'object':
            return <ObjectForm settings={settings} path={path} themeToken={themeToken}/>

        default:
            return null
    }
}


const DynamicForm: React.FC<CommonFormParams> = ({settings, path, themeToken}) => {
    return <_DynamicForm settings={settings} path={path} themeToken={themeToken}/>
}

export default DynamicForm