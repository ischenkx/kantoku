import {DataProvider} from "@refinedev/core";
import {AxiosInstance} from "axios";
import {
    DefaultApi,
    TasksFilterPostRequest,
    TaskSpecification,
    Configuration
} from "../api/generated";
import {API, ConvertFilter} from "./common"; // Update the import path based on your file structure

export const TaskProvider: DataProvider = {
    getList: async ({resource, pagination, sorters, filters}: {
                        resource?: string;
                        pagination?: { current?: number; pageSize?: number };
                        sorters?: { field: string; order: "asc" | "desc" }[];
                        filters?: Record<string, any>
                    } = {}
    ) => {
        const {current = 1, pageSize = 10} = pagination ?? {};

        const query: TasksFilterPostRequest = {
            filter: ConvertFilter(filters?? []),
            cursor: {
                skip: (current - 1) * pageSize,
                limit: pageSize,
                sort: sorters?.map((sorter) => ({
                    key: sorter.field,
                    ordering: sorter.order === "asc" ? "ASC" : "DESC",
                })),
            },
        };
        const {data: tasks} = await API.tasksFilterPost(query);
        const {data: totalCount} = await API.tasksCountPost(query);

        return {
            data: tasks,
            total: totalCount,
        };
    },

    getOne: async ({resource, id}: {resource: string; id: string}) => {
        const {data: task} = await API.tasksLoadPost([id]);
        return {
            data: task[0],
        };
    },

    create: async (resource: string, {variables}: {
        variables: { inputs: string[]; outputs: string[]; info: Record<string, any> }
    }) => {
        const taskSpecification: TaskSpecification = {
            inputs: variables.inputs,
            outputs: variables.outputs,
            info: variables.info,
        };

        const {data: taskSpawnResponse} = await API.tasksSpawnPost(taskSpecification);

        return {
            data: {id: taskSpawnResponse.id},
        };
    },
};
