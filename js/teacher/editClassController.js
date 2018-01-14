app.controller('editClassController', function($scope, $rootScope, objService, $timeout, $mdDialog) {
  var controller = this;
  controller.rootScope = $rootScope;
  controller.rootScope.edit = this;
  controller.currentBlockInfo = [{subject:'',description:'',roomNumber:'',curSize: '',maxSize:'',status:''},
  {subject:'',description:'',roomNumber:'',curSize: '',maxSize:'',status:''}];
  
  controller.updateBlockInfo = function() {
    var block = $rootScope.block;
    google.script.run.withSuccessHandler(controller.showBlockInfo).getBlockInfo();
  }
  
  controller.showBlockInfo = function(message) {
    controller.currentBlockInfo[0] = objService.getBlockInfoObjFromList(message.blockInfo.slice(0,6));
    controller.currentBlockInfo[1] = objService.getBlockInfoObjFromList(message.blockInfo.slice(6,12));
    controller.currentBlockInfo[0].blockNum = 0;
    controller.currentBlockInfo[1].blockNum = 1;
    $scope.$apply();
  }
  
  controller.submitEdit = function() {
    google.script.run.editBlockInfo(controller.currentBlockInfo[$rootScope.block]);
    controller.closeDialog();
  }
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  }
  
  controller.onClose = function() {
    controller.updateBlockInfo();
  }
  
  controller.updateBlockInfo()
})