import {DataProvider} from "@refinedev/core";
import {
    TasksFilterPostRequest,
    TaskSpecification
} from "../api/generated";
import {API, ConvertFilter} from "./common";

export const SpecificationProvider: DataProvider = {
    getList: async ({resource}: {
                        resource?: string;
                    } = {}
    ) => {
        if (resource === 'specification') {
            return await API.tasksSpecificationsGetAllPost();
        } else {
            return await API.tasksSpecificationsTypesGetAllPost();
        }
    },

    getOne: async ({resource, id}: { resource: string; id: string }) => {
        if (resource === 'specification') {
            return await API.tasksSpecificationsGetPost({id: id});
        } else {
            return await API.tasksSpecificationsTypesGetPost({id: id});
        }
    },
};
