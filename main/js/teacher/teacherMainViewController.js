app.controller('teacherMainViewController', function($scope, $rootScope, objService, $timeout, $mdDialog, $mdToast) {
  var controller = this;
  $rootScope.mainView = controller;
  controller.NextStudents = [[{ID:-1,name:'loading'},{ID:-1,name:''}]];
  controller.currentStudents = [[{ID:-1,name:'loading'},{ID:-1,name:''}]];
  
  controller.updateStudents = function (currentWeek) {
    if (currentWeek) {
      controller.currentStudents = [[{id:-1,name:'loading'},{id:-1,name:''}]];
    } else {
      controller.NextStudents = [[{id:-1,name:'loading'},{id:-1,name:''}]];
    }
    controller.currentWeek = currentWeek;
    callMethod("getCurrentStudents", currentWeek, controller.showStudents);
  };
  
  controller.showStudents = function(students) {
    if (controller.currentWeek) {
      controller.currentStudents = students;
    } else {
      controller.NextStudents = students;
      controller.updateStudents(true);
    }
    $scope.$apply();
  };
  
  controller.showAssignDialog = function(ev,block) {
    controller.setBlock(block);
    $mdDialog.show({
      contentElement: '#assign',
      parent: angular.element(document.body),
      targetEvent: ev,
      clickOutsideToClose: true,
      onRemoving: function() {
        $rootScope.assign.onClose();
      }
    });
  };
  
  controller.showEditDialog = function(ev,block) {
    controller.setBlock(block);
    $mdDialog.show({
      contentElement: '#edit',
      parent: angular.element(document.body),
      targetEvent: ev,
      clickOutsideToClose: true,
      onRemoving: function() {
        $rootScope.edit.onClose();
      }
    });
  };
  
  controller.setBlock = function (block) {
    $rootScope.block = block;      
  };
  
  controller.removeStudent = function(index, block){
    var studentId = controller.NextStudents[block][index].ID;
    callMethod("removeFromClass", {Key: studentId, Block: block}, controller.showStudents);
    var list = controller.NextStudents;
    list[block].splice(index, 1);
  };
  
  controller.updateStudents(false);
});