import {Configuration, DefaultApi, Task} from './generated'
import {AxiosResponse} from 'axios'

export class Api {
    rawApi: DefaultApi

    constructor(public readonly url: string) {
        this.rawApi = new DefaultApi(new Configuration(), url)
    }

    getRaw(): DefaultApi {
        return this.rawApi
    }

    async getTaskRestarts(task: Task): Promise<Task[]> {
        const info = task.info as Record<string, any>;

        if (!info['restart_root']) return []

        const rootId: string = info['restart_root']

        const result: AxiosResponse<Task[]> = await this.getRaw().tasksStorageGetWithPropertiesPost({
            'properties_to_values': {
                'info.restart_root': [rootId]
            }
        })

        if (result.status != 200) {
            throw new Error(result.statusText)
        }

        return result.data || []
    }
}