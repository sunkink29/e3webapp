app.controller('editClassController', function($scope, $rootScope, objService, $timeout, $mdDialog) {
  var controller = this;
  controller.rootScope = $rootScope;
  controller.rootScope.edit = this;
  controller.currentBlockInfo = [{Subject:'',Description:'',RoomNumber:'',CurSize: '',MaxSize:'',BlockOpen:''},
  {Subject:'',Description:'',RoomNumber:'',CurSize: '',MaxSize:'',BlockOpen:''}];
  
  controller.updateBlockInfo = function() {
    callMethod("getBlocks", null, controller.showBlockInfo);
  };
  
  controller.showBlockInfo = function(message) {
    controller.currentBlockInfo[0] = message[0];
    controller.currentBlockInfo[1] = message[1];
    $scope.$apply();
  };
  
  controller.submitEdit = function() {
    callMethod("setBlocks", controller.currentBlockInfo, controller.showBlockInfo);
    controller.closeDialog();
  };
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  };
  
  controller.onClose = function() {
//    controller.updateBlockInfo();
  };
  
  controller.updateBlockInfo();
});