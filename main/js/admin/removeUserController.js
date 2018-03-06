app.controller('removeUserController', function($scope, $rootScope, $mdDialog) {
  var controller = this;
  controller.selectedUser;
  
  controller.rootScope = $rootScope;
  controller.rootScope.removeUser = this;
  
  controller.querySearch = function(query) {
    return query ? $rootScope.adminControl.userList.filter( controller.createFilterFor(query) ) : $rootScope.adminControl.userList;
  };
  
  controller.createFilterFor = function(query) {
    var lowercaseQuery = angular.lowercase(query);

    return function filterFn(user) {
      return angular.lowercase(user.Name).indexOf(lowercaseQuery) === 0;
    };
  };
  
  controller.removeUser = function() {
    postMethod("/admin/deleteUser", {ID: controller.selectedUser.ID}, null);
    controller.closeDialog();
  };
  
  $(window).on('hashchange', function() {
	if (window.location.hash === "#removeUser") {
		controller.showDialog();
	}
  });
  
  controller.showDialog = function() {
    $mdDialog.show({
      contentElement: '#removeUser',
      parent: angular.element(document.body),
      targetEvent: $rootScope.ev,
      clickOutsideToClose: true,
      onRemoving: function() {
        controller.resetUser();
        window.location.hash = "admin";
      },
      onShowing: function() {
      	$rootScope.adminControl.updateUsers();
      }
    });
  };
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  };
  
  controller.resetUser = function() {
    controller.selectedUser = "";
  };
});