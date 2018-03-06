app.controller('adminMainViewController', function($scope, $rootScope, $mdDialog) {
  var controller = this;
  $rootScope.adminControl = this;
  controller.rootScope = $rootScope;
  
  controller.userList;
  
  controller.updateUsers = function() {
   	getMethod("/admin/getallusers", null, controller.showUsers);
  };
  
  controller.showUsers = function(users) {
    $scope.$apply(function() {
      controller.userList = users;
    });
  };
});