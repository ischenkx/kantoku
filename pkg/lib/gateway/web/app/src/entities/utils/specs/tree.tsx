import {Specification} from './specification'
import {ReactElement} from 'react'

export type TreeValue = {
    path: string
    spec: Specification | null
}

export type TreeChildren = Record<string, Tree>

export type Tree = {
    value: TreeValue
    children: TreeChildren
}

export const newTree = (path: string, spec: Specification | null): Tree => {
    return {value: {path, spec}, children: {}}
}

function compressSpecificationTree(tree: Tree) {
    Object.values(tree.children).forEach(compressSpecificationTree)

    while (objectSize(tree.children) === 1) {
        const [childKey, child] = Object.entries(tree.children)[0]

        const childChildren = Object.entries(child.children)

        if (childChildren.length !== 1) break

        tree.children = Object.fromEntries(
            childChildren.map(([subChildKey, subChild]) => {
                return [childKey + '/' + subChildKey, subChild]
            })
        )
        tree.value = child.value
    }

}

export function buildSpecificationTree(specifications: Specification[]): Tree {
    const tree = newTree('', null)

    for (const spec of specifications) {
        const path = spec.id.split('/').filter(part => !!part)
        let currentNode = tree
        let currentPath = ''
        for (const part of path) {
            if (currentPath.length > 0) currentPath += '/'
            currentPath += part

            if (!currentNode.children[part]) currentNode.children[part] = newTree(currentPath, null)

            currentNode = currentNode.children[part]
        }

        currentNode.value.spec = spec
    }

    console.log(tree)

    compressSpecificationTree(tree)

    console.log(tree)

    return tree
}

export type AntdTree = {
    value: string,
    key: string,
    title: string | ReactElement,
    selectable: boolean,
    children: AntdTree[],
    __spec: Specification | null,
}

export function convertTreeToAntdTree(headlessTree: TreeChildren): AntdTree[] {
    return Object.entries(headlessTree).map(([key, subTree]) => {
        return {
            value: subTree.value.path,
            key: subTree.value.path,
            title: key,
            selectable: Object.keys(subTree.children).length === 0,
            children: convertTreeToAntdTree(subTree.children),
            __spec: subTree.value.spec,
        }
    })
}

function objectSize(obj) {
    let length = 0
    for (const key in obj) {
        length++
    }

    return length
}