import React, {ChangeEvent, ChangeEventHandler, ReactNode, useMemo, useState} from 'react'
import {IResourceComponentsProps, OpenNotificationParams, useList} from '@refinedev/core'
import {Card, Descriptions, Input, Layout, Spin, Tree, Typography} from 'antd'
import './style.css'
import {AntdTree, buildSpecificationTree, convertTreeToAntdTree} from '../utils/specs/tree'
import {Specification} from '../utils/specs/specification'
import {EventDataNode} from 'antd/lib/tree'
import {Key as TableKey} from 'antd/es/table/interface'

const {Content} = Layout
const {Text} = Typography
const {Search} = Input

type SearchResult = {
    key: string,
    title: string,
}

const getParentKey = (key: string, tree: AntdTree[]): string => {
    let parentKey = ''
    for (let i = 0; i < tree.length; i++) {
        const node = tree[i]
        if (node.children) {
            if (node.children.some(item => item.value === key)) {
                parentKey = node.value
            } else if (getParentKey(key, node.children)) {
                parentKey = getParentKey(key, node.children)
            }
        }
    }
    return parentKey
}

const generateList = (data: AntdTree[], dataList: SearchResult[] = []) => {
    for (let i = 0; i < data.length; i++) {
        const node = data[i]
        const {value} = node
        dataList.push({key: value, title: value})
        if (node.children) {
            generateList(node.children, dataList)
        }
    }
    return dataList
}

export const SpecificationList: React.FC<IResourceComponentsProps> = () => {
    const [selectedSpecification, setSelectedSpecification] = useState<Specification | null>(null)
    const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([])
    const [searchValue, setSearchValue] = useState('')
    const [autoExpandParent, setAutoExpandParent] = useState(true)
    const [specificationTree, setSpecificationTree] = useState<AntdTree[]>([])

    const onExpand = (newExpandedKeys: React.Key[]) => {
        setExpandedKeys(newExpandedKeys)
        setAutoExpandParent(false)
    }

    const onTreeSelect = (_: TableKey[], info: { node: EventDataNode<AntdTree> }) => {
        let value = info.node.__spec
        if (value == selectedSpecification) value = null
        if (!value) return

        setSelectedSpecification(value)
    }

    const onSearchChange: ChangeEventHandler<HTMLInputElement> = (event: ChangeEvent<HTMLInputElement>) => {
        const {value} = event.target
        const dataList = generateList(specificationTree)

        console.log('search change:', value, dataList)

        const newExpandedKeys: string[] = dataList
            .map(item => {
                if (item.title.indexOf(value) > -1 && value.length > 0) {
                    return getParentKey(item.key, specificationTree)
                }
                return ''
            })
            .filter((item, i, self) => {
                return item && self.indexOf(item) === i
            })

        setExpandedKeys(newExpandedKeys)
        setSearchValue(value)
        setAutoExpandParent(true)
    }

    const {
        isLoading: areSpecificationsLoading,
        error: specificationsLoadingError,
    } = useList<Specification>({
        resource: 'specifications',
        successNotification(response): (false | OpenNotificationParams) {
            if (!response) return false

            const specs: Specification[] = response.data
            const tree = buildSpecificationTree(specs)
            const antdTree = convertTreeToAntdTree(tree.children)
            setSpecificationTree(antdTree)

            return false
        }
    })

    const treeData = useMemo(() => {
        const loop = (tree: AntdTree[]): AntdTree[] => tree.map(item => {
            let strTitle = ''
            if (typeof item.title === 'string') {
                strTitle = item.title
            } else {
                strTitle = item.title.key || ''
            }
            const index = strTitle.indexOf(searchValue)

            // TODO: refactor this code
            // const beforeStr = strTitle.substring(0, index);
            // const afterStr = strTitle.slice(index + searchValue.length);

            let title: ReactNode
            if (index > -1 && searchValue.length > 0) {
                title = <span className='site-tree-search-value'>{strTitle}</span>
            } else {
                title = <span>{strTitle}</span>
            }

            if (item.children) {
                return {...item, title, children: loop(item.children)}
            }

            return {...item, title}
        })

        return loop(specificationTree)
    }, [searchValue, specificationTree])

    if (areSpecificationsLoading) {
        return <Spin/>
    }

    if (specificationsLoadingError) {
        return <>Failed to load specifications: {specificationsLoadingError}</>
    }


    return (<Layout>
        <Content style={{display: 'flex'}}>
            <div style={{width: '300px', padding: '16px'}}>
                <Search style={{marginBottom: 8}} placeholder='Search' onChange={onSearchChange}/>
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
                {selectedSpecification ? (<Card title={<Text copyable>{selectedSpecification.id}</Text>}>
                    <Descriptions
                        layout={'horizontal'}
                        column={1}
                        // bordered={true}
                    >
                        {selectedSpecification?.executable?.type && <Descriptions.Item
                            label={'Executable Type'}>{selectedSpecification?.executable.type}</Descriptions.Item>}

                    </Descriptions>
                </Card>) : (<Card title='Select a specification'>
                    <p>Nothing</p>
                </Card>)}
            </div>
        </Content>
    </Layout>)
}