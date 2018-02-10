app.controller('changeClassController', function($scope, $rootScope,objService, $mdDialog) {
  var controller = this;
  controller.rootScope = $rootScope;
  $rootScope.changeClass = $scope;
  controller.classes = [{name:'Loading'}];
  
  controller.requestTeachers = function () {
    controller.classes = [{name:'Loading'}];
    callMethod("getAllTeachers", false, controller.updateTeachers);
  };
  
  controller.updateTeachers = function (teachers) {
    controller.classes = teachers;
    $scope.$apply();
  };
  
  controller.selectClass = function (teacher) {
    teacher.curBlock = $rootScope.requestedChange.curBlock;
    controller.requestTeachers();
    if (teacher.curBlock === 1) {
    	teacher.Block1 = teacher.Block2;
    }
    teacher.Block1.CurSize++;
    $rootScope.mainView.nextClasses[teacher.curBlock] = teacher;
    callMethod("setTeacher", {Teacher: teacher.Email,"Block": teacher.curBlock}, controller.updateTeachers);
    controller.closeDialog();
  };
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  };
  
  controller.onClose = function() {
    
  };
  
  controller.requestTeachers();
});