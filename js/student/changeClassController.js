app.controller('changeClassController', function($scope, $rootScope,objService, $mdDialog) {
  var controller = this;
  controller.rootScope = $rootScope;
  $rootScope.error = '';
  $rootScope.changeClass = $scope;
  controller.classes = [{name:'Loading'}];
  
  controller.requestTeachers = function () {
    controller.classes = [{name:'Loading'}];
    // google.script.run.withSuccessHandler(controller.updateTeachers).getTeachers();
  }
  
  controller.updateTeachers = function (teachers) {
    var list = [];
    teachers.forEach(function(item, index) {
      var teacherObj = objService.getTeacherObjFromList(item);
      list.push(teacherObj);
    })
    controller.classes = list;
    $scope.$apply();
  }
  
  controller.selectClass = function (teacher) {
    teacher.curBlock = $rootScope.requestedChange.curBlock;
    if ($rootScope.requestedChange.id >= 0) {
      controller.classes[$rootScope.requestedChange.id].blocks[teacher.curBlock].curSize--;
    }
    teacher.blocks[teacher.curBlock].curSize++;
    $rootScope.mainView.nextClasses[teacher.curBlock] = teacher;
    // google.script.run.withSuccessHandler(controller.handleError).setClass(teacher.curBlock, teacher.id);
    $rootScope.error = '';
    controller.closeDialog()
  }
  
  controller.handleError = function (message) {
    if (message.succeed == false) {
      Materialize.toast(message.error, 4000, 'error');
      $rootScope.mainView.requestTeachers();
      controller.requestTeachers();
    }
  }
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  }
  
  controller.onClose = function() {
    
  }
  
  controller.requestTeachers();
});