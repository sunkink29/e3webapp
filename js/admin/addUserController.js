app.controller('addUserController', function($scope, $rootScope, objService, $mdDialog) {
  var controller = this;
  controller.newUser = {name: '', email: '', isTeacher: '', isAdmin: ''};
  
  controller.rootScope = $rootScope;
  controller.rootScope.addUser = this;
  
  controller.addUser = function(user) {
    google.script.run.addUser(user);
    controller.closeDialog();
  }
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  }
  
  controller.resetUser = function() {
    controller.newUser = {name: '', email: '', isTeacher: '', isAdmin: ''};
  }
})
