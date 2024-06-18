import React from 'react';
import {Badge, Card, Col, Collapse, Descriptions, Divider, Row, Space, theme, Typography} from 'antd';
import {CaretRightOutlined} from "@ant-design/icons";

const {Panel} = Collapse;
const {Text, Paragraph} = Typography;

const getPanelStyles = (token) => ({
    // marginBottom: 24,
    // background: token.colorFillAlter,
    // borderRadius: token.borderRadiusLG,
    // border: 'none',
});

const StringViewer = ({label, data}) => (
    // <Space direction="vertical">
    //     {label && <Text strong>{label}</Text>}
    //     <Text>{data || '-'}</Text>
    // </Space>
    // <Descriptions size="small" bordered>
    //     <Descriptions.Item label={label}>{data || '-'}</Descriptions.Item>
    // </Descriptions>
    // <Row style={{ marginBottom: 16 }}>
    //     <Col span={2}>
    //         {label && <Text strong>{label}</Text>}
    //     </Col>
    //     <Col span={16}>
    //         <Text>{data || '-'}</Text>
    //     </Col>
    // </Row>
    <Space direction="vertical" style={{width: '100%'}}>
        {label && <Text strong>{label}</Text>}
        <Text code ellipsis style={{maxWidth: '100%'}}>{data}</Text>
    </Space>
);

const NumberViewer = ({label, data}) => (
    <Space direction="vertical" style={{width: '100%'}}>
        {label && <Text strong>{label}</Text>}
        <Text code ellipsis style={{maxWidth: '100%'}}>{data}</Text>
    </Space>
);

const ArrayViewer = ({label, data, themeToken}) => (
    <Space direction="vertical">
        {label && <Text strong>{label}</Text>}
        {data && data.length > 0 ? data.map((item, index) => (
            <div style={{marginLeft: 16}}>
                <DynamicViewer key={index} label={`# ${index + 1}`} data={item} themeToken={themeToken}/>
            </div>
        )) : <Text>-</Text>}
    </Space>
);

const ObjectViewer = ({label, data, themeToken}) => (
    <Collapse
        // style={{width: 'fit-content', minWidth: '40%'}}
        // bordered={false}
        expandIcon={({isActive}) => <CaretRightOutlined rotate={isActive ? 90 : 0}/>}
    >
        <Panel header={label} key="1" style={{display: 'flex', flexDirection: 'column', ...getPanelStyles(themeToken)}}>
            <div style={{display: 'flex', flexDirection: 'column'}}>
                {Object.keys(data).map(subKey => (
                    <div style={{marginBottom: 8}}>
                        <DynamicViewer key={subKey} label={subKey} data={data[subKey]} themeToken={themeToken}/>
                    </div>
                ))}
            </div>
        </Panel>
    </Collapse>
);

const DynamicViewer = ({label, data, themeToken}) => {
    if (Array.isArray(data)) {
        return <ArrayViewer label={label} data={data} themeToken={themeToken}/>;
    } else if (data !== null && typeof data === 'object') {
        return <ObjectViewer label={label} data={data} themeToken={themeToken}/>;
    } else if (typeof data === 'number') {
        return <NumberViewer label={label} data={data}/>;
    } else {
        return <StringViewer label={label} data={data}/>;
    }
};

const Viewer = ({data, label}) => {
    const {token} = theme.useToken();
    return <DynamicViewer label={label} data={data} themeToken={token}/>;
};

export default Viewer;