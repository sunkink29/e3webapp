app.controller('editUserController', function($scope, $rootScope, objService, $mdDialog) {
  var controller = this;
  controller.selectedUser;
  controller.userList;
  
  controller.rootScope = $rootScope;
  controller.rootScope.editUser = this;
  
  controller.querySearch = function(query) {
    return query ? controller.userList.filter( controller.createFilterFor(query) ) : controller.userList;
  }
  
  controller.createFilterFor = function(query) {
    var lowercaseQuery = angular.lowercase(query);

    return function filterFn(user) {
      return (angular.lowercase(user.name).indexOf(lowercaseQuery) === 0);
    };
  }
  
  controller.updateUsers = function() {
    controller.selectedStudent = null;
    //google.script.run.withSuccessHandler(controller.showUsers).getUserList();
  }
  
  controller.showUsers = function(users) {
    $scope.$apply(function() {
      var list = [];
      users.forEach(function(item, index) {
        var userObj = objService.getUserObjFromList(item);
        list.push(userObj);
      })
      controller.userList = list;
    })
  }
  
  controller.editUser = function(user) {
    //google.script.run.editUser(user);
    controller.closeDialog();
  }
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  }
  
  controller.resetUser = function() {
    controller.selectedUser = "";
  }
  
  controller.updateUsers();
})