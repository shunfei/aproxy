'use strict';

/* App Module */

var aproxyApp = angular.module('aproxyApp', [
  'ngRoute',
  'aproxyControllers',
  'aproxyFilters',
  'aproxyServices'
]);

aproxyApp.config(['$httpProvider', function($httpProvider) {
    $httpProvider.defaults.headers.common["X-Requested-With"] = 'XMLHttpRequest';
}]);

aproxyApp.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/backends', {
        templateUrl: 'partials/backend-conf/list.html',
        controller: 'BackendConfListCtrl'
      }).
      when('/backend-conf/new', {
        templateUrl: 'partials/backend-conf/new.html',
        controller: 'BackendConfAddNewCtrl'
      }).
      when('/backend-conf/:hostname', {
        templateUrl: 'partials/backend-conf/detail.html',
        controller: 'BackendConfDetailCtrl'
      }).
      otherwise({
        redirectTo: '/backends'
      });

    $routeProvider.
      when('/role', {
        templateUrl: 'partials/role/list.html',
        controller: 'RoleListCtrl'
      }).
      when('/role/new', {
        templateUrl: 'partials/role/new.html',
        controller: 'RoleAddNewCtrl'
      }).
      when('/role/:id', {
        templateUrl: 'partials/role/detail.html',
        controller: 'RoleDetailCtrl'
      });

    $routeProvider.
      when('/authority', {
        templateUrl: 'partials/authority/list.html',
        controller: 'AuthorityListCtrl'
      }).
      when('/authority/new', {
        templateUrl: 'partials/authority/new.html',
        controller: 'AuthorityAddNewCtrl'
      }).
      when('/authority/:id', {
        templateUrl: 'partials/authority/detail.html',
        controller: 'AuthorityDetailCtrl'
      });

    $routeProvider.
      when('/users', {
        templateUrl: 'partials/users/list.html',
        controller: 'UserListCtrl'
      }).
      when('/users/new', {
        templateUrl: 'partials/users/new.html',
        controller: 'UserAddNewCtrl'
      }).
      when('/users/:email', {
        templateUrl: 'partials/users/detail.html',
        controller: 'UserDetailCtrl'
      });

  }]);
