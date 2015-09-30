'use strict';

var authTypes = {
            0: "Public",
            1: "Need Login",
            2: "Need Authority"
          };

var adminLevel = {
            0: "Not Administrator",
            50: "System Administrator",
            99: "Super Administrator"
          };

/* Controllers */

var aproxyControllers = angular.module('aproxyControllers', ['checklist-model']);

/*
list backend-config
 */
aproxyControllers.controller('BackendConfListCtrl', ['$scope', 'BackendConf',
  function($scope, BackendConf) {
    BackendConf.query(null, function (res) {
      if (res.success) {
        $scope.backends = res.data;
        $scope.orderProp = 'hostname';
      }else{
        alert('Get BackendConf list error: ' + res.error);
      }
    });

    $scope.delBackendConf = function(bc) {
      if (!confirm('DEL BackendConf for ' + bc.HostName + " ?")) { return }
      BackendConf.remove({id: bc.Id}, function (res) {
        if (res.success === true) {
          $("#bc-"+bc.Id).remove();
        } else {
          alert(res.error);
        }
      });
    };
    
  }]);

/*
edit backend-config
 */
aproxyControllers.controller('BackendConfDetailCtrl', ['$scope', '$location', '$routeParams', 'BackendConf',
  function($scope, $location, $routeParams, BackendConf) {
    $scope.authTypes = authTypes;
    BackendConf.get({hostname: $routeParams.hostname}, function(res) {
      if (res.success) {
        $scope.backend = res.data;
        if (res.data) {
            $scope.oldHostName = res.data.HostName;
        }
      }else{
        alert('Get BackendConf error: ' + res.error);
      }
    });

    $scope.saveBackendConf = function () {
      var bc = $scope.backend;
      if (bc.AuthType === undefined) {
        alert('please select auth-type.');
        return;
      }
      if (!bc.HostName || !bc.UpStreams || bc.UpStreams.length < 1) {
        alert('please fill all items.');
        return;
      }
      $scope.isSaving = true;
      BackendConf.update({id: bc.Id, hostname: $scope.oldHostName}, bc, 
        function (res) {
          if (res.success === true) {
            $location.path("/backend-conf/" + res.data.HostName);
          } else {
            alert(res.error);
          }
          $scope.isSaving = false;
        });
    }
  }]);

/*
add new backend-config
 */
aproxyControllers.controller('BackendConfAddNewCtrl', ['$scope', '$location', 'BackendConf',
  function($scope, $location, BackendConf) {
    $scope.authTypes = authTypes;
    $scope.backend = new BackendConf();

    $scope.saveBackendConf = function () {
      var bc = angular.copy($scope.backend);
      if (bc.AuthType === undefined) {
        alert('please select auth-type.');
        return;
      }
      if (!bc.HostName || !bc.UpStreams || bc.UpStreams.length < 1) {
        alert('please fill all items.');
        return;
      }
      $scope.isSaving = true;
      bc.$save(null, function (res) {
        if (res.success) {
          $location.path("/backend-conf/" + res.data.HostName);
        } else {
          alert(res.error);
        }
        $scope.isSaving = false;
      });
    }
  }]);




/*
list role
 */
aproxyControllers.controller('RoleListCtrl', ['$scope', 'Role',
  function($scope, Role) {
    Role.query(null, function (res) {
      if (res.success) {
        $scope.roles = res.data;
        $scope.orderProp = 'Name';
      }else{
        alert('Get role list error: ' + res.error);
      }
    });

    $scope.delRole = function(role) {
      if (!confirm('DEL Role ' + role.Name + " ?")) { return }
      Role.remove({id: role.Id}, function (res) {
        if (res.success === true) {
          $("#role-"+role.Id).remove();
        } else {
          alert(res.error);
        }
      });
    };
    
  }]);

/*
edit role
 */
aproxyControllers.controller('RoleDetailCtrl', ['$scope', '$location', '$routeParams', 'Role',
  function($scope, $location, $routeParams, Role) {
    Role.get({id: $routeParams.id}, function(res) {
      if (res.success) {
        $scope.role = res.data;
      }else{
        alert('Get role error: ' + res.error);
      }
    });

    $scope.saveRoleConf = function () {
      var role = $scope.role;
      var noAllow = !role.Allow || role.Allow.length < 1;
      var noDeny = !role.Deny || role.Deny.length < 1;
      if (!role.Name || (noAllow && noDeny)) {
        alert('please fill all items.');
        return;
      }
      $scope.isSaving = true;
      Role.update({id: role.Id}, role, 
        function (res) {
          if (res.success === true) {
            $location.path("/role/" + res.data.Id);
          } else {
            alert(res.error);
          }
          $scope.isSaving = false;
        });
    }
  }]);

/*
add new role
 */
aproxyControllers.controller('RoleAddNewCtrl', ['$scope', '$location', 'Role',
  function($scope, $location, Role) {
    $scope.role = new Role();

    $scope.saveRoleConf = function () {
      var role = angular.copy($scope.role);
      var noAllow = !role.Allow || role.Allow.length < 1;
      var noDeny = !role.Deny || role.Deny.length < 1;
      if (!role.Name || (noAllow && noDeny)) {
        alert('please fill all items.');
        return;
      }
      $scope.isSaving = true;
      role.$save(null, function (res) {
        if (res.success) {
          $location.path("/role/" + res.data.Id);
        } else {
          alert(res.error);
        }
        $scope.isSaving = false;
      });
    }
  }]);




/*
list authority
 */
