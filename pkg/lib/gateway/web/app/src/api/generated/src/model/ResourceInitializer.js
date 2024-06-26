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
 * The ResourceInitializer model module.
 * @module model/ResourceInitializer
 * @version 1.0.0
 */
class ResourceInitializer {
    /**
     * Constructs a new <code>ResourceInitializer</code>.
     * @alias module:model/ResourceInitializer
     * @param id {String} 
     * @param value {String} 
     */
    constructor(id, value) { 
        
        ResourceInitializer.initialize(this, id, value);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, id, value) { 
        obj['id'] = id;
        obj['value'] = value;
    }

    /**
     * Constructs a <code>ResourceInitializer</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/ResourceInitializer} obj Optional instance to populate.
     * @return {module:model/ResourceInitializer} The populated <code>ResourceInitializer</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new ResourceInitializer();

            if (data.hasOwnProperty('id')) {
                obj['id'] = ApiClient.convertToType(data['id'], 'String');
            }
            if (data.hasOwnProperty('value')) {
                obj['value'] = ApiClient.convertToType(data['value'], 'String');
            }
        }
        return obj;
    }


}

/**
 * @member {String} id
 */
ResourceInitializer.prototype['id'] = undefined;

/**
 * @member {String} value
 */
ResourceInitializer.prototype['value'] = undefined;






export default ResourceInitializer;

