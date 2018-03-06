app.controller('addUserController', function($scope, $rootScope, $mdDialog) {
  var controller = this;
  controller.newUser = {name: '', email: '', Teacher: '', Admin: ''};
  
  controller.rootScope = $rootScope;
  controller.rootScope.addUser = this;
  
  $(window).on('hashchange', function() {
	if (window.location.hash === "#addUser") {
		controller.showDialog();
	}
  });
  
  controller.showDialog = function() {
    $mdDialog.show({
      contentElement: '#addUser',
      parent: angular.element(document.body),
      targetEvent: $rootScope.ev,
      clickOutsideToClose: true,
      onRemoving: function() {
        controller.resetUser();
        window.location.hash = "admin";
      }
    });
  };
  
  controller.addUser = function(user) {
  	user.Teacher = user.Teacher || false;
  	user.Admin = user.Admin || false;
    postMethod("/admin/newuser",user);
    controller.closeDialog();
  };
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  };
  
  controller.resetUser = function() {
    controller.newUser = {name: '', email: '', Teacher: '', Admin: ''};
  };
});
