/**
 * Kantoku API
 * Create and execute distributed workflows
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 *
 */

import ApiClient from '../ApiClient';

/**
 * The TasksUpdatePostRequest model module.
 * @module model/TasksUpdatePostRequest
 * @version 1.0.0
 */
class TasksUpdatePostRequest {
    /**
     * Constructs a new <code>TasksUpdatePostRequest</code>.
     * @alias module:model/TasksUpdatePostRequest
     * @param filter {Object} 
     * @param update {Object} 
     */
    constructor(filter, update) { 
        
        TasksUpdatePostRequest.initialize(this, filter, update);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, filter, update) { 
        obj['filter'] = filter;
        obj['update'] = update;
    }

    /**
     * Constructs a <code>TasksUpdatePostRequest</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/TasksUpdatePostRequest} obj Optional instance to populate.
     * @return {module:model/TasksUpdatePostRequest} The populated <code>TasksUpdatePostRequest</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new TasksUpdatePostRequest();

            if (data.hasOwnProperty('filter')) {
                obj['filter'] = ApiClient.convertToType(data['filter'], Object);
            }
            if (data.hasOwnProperty('update')) {
                obj['update'] = ApiClient.convertToType(data['update'], Object);
            }
            if (data.hasOwnProperty('upsert')) {
                obj['upsert'] = ApiClient.convertToType(data['upsert'], Object);
            }
        }
        return obj;
    }


}

/**
 * @member {Object} filter
 */
TasksUpdatePostRequest.prototype['filter'] = undefined;

/**
 * @member {Object} update
 */
TasksUpdatePostRequest.prototype['update'] = undefined;

/**
 * @member {Object} upsert
 */
TasksUpdatePostRequest.prototype['upsert'] = undefined;






export default TasksUpdatePostRequest;

