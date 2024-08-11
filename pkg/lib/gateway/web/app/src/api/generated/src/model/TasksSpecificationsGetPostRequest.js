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
 * The TasksSpecificationsGetPostRequest model module.
 * @module model/TasksSpecificationsGetPostRequest
 * @version 1.0.0
 */
class TasksSpecificationsGetPostRequest {
    /**
     * Constructs a new <code>TasksSpecificationsGetPostRequest</code>.
     * @alias module:model/TasksSpecificationsGetPostRequest
     * @param id {String} 
     */
    constructor(id) { 
        
        TasksSpecificationsGetPostRequest.initialize(this, id);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, id) { 
        obj['id'] = id;
    }

    /**
     * Constructs a <code>TasksSpecificationsGetPostRequest</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/TasksSpecificationsGetPostRequest} obj Optional instance to populate.
     * @return {module:model/TasksSpecificationsGetPostRequest} The populated <code>TasksSpecificationsGetPostRequest</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new TasksSpecificationsGetPostRequest();

            if (data.hasOwnProperty('id')) {
                obj['id'] = ApiClient.convertToType(data['id'], 'String');
            }
        }
        return obj;
    }


}

/**
 * @member {String} id
 */
TasksSpecificationsGetPostRequest.prototype['id'] = undefined;






export default TasksSpecificationsGetPostRequest;
