app.controller('editClassController', function($scope, $rootScope, $timeout, $mdDialog) {
  var controller = this;
  controller.rootScope = $rootScope;
  controller.rootScope.edit = this;
  controller.currentBlockInfo = [{Subject:'',Description:'',RoomNumber:'',CurSize: '',MaxSize:'',BlockOpen:''},
  {Subject:'',Description:'',RoomNumber:'',CurSize: '',MaxSize:'',BlockOpen:''}];
  controller.addButton1 = true;
  controller.addButton2 = true;
  
  controller.updateBlockInfo = function() {
    callMethod("getBlocks", null, controller.showBlockInfo);
  };
  
  controller.showBlockInfo = function(message) {
    controller.currentBlockInfo[0] = message[0];
    controller.currentBlockInfo[1] = message[1];
    controller.checkAddButton();
    $scope.$apply();
  };
  
  controller.submitEdit = function() {
    callMethod("setBlocks", controller.currentBlockInfo, controller.showBlockInfo);
    controller.checkAddButton();
    controller.closeDialog();
  };
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  };
  
  controller.checkAddButton = function() {
  	if (controller.currentBlockInfo[0].MaxSize > 0) {
    	controller.addButton0 = false;
    }
    if (controller.currentBlockInfo[1].MaxSize > 0) {
    	controller.addButton1 = false;
    }  
  };
  
  controller.onClose = function() {
//    controller.updateBlockInfo();
  };
  
  controller.updateBlockInfo();
});