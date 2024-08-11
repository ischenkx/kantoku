import {DataProvider} from '@refinedev/core'
import {API, filtersToMongoFilter} from './common' // Update the import path based on your file structure

export const TaskProvider: DataProvider = {
    getList: async ({resource, pagination, sorters, filters}: {
                        resource?: string;
                        pagination?: { current?: number; pageSize?: number };
                        sorters?: { field: string; order: 'asc' | 'desc' }[];
                        filters?: Record<string, any>
                    } = {}
    ) => {
        console.log('Params:', resource, pagination, sorters, filters)
        const {current = 1, pageSize = 10} = pagination ?? {}

        sorters ||= []

        if (sorters.length === 0) {
            sorters = [
                {
                    field: 'info.updated_at',
                    order: 'desc',
                }
            ]
        }

        const findCommand = {
            operation: 'find',
            params: [
                {name: 'filter', value: filtersToMongoFilter(filters) ?? []},
                {name: 'skip', value: (current - 1) * pageSize},
                {name: 'limit', value: pageSize},
                {
                    name: 'sort',
                    value: Object.assign({}, ...sorters?.map((sorter) => ({
                        [sorter.field]: sorter.order === 'asc' ? 1 : -1,
                    })))
                }
            ]
        }

        const countCommand = {
            operation: 'count',
            params: [
                {name: 'query', value: filtersToMongoFilter(filters) ?? []},
            ]
        }

        const {data: tasks} = await API.tasksStorageExecPost(findCommand)
        const {data: totalCount} = await API.tasksStorageExecPost(countCommand)

        return {
            data: tasks[0]?.cursor?.firstBatch || [],
            total: totalCount[0]?.n || 0,
        }
    },

    getOne: async ({resource, id}: { resource: string; id: string }) => {
        const {data: task} = await API.tasksLoadPost([id])
        return {
            data: task[0],
        }
    },

    create: async ({resource, variables}: {
        resource: string,
        variables: { specification: string, parameters: any, info: any }
    }) => {
        const taskSpecification = {
            specification: variables.specification,
            parameters: variables.parameters,
            info: variables.info,
        }

        console.log('spawning:', resource, variables, taskSpecification)

        const {data: taskSpawnResponse} = await API.tasksSpawnFromSpecPost(taskSpecification)

        return {
            data: {id: taskSpawnResponse.id},
        }
    },
}
