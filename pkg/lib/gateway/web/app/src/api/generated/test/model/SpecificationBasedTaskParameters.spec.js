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

(function(root, factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD.
    define(['expect.js', process.cwd()+'/src/index'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    factory(require('expect.js'), require(process.cwd()+'/src/index'));
  } else {
    // Browser globals (root is window)
    factory(root.expect, root.KantokuApi);
  }
}(this, function(expect, KantokuApi) {
  'use strict';

  var instance;

  beforeEach(function() {
    instance = new KantokuApi.SpecificationBasedTaskParameters();
  });

  var getProperty = function(object, getter, property) {
    // Use getter method if present; otherwise, get the property directly.
    if (typeof object[getter] === 'function')
      return object[getter]();
    else
      return object[property];
  }

  var setProperty = function(object, setter, property, value) {
    // Use setter method if present; otherwise, set the property directly.
    if (typeof object[setter] === 'function')
      object[setter](value);
    else
      object[property] = value;
  }

  describe('SpecificationBasedTaskParameters', function() {
    it('should create an instance of SpecificationBasedTaskParameters', function() {
      // uncomment below and update the code to test SpecificationBasedTaskParameters
      //var instance = new KantokuApi.SpecificationBasedTaskParameters();
      //expect(instance).to.be.a(KantokuApi.SpecificationBasedTaskParameters);
    });

    it('should have the property parameters (base name: "parameters")', function() {
      // uncomment below and update the code to test the property parameters
      //var instance = new KantokuApi.SpecificationBasedTaskParameters();
      //expect(instance).to.be();
    });

    it('should have the property specification (base name: "specification")', function() {
      // uncomment below and update the code to test the property specification
      //var instance = new KantokuApi.SpecificationBasedTaskParameters();
      //expect(instance).to.be();
    });

    it('should have the property info (base name: "info")', function() {
      // uncomment below and update the code to test the property info
      //var instance = new KantokuApi.SpecificationBasedTaskParameters();
      //expect(instance).to.be();
    });

  });

}));
