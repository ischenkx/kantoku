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


import ApiClient from "../ApiClient";
import Error from '../model/Error';
import Resource from '../model/Resource';
import ResourceInitializer from '../model/ResourceInitializer';
import Specification from '../model/Specification';
import SpecificationBasedTaskParameters from '../model/SpecificationBasedTaskParameters';
import Task from '../model/Task';
import TaskParameters from '../model/TaskParameters';
import TaskSpawnResponse from '../model/TaskSpawnResponse';
import TasksFilterPostRequest from '../model/TasksFilterPostRequest';
import TasksSpecificationsGetPostRequest from '../model/TasksSpecificationsGetPostRequest';
import TasksUpdatePostRequest from '../model/TasksUpdatePostRequest';
import TypeWithID from '../model/TypeWithID';

/**
* Default service.
* @module api/DefaultApi
* @version 1.0.0
*/
export default class DefaultApi {

    /**
    * Constructs a new DefaultApi. 
    * @alias module:api/DefaultApi
    * @class
    * @param {module:ApiClient} [apiClient] Optional API client implementation to use,
    * default to {@link module:ApiClient#instance} if unspecified.
    */
    constructor(apiClient) {
        this.apiClient = apiClient || ApiClient.instance;
    }


    /**
     * Callback function to receive the result of the resourcesAllocatePost operation.
     * @callback module:api/DefaultApi~resourcesAllocatePostCallback
     * @param {String} error Error message, if any.
     * @param {Array.<String>} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Allocates N resources
     * @param {Number} amount 
     * @param {module:api/DefaultApi~resourcesAllocatePostCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link Array.<String>}
     */
    resourcesAllocatePost(amount, callback) {
      let postBody = null;
      // verify the required parameter 'amount' is set
      if (amount === undefined || amount === null) {
        throw new Error("Missing the required parameter 'amount' when calling resourcesAllocatePost");
      }

      let pathParams = {
      };
      let queryParams = {
        'amount': amount
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = [];
      let accepts = ['application/json'];
      let returnType = ['String'];
      return this.apiClient.callApi(
        '/resources/allocate', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the resourcesDeallocatePost operation.
     * @callback module:api/DefaultApi~resourcesDeallocatePostCallback
     * @param {String} error Error message, if any.
     * @param {Object} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Deallocate resources
     * @param {Array.<String>} requestBody A list of resource_db identifiers
     * @param {module:api/DefaultApi~resourcesDeallocatePostCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link Object}
     */
    resourcesDeallocatePost(requestBody, callback) {
      let postBody = requestBody;
      // verify the required parameter 'requestBody' is set
      if (requestBody === undefined || requestBody === null) {
        throw new Error("Missing the required parameter 'requestBody' when calling resourcesDeallocatePost");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = Object;
      return this.apiClient.callApi(
        '/resources/deallocate', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the resourcesInitializePost operation.
     * @callback module:api/DefaultApi~resourcesInitializePostCallback
     * @param {String} error Error message, if any.
     * @param {Object} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Initialize resources
     * @param {Array.<module:model/ResourceInitializer>} resourceInitializer A dictionary (ResourceID -> Value)
     * @param {module:api/DefaultApi~resourcesInitializePostCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link Object}
     */
    resourcesInitializePost(resourceInitializer, callback) {
      let postBody = resourceInitializer;
      // verify the required parameter 'resourceInitializer' is set
      if (resourceInitializer === undefined || resourceInitializer === null) {
        throw new Error("Missing the required parameter 'resourceInitializer' when calling resourcesInitializePost");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = Object;
      return this.apiClient.callApi(
        '/resources/initialize', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the resourcesLoadPost operation.
     * @callback module:api/DefaultApi~resourcesLoadPostCallback
     * @param {String} error Error message, if any.
     * @param {Array.<module:model/Resource>} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Load resources
     * @param {Array.<String>} requestBody A list of resource_db identifiers
     * @param {module:api/DefaultApi~resourcesLoadPostCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link Array.<module:model/Resource>}
     */
    resourcesLoadPost(requestBody, callback) {
      let postBody = requestBody;
      // verify the required parameter 'requestBody' is set
      if (requestBody === undefined || requestBody === null) {
        throw new Error("Missing the required parameter 'requestBody' when calling resourcesLoadPost");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = [Resource];
      return this.apiClient.callApi(
        '/resources/load', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the tasksCountPost operation.
     * @callback module:api/DefaultApi~tasksCountPostCallback
     * @param {String} error Error message, if any.
     * @param {Number} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Count records using a filter
     * @param {module:model/TasksFilterPostRequest} tasksFilterPostRequest A query
     * @param {module:api/DefaultApi~tasksCountPostCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link Number}
     */
    tasksCountPost(tasksFilterPostRequest, callback) {
      let postBody = tasksFilterPostRequest;
      // verify the required parameter 'tasksFilterPostRequest' is set
      if (tasksFilterPostRequest === undefined || tasksFilterPostRequest === null) {
        throw new Error("Missing the required parameter 'tasksFilterPostRequest' when calling tasksCountPost");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = 'Number';
      return this.apiClient.callApi(
        '/tasks/count', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the tasksFilterPost operation.
     * @callback module:api/DefaultApi~tasksFilterPostCallback
     * @param {String} error Error message, if any.
     * @param {Array.<module:model/Task>} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Load records using a filter
     * @param {module:model/TasksFilterPostRequest} tasksFilterPostRequest A query
     * @param {module:api/DefaultApi~tasksFilterPostCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link Array.<module:model/Task>}
     */
    tasksFilterPost(tasksFilterPostRequest, callback) {
      let postBody = tasksFilterPostRequest;
      // verify the required parameter 'tasksFilterPostRequest' is set
      if (tasksFilterPostRequest === undefined || tasksFilterPostRequest === null) {
        throw new Error("Missing the required parameter 'tasksFilterPostRequest' when calling tasksFilterPost");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = [Task];
      return this.apiClient.callApi(
        '/tasks/filter', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the tasksLoadPost operation.
     * @callback module:api/DefaultApi~tasksLoadPostCallback
     * @param {String} error Error message, if any.
     * @param {Array.<module:model/Task>} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Load a set of tasks
     * @param {Array.<String>} requestBody A list of task identifiers
     * @param {module:api/DefaultApi~tasksLoadPostCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link Array.<module:model/Task>}
     */
    tasksLoadPost(requestBody, callback) {
      let postBody = requestBody;
      // verify the required parameter 'requestBody' is set
      if (requestBody === undefined || requestBody === null) {
        throw new Error("Missing the required parameter 'requestBody' when calling tasksLoadPost");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = [Task];
      return this.apiClient.callApi(
        '/tasks/load', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the tasksSpawnFromSpecPost operation.
     * @callback module:api/DefaultApi~tasksSpawnFromSpecPostCallback
     * @param {String} error Error message, if any.
     * @param {module:model/TaskSpawnResponse} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Spawn a new task from specification
     * @param {module:model/SpecificationBasedTaskParameters} specificationBasedTaskParameters The specification of a task to be spawned
     * @param {module:api/DefaultApi~tasksSpawnFromSpecPostCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/TaskSpawnResponse}
     */
    tasksSpawnFromSpecPost(specificationBasedTaskParameters, callback) {
      let postBody = specificationBasedTaskParameters;
      // verify the required parameter 'specificationBasedTaskParameters' is set
      if (specificationBasedTaskParameters === undefined || specificationBasedTaskParameters === null) {
        throw new Error("Missing the required parameter 'specificationBasedTaskParameters' when calling tasksSpawnFromSpecPost");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = TaskSpawnResponse;
      return this.apiClient.callApi(
        '/tasks/spawn_from_spec', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the tasksSpawnPost operation.
     * @callback module:api/DefaultApi~tasksSpawnPostCallback
     * @param {String} error Error message, if any.
     * @param {module:model/TaskSpawnResponse} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Spawn a new task
     * @param {module:model/TaskParameters} taskParameters The specification of a task to be spawned
     * @param {module:api/DefaultApi~tasksSpawnPostCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/TaskSpawnResponse}
     */
    tasksSpawnPost(taskParameters, callback) {
      let postBody = taskParameters;
      // verify the required parameter 'taskParameters' is set
      if (taskParameters === undefined || taskParameters === null) {
        throw new Error("Missing the required parameter 'taskParameters' when calling tasksSpawnPost");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = TaskSpawnResponse;
      return this.apiClient.callApi(
        '/tasks/spawn', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the tasksSpecificationsCreatePost operation.
     * @callback module:api/DefaultApi~tasksSpecificationsCreatePostCallback
     * @param {String} error Error message, if any.
     * @param data This operation does not return a value.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Create a specification
     * @param {module:model/Specification} specification 
     * @param {module:api/DefaultApi~tasksSpecificationsCreatePostCallback} callback The callback function, accepting three arguments: error, data, response
     */
    tasksSpecificationsCreatePost(specification, callback) {
      let postBody = specification;
      // verify the required parameter 'specification' is set
      if (specification === undefined || specification === null) {
        throw new Error("Missing the required parameter 'specification' when calling tasksSpecificationsCreatePost");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = null;
      return this.apiClient.callApi(
        '/tasks/specifications/create', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the tasksSpecificationsGetAllPost operation.
     * @callback module:api/DefaultApi~tasksSpecificationsGetAllPostCallback
     * @param {String} error Error message, if any.
     * @param {Array.<module:model/Specification>} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Get all specifications
     * @param {module:api/DefaultApi~tasksSpecificationsGetAllPostCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link Array.<module:model/Specification>}
     */
    tasksSpecificationsGetAllPost(callback) {
      let postBody = null;

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = [];
      let accepts = ['application/json'];
      let returnType = [Specification];
      return this.apiClient.callApi(
        '/tasks/specifications/get_all', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the tasksSpecificationsGetPost operation.
     * @callback module:api/DefaultApi~tasksSpecificationsGetPostCallback
     * @param {String} error Error message, if any.
     * @param {module:model/Specification} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Get specifications by id
     * @param {module:model/TasksSpecificationsGetPostRequest} tasksSpecificationsGetPostRequest 
     * @param {module:api/DefaultApi~tasksSpecificationsGetPostCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/Specification}
     */
    tasksSpecificationsGetPost(tasksSpecificationsGetPostRequest, callback) {
      let postBody = tasksSpecificationsGetPostRequest;
      // verify the required parameter 'tasksSpecificationsGetPostRequest' is set
      if (tasksSpecificationsGetPostRequest === undefined || tasksSpecificationsGetPostRequest === null) {
        throw new Error("Missing the required parameter 'tasksSpecificationsGetPostRequest' when calling tasksSpecificationsGetPost");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = Specification;
      return this.apiClient.callApi(
        '/tasks/specifications/get', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the tasksSpecificationsRemovePost operation.
     * @callback module:api/DefaultApi~tasksSpecificationsRemovePostCallback
     * @param {String} error Error message, if any.
     * @param data This operation does not return a value.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Remove a specification
     * @param {module:model/TasksSpecificationsGetPostRequest} tasksSpecificationsGetPostRequest 
     * @param {module:api/DefaultApi~tasksSpecificationsRemovePostCallback} callback The callback function, accepting three arguments: error, data, response
     */
    tasksSpecificationsRemovePost(tasksSpecificationsGetPostRequest, callback) {
      let postBody = tasksSpecificationsGetPostRequest;
      // verify the required parameter 'tasksSpecificationsGetPostRequest' is set
      if (tasksSpecificationsGetPostRequest === undefined || tasksSpecificationsGetPostRequest === null) {
        throw new Error("Missing the required parameter 'tasksSpecificationsGetPostRequest' when calling tasksSpecificationsRemovePost");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = null;
      return this.apiClient.callApi(
        '/tasks/specifications/remove', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the tasksSpecificationsTypesCreatePost operation.
     * @callback module:api/DefaultApi~tasksSpecificationsTypesCreatePostCallback
     * @param {String} error Error message, if any.
     * @param data This operation does not return a value.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Create a type
     * @param {module:model/TypeWithID} typeWithID 
     * @param {module:api/DefaultApi~tasksSpecificationsTypesCreatePostCallback} callback The callback function, accepting three arguments: error, data, response
     */
    tasksSpecificationsTypesCreatePost(typeWithID, callback) {
      let postBody = typeWithID;
      // verify the required parameter 'typeWithID' is set
      if (typeWithID === undefined || typeWithID === null) {
        throw new Error("Missing the required parameter 'typeWithID' when calling tasksSpecificationsTypesCreatePost");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = null;
      return this.apiClient.callApi(
        '/tasks/specifications/types/create', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the tasksSpecificationsTypesGetAllPost operation.
     * @callback module:api/DefaultApi~tasksSpecificationsTypesGetAllPostCallback
     * @param {String} error Error message, if any.
     * @param {Array.<module:model/TypeWithID>} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Get all types
     * @param {module:api/DefaultApi~tasksSpecificationsTypesGetAllPostCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link Array.<module:model/TypeWithID>}
     */
    tasksSpecificationsTypesGetAllPost(callback) {
      let postBody = null;

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = [];
      let accepts = ['application/json'];
      let returnType = [TypeWithID];
      return this.apiClient.callApi(
        '/tasks/specifications/types/get_all', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the tasksSpecificationsTypesGetPost operation.
     * @callback module:api/DefaultApi~tasksSpecificationsTypesGetPostCallback
     * @param {String} error Error message, if any.
     * @param {module:model/TypeWithID} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Get a type by id
     * @param {module:model/TasksSpecificationsGetPostRequest} tasksSpecificationsGetPostRequest 
     * @param {module:api/DefaultApi~tasksSpecificationsTypesGetPostCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/TypeWithID}
     */
    tasksSpecificationsTypesGetPost(tasksSpecificationsGetPostRequest, callback) {
      let postBody = tasksSpecificationsGetPostRequest;
      // verify the required parameter 'tasksSpecificationsGetPostRequest' is set
      if (tasksSpecificationsGetPostRequest === undefined || tasksSpecificationsGetPostRequest === null) {
        throw new Error("Missing the required parameter 'tasksSpecificationsGetPostRequest' when calling tasksSpecificationsTypesGetPost");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = TypeWithID;
      return this.apiClient.callApi(
        '/tasks/specifications/types/get', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the tasksSpecificationsTypesRemovePost operation.
     * @callback module:api/DefaultApi~tasksSpecificationsTypesRemovePostCallback
     * @param {String} error Error message, if any.
     * @param data This operation does not return a value.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Remove a type
     * @param {module:model/TasksSpecificationsGetPostRequest} tasksSpecificationsGetPostRequest 
     * @param {module:api/DefaultApi~tasksSpecificationsTypesRemovePostCallback} callback The callback function, accepting three arguments: error, data, response
     */
    tasksSpecificationsTypesRemovePost(tasksSpecificationsGetPostRequest, callback) {
      let postBody = tasksSpecificationsGetPostRequest;
      // verify the required parameter 'tasksSpecificationsGetPostRequest' is set
      if (tasksSpecificationsGetPostRequest === undefined || tasksSpecificationsGetPostRequest === null) {
        throw new Error("Missing the required parameter 'tasksSpecificationsGetPostRequest' when calling tasksSpecificationsTypesRemovePost");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = null;
      return this.apiClient.callApi(
        '/tasks/specifications/types/remove', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }

    /**
     * Callback function to receive the result of the tasksUpdatePost operation.
     * @callback module:api/DefaultApi~tasksUpdatePostCallback
     * @param {String} error Error message, if any.
     * @param {Object} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Update a record
     * @param {module:model/TasksUpdatePostRequest} tasksUpdatePostRequest 
     * @param {module:api/DefaultApi~tasksUpdatePostCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link Object}
     */
    tasksUpdatePost(tasksUpdatePostRequest, callback) {
      let postBody = tasksUpdatePostRequest;
      // verify the required parameter 'tasksUpdatePostRequest' is set
      if (tasksUpdatePostRequest === undefined || tasksUpdatePostRequest === null) {
        throw new Error("Missing the required parameter 'tasksUpdatePostRequest' when calling tasksUpdatePost");
      }

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = Object;
      return this.apiClient.callApi(
        '/tasks/update', 'POST',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }


}
