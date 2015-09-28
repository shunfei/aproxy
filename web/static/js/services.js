'use strict';

/* Services */

var aproxyServices = angular.module('aproxyServices', ['ngResource']);

aproxyServices.factory('BackendConf', ['$resource',
  function($resource){
    return $resource('api/backends/:hostname', {}, {
      query: {method:'GET', params:{hostname:'all'}},
      remove: {method:'DELETE'},
      update: {method:'PUT'}
    });
  }]);

aproxyServices.factory('Role', ['$resource',
  function($resource){
    return $resource('api/role/:id', {}, {
      query: {method:'GET', params:{id:'all'}},
      remove: {method:'DELETE'},
      update: {method:'PUT'}
    });
  }]);


aproxyServices.factory('Authority', ['$resource',
  function($resource){
    return $resource('api/authority/:id', {}, {
      query: {method:'GET', params:{id:'all'}},
      remove: {method:'DELETE'},
      update: {method:'PUT'}
    });
  }]);

aproxyServices.factory('User', ['$resource',
  function($resource){
    return $resource('api/users/:email', {}, {
      query: {method:'GET', params:{email:'all'}},
      remove: {method:'DELETE'},
      update: {method:'PUT'}
    });
  }]);