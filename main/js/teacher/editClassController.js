app.controller('editClassController', function($scope, $rootScope, $timeout, $mdDialog) {
  var controller = this;
  controller.rootScope = $rootScope;
  controller.rootScope.edit = this;
  controller.currentBlockInfo = [{Subject:'',Description:'',RoomNumber:'',CurSize: '',MaxSize:'',BlockOpen:''},
  {Subject:'',Description:'',RoomNumber:'',CurSize: '',MaxSize:'',BlockOpen:''}];
  controller.addButton0 = true;
  controller.addButton1 = true;
  
  controller.updateBlockInfo = function() {
    getMethod("/teacher/getblocks", null, controller.showBlockInfo);
  };
  
  controller.showBlockInfo = function(message) {
    controller.currentBlockInfo[0] = message[0];
    controller.currentBlockInfo[1] = message[1];
    $scope.$apply();
  };
  
  controller.submitEdit = function() {
    postMethod("/teacher/setblocks", controller.currentBlockInfo);
    controller.closeDialog();
  };
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  };
  
  $(window).on('hashchange', function() {
	if (window.location.hash === "#editBlock0") {
		$rootScope.block = 0;
		controller.showDialog();
	} else if (window.location.hash === "#editBlock1") {
		$rootScope.block = 1;
		controller.showDialog();
	}
  });
  
  controller.showDialog = function() {
    $mdDialog.show({
      contentElement: '#editBlock',
      parent: angular.element(document.body),
      targetEvent: $rootScope.ev,
      clickOutsideToClose: true,
      onRemoving: function() {
        window.location.hash = "teacher";
      }
    });
  };
  
  controller.updateBlockInfo();
});