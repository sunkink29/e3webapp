app.controller('editUserController', function($scope, $rootScope, objService, $mdDialog) {
  var controller = this;
  controller.selectedUser;
  
  controller.rootScope = $rootScope;
  controller.rootScope.editUser = this;
  
  controller.querySearch = function(query) {
    return query ? $rootScope.adminControl.userList.filter( controller.createFilterFor(query) ) : $rootScope.adminControl.userList;
  };
  
  controller.createFilterFor = function(query) {
    var lowercaseQuery = angular.lowercase(query);

    return function filterFn(user) {
      return angular.lowercase(user.Name).indexOf(lowercaseQuery) === 0;
    };
  };
  
  controller.editUser = function(user) {
    callMethod("editUser", user, null);
    $rootScope.adminControl.updateUsers();
    controller.closeDialog();
  };
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  };
  
  controller.resetUser = function() {
    controller.selectedUser = "";
  };
});