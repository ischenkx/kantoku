import React, {useState, useMemo} from "react";
import {IResourceComponentsProps, useList} from "@refinedev/core";
import {Tree, Card, Layout, Spin, Descriptions, Typography, TreeSelect, Input} from 'antd';
import './style.css'

const {TreeNode} = Tree;
const {Content, Sider} = Layout;
const {Text} = Typography;
const {Search} = Input;

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
                key: subTree.value.path,
                __spec: subTree.value.spec,
                title: key,
                selectable: Object.keys(subTree.children).length === 0,
                children: convertTreeToAntdTree(subTree.children)
            }
        })
    }

    return convertTreeToAntdTree(tree.children)
}

const getParentKey = (key, tree) => {
    let parentKey;
    for (let i = 0; i < tree.length; i++) {
        const node = tree[i];
        if (node.children) {
            if (node.children.some(item => item.value === key)) {
                parentKey = node.value;
            } else if (getParentKey(key, node.children)) {
                parentKey = getParentKey(key, node.children);
            }
        }
    }
    return parentKey;
};

const generateList = (data, dataList = []) => {
    for (let i = 0; i < data.length; i++) {
        const node = data[i];
        const {value} = node;
        dataList.push({key: value, title: value});
        if (node.children) {
            generateList(node.children, dataList);
        }
    }
    return dataList;
};

export const SpecificationList: React.FC<IResourceComponentsProps> = () => {
    const [selectedSpecification, setSelectedSpecification] = useState(null);
    const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([]);
    const [searchValue, setSearchValue] = useState('');
    const [autoExpandParent, setAutoExpandParent] = useState(true);
    const [specificationTree, setSpecificationTree] = useState([]);

    const onExpand = (newExpandedKeys: React.Key[]) => {
        setExpandedKeys(newExpandedKeys);
        setAutoExpandParent(false);
    };

    const onTreeSelect = (selectedKeys, info) => {
        let value = info.node.__spec
        if (value == selectedSpecification) value = null;

        setSelectedSpecification(value);
    };

    const onSearchChange = (e) => {
        const {value} = e.target;
        const dataList = generateList(specificationTree);
        console.log(dataList, value, expandedKeys, autoExpandParent, treeData)
        const newExpandedKeys = dataList
            .map(item => {
                if (item.title.indexOf(value) > -1 && value.length > 0) {
                    return getParentKey(item.key, specificationTree);
                }
                return null;
            })
            .filter((item, i, self) => item && self.indexOf(item) === i);
        setExpandedKeys(newExpandedKeys);
        setSearchValue(value);
        setAutoExpandParent(true);
    };

    const {
        data: specifications,
        isLoading: areSpecificationsLoading,
        error: specificationsLoadingError
    } =
        useList({
            resource: 'specifications',
            successNotification(response) {
                setSpecificationTree(buildSpecificationTree(response.data))
            }
        })

    const treeData = useMemo(() => {
        const loop = (data) =>
            data.map(item => {
                const strTitle = item.title;
                const index = strTitle.indexOf(searchValue);

                // const beforeStr = strTitle.substring(0, index);
                // const afterStr = strTitle.slice(index + searchValue.length);
                const title =
                    (index > -1 && searchValue.length > 0) ? (
                        <span className="site-tree-search-value">{strTitle}</span>

                    ) : (
                        <span>{strTitle}</span>
                    );
                if (item.children) {
                    return {...item, title, children: loop(item.children)};
                }

                return {
                    ...item,
                    title,
                };
            });

        return loop(specificationTree);
    }, [searchValue, specificationTree]);

    if (areSpecificationsLoading) {
        return <Spin/>
    }

    if (specificationsLoadingError) {
        return <div>Failed to load specifications: {specificationsLoadingError}</div>
    }


    return (
        <Layout>
            <Content style={{display: 'flex'}}>
                <div style={{width: '300px', padding: '16px'}}>
                    <Search style={{marginBottom: 8}} placeholder="Search" onChange={onSearchChange}/>
                    <Tree
                        onSelect={onTreeSelect}
                        treeData={treeData}
                        expandedKeys={expandedKeys}
                        onExpand={onExpand}
                        autoExpandParent={autoExpandParent}
                        showLine
                        style={{padding: 10}}
                    />
                </div>
                <div style={{flex: 1, padding: '16px'}}>
                    {selectedSpecification ? (
                        <Card title={<Text copyable>{selectedSpecification.id}</Text>}>
                            <Descriptions
                                layout={'horizontal'}
                                column={1}
                                // bordered={true}
                            >
                                {
                                    selectedSpecification.executable.type &&
                                    <Descriptions.Item
                                        label={'Executable Type'}>{selectedSpecification.executable.type}</Descriptions.Item>
                                }

                            </Descriptions>
                        </Card>
                    ) : (
                        <Card title="Select a specification">
                            <p>Nothing</p>
                        </Card>
                    )}
                </div>
            </Content>
        </Layout>
    );
}