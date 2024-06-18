import React, {useState} from 'react';
import {Form, Input, InputNumber, Button, Collapse, Space, theme, Select, GlobalToken} from 'antd';
import {CaretRightOutlined, MinusCircleOutlined, PlusOutlined} from "@ant-design/icons";
import {List} from "@refinedev/antd";

const {Panel} = Collapse;
const {Option} = Select;

const getPanelStyles = (token: GlobalToken) => ({
    marginBottom: 24,
    background: token.colorFillAlter,
    borderRadius: token.borderRadiusLG,
    border: 'none',
});

const getName = (path) => path[path.length - 1];

const StringForm = ({settings, path, themeToken}) => {
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

const NumberForm = ({settings, path, themeToken}) => {
    return (
        <Form.Item
            required={settings.required || false}
            name={path}
            label={settings.disableLabel ? undefined : settings.label}
            initialValue={settings.initialValue || 0}
        >
            <InputNumber/>
        </Form.Item>
    );
}

const ArrayForm = ({settings, path, themeToken}) => {
    return (
        <Form.Item label={settings.disableLabel ? undefined : settings.label}>
            <Form.List name={path}>
                {(fields, {add, remove}) => (
                    <div>
                        {fields.map((field) => {
                            return (
                                <Space key={field.key} style={{display: 'flex', marginBottom: 8}} align="baseline">
                                    <_DynamicForm
                                        settings={{...settings.items, label: `# ${field.name}`}}
                                        path={[field.name]}
                                        themeToken={themeToken}
                                    />
                                    <MinusCircleOutlined onClick={() => remove(field.name)}/>
                                </Space>
                            );
                        })}
                        <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined/>}>
                            Add {settings.label}
                        </Button>
                    </div>
                )}
            </Form.List>
        </Form.Item>
    );
};

const ObjectForm = ({settings, path, themeToken}) => {
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
                    Object.keys(settings.properties).map(subKey => {
                        const subSettings = settings.properties[subKey]
                        return <_DynamicForm
                            settings={subSettings}
                            path={[...path, subKey]}
                            themeToken={themeToken}
                        />
                    })
                }
            </Panel>
        </Collapse>
    );
}

// Helper function to render form items based on type
const _DynamicForm = ({settings, path, themeToken}) => {
    switch (settings.type) {
        case 'string':
            return <StringForm settings={settings} path={path} themeToken={themeToken}/>;

        case 'number':
            return <NumberForm settings={settings} path={path} themeToken={themeToken}/>;

        case 'array':
            return <ArrayForm settings={settings} path={path} themeToken={themeToken}/>;

        case 'object':
            return <ObjectForm settings={settings} path={path} themeToken={themeToken}/>;

        default:
            return null;
    }
};


const DynamicForm = ({schema, path}) => {
    const {token} = theme.useToken();

    return <_DynamicForm settings={schema} path={path} themeToken={token}/>
};

export default DynamicForm;