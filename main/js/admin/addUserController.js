app.controller('addUserController', function($scope, $rootScope, $mdDialog) {
  var controller = this;
  controller.newUser = {name: '', email: '', Teacher: '', Admin: ''};
  
  controller.rootScope = $rootScope;
  controller.rootScope.addUser = this;
  
  controller.addUser = function(user) {
  	user.Teacher = user.Teacher || false;
  	user.Admin = user.Admin || false;
    callMethod("newUser",user);
    controller.closeDialog();
  };
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  };
  
  controller.resetUser = function() {
    controller.newUser = {name: '', email: '', Teacher: '', Admin: ''};
  };
});