aproxyControllers.controller('AuthorityListCtrl', ['$scope', '$location', 'Authority', 
  function($scope, $location, Authority) {
    var promise = Authority.query(null, function (res) {
      if (res.success) {
        $scope.authorities = res.data;
        $scope.orderProp = 'Email';
      }else{
        alert('Get Authority list error: ' + res.error);
      }
    });

    $scope.emailToEdit = "";

    $scope.delAuthority = function(athy) {
      if (!confirm('DEL Authority for ' + athy.Email + " ?")) { return }
      Authority.remove({id: athy.Id}, function (res) {
        if (res.success === true) {
          $("#ahty-"+athy.Id).remove();
        } else {
          alert(res.error);
        }
      });
    };

    $scope.AddOrEdit = function() {
      if ($scope.emailToEdit == "") {
        $location.path("/authority/new");
      } else {
        Authority.get({email: $scope.emailToEdit}, function(res) {
          if (res.success && res.data && res.data.Id) {
            $location.path("/authority/" + res.data.Id);
          }else{
            $location.path("/authority/new").search({email: $scope.emailToEdit});
          }
        });
      }
    };
    
  }]);

/*
edit authority
 */
aproxyControllers.controller('AuthorityDetailCtrl', 
  ['$scope', '$location', '$routeParams', 'Authority', 'Role',
  function($scope, $location, $routeParams, Authority, Role) {
    $scope.adminLevel = adminLevel;
    Authority.get({id: $routeParams.id}, function(res) {
      if (res.success) {
        $scope.authority = res.data;
      }else{
        alert('Get authority error: ' + res.error);
      }
    });
    Role.query(null, function (res) {
      if (res.success) {
        $scope.roles = res.data;
        $scope.orderProp = 'Name';
      }else{
        alert('Get Role list error: ' + res.error);
      }
    });

    $scope.saveAuthorityConf = function () {
      var ahty = angular.copy($scope.authority);
      var noAllow = !ahty.Allow || ahty.Allow.length < 1;
      var noDeny = !ahty.Deny || ahty.Deny.length < 1;
      var noRoles = !ahty.Roles || ahty.Roles.length < 1;
      if (noAllow && noDeny && noRoles) {
        alert('please fill all items.');
        return;
      }
      $scope.isSaving = true;
      Authority.update({id: ahty.Id}, ahty, 
        function (res) {
          if (res.success === true) {
            $location.path("/authority/" + res.data.Id);
          } else {
            alert(res.error);
          }
          $scope.isSaving = false;
        });
    }
  }]);

/*
add new authority
 */
aproxyControllers.controller('AuthorityAddNewCtrl', 
  ['$scope', '$location', 'Authority', 'Role',
  function($scope, $location, Authority, Role) {
    $scope.adminLevel = adminLevel;
    $scope.authority = new Authority();
    $scope.authority.AdminLevel = 0
    $scope.authority.Email = $location.search().email || '';
    Role.query(null, function (res) {
      if (res.success) {
        $scope.roles = res.data;
        $scope.orderProp = 'Name';
      }else{
        alert('Get Role list error: ' + res.error);
      }
    });

    $scope.saveAuthorityConf = function () {
      var ahty = angular.copy($scope.authority);
      var noAllow = !ahty.Allow || ahty.Allow.length < 1;
      var noDeny = !ahty.Deny || ahty.Deny.length < 1;
      var noRoles = !ahty.Roles || ahty.Roles.length < 1;
      if (!ahty.Email || (noAllow && noDeny && noRoles)) {
        alert('please fill all items.');
        return;
      }
      $scope.isSaving = true;
      ahty.$save(null, function (res) {
        if (res.success) {
          $location.path("/authority/" + res.data.Id);
        } else {
          alert(res.error);
        }
        $scope.isSaving = false;
      });
    };

  }]);



/*
list user
 */
aproxyControllers.controller('UserListCtrl', ['$scope', 'User',
  function($scope, User) {
    User.query(null, function (res) {
      if (res.success) {
        $scope.users = res.data;
        $scope.orderProp = 'Email';
      }else{
        alert('Get User list error: ' + res.error);
      }
    });
    
    $scope.delUser = function(user) {
      if (!confirm('DEL User ' + user.Email + " ?")) { return }
      User.remove({id: user.Id}, function (res) {
        if (res.success === true) {
          $("#user-"+user.Id).remove();
        } else {
          alert(res.error);
        }
      });
    };

  }]);

/*
edit user
 */
aproxyControllers.controller('UserDetailCtrl', 
  ['$scope', '$location', '$routeParams', 'User',
  function($scope, $location, $routeParams, User) {
    $scope.adminLevel = adminLevel;
    User.get({email: $routeParams.email}, function(res) {
      if (res.success) {
        $scope.user = res.data;
      }else{
        alert('Get user error: ' + res.error);
      }
    });

    $scope.saveUser = function () {
      var user = angular.copy($scope.user);
      $scope.isSaving = true;
      User.update({id: user.Id}, user, 
        function (res) {
          if (res.success === true) {
            $location.path("/users/" + res.data.Email);
          } else {
            alert(res.error);
          }
          $scope.isSaving = false;
        });
    }
  }]);

/*
add new user
 */
aproxyControllers.controller('UserAddNewCtrl', 
  ['$scope', '$location', 'User', 
  function($scope, $location, User) {
    $scope.adminLevel = adminLevel;
    $scope.user = new User();
    $scope.user.AdminLevel = 0

    $scope.saveUser = function () {
      var user = angular.copy($scope.user);
      if (!user.Email || !user.Pwd) {
        alert('please fill all items.');
        return;
      }
      $scope.isSaving = true;
      user.$save(null, function (res) {
        if (res.success) {
          $location.path("/users/" + res.data.Email);
        } else {
          alert(res.error);
        }
        $scope.isSaving = false;
      });
    };

  }]);
