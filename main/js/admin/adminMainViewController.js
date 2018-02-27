app.controller('adminMainViewController', function($scope, $rootScope, $mdDialog) {
  var controller = this;
  $rootScope.adminControl = this;
  
  controller.userList;
  
  controller.updateUsers = function() {
   	callMethod("getAllUsers", null, controller.showUsers);
  };
  
  controller.showUsers = function(users) {
    $scope.$apply(function() {
      controller.userList = users;
    });
  };
  
  controller.showAddDialog = function(ev) {
    $mdDialog.show({
      contentElement: '#addUser',
      parent: angular.element(document.body),
      targetEvent: ev,
      clickOutsideToClose: true,
      onRemoving: function() {
        $rootScope.addUser.resetUser();
      }
    });
  };
  
  controller.showEditDialog = function(ev) {
    $mdDialog.show({
      contentElement: '#editUser',
      parent: angular.element(document.body),
      targetEvent: ev,
      clickOutsideToClose: true,
      onRemoving: function() {
        $rootScope.editUser.resetUser();
      },
      onShowing: function() {
      	controller.updateUsers();
      }
    });
  };
  
  controller.showRemoveDialog = function(ev) {
    $mdDialog.show({
      contentElement: '#removeUser',
      parent: angular.element(document.body),
      targetEvent: ev,
      clickOutsideToClose: true,
      onRemoving: function() {
        $rootScope.removeUser.resetUser();
      },
      onShowing: function() {
      	controller.updateUsers();
      }
    });
  };
  
  controller.showScheduleDialog = function(ev) {
    $mdDialog.show({
      contentElement: '#viewSchedule',
      parent: angular.element(document.body),
      targetEvent: ev,
      clickOutsideToClose: true,
      onShowing: function() {
      	$rootScope.viewSchedule.updateTeachers(false);
      }
    });
  };
});