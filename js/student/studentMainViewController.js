app.controller('mainViewController', function($scope, $rootScope,objService, $timeout, $mdDialog) {
  var controller = this;
  controller.nextClasses = [{name:'Loading'}];
  controller.currentClasses = [{name:"Unassigned", id: -1, curBlock: 0, blocks: [{status: 'closed'}, {status: 'closed'}]}
                                ,{name:"Unassigned", id: -1, curBlock: 1, blocks: [{status: 'closed'}, {status: 'closed'}]}];
  controller.rootScope = $rootScope;
  $rootScope.mainView = this;
  $rootScope.mainViewScope = $scope;
  $rootScope.requestedChange = {curBlock: 0};
  
  controller.requestTeachers = function (currentWeek) {
    if (currentWeek) {
      controller.currentClasses = [{name:'Loading'}];
    } else {
      controller.nextClasses = [{name:'Loading'}];
    }
    controller.currentWeek = currentWeek;
    // google.script.run.withSuccessHandler(controller.updateTeachers).getTeachers(true,currentWeek);
  }
  
  controller.updateTeachers = function (teachers) {
    var list = [];
    teachers.forEach(function(item, index) {
      var teacherObj = objService.getTeacherObjFromList(item);
      teacherObj.curBlock = index;
      list.push(teacherObj);
    })
    if (controller.currentWeek) {
      controller.currentClasses = list;
    } else {
      controller.nextClasses = list;
      controller.requestTeachers(true);
    }
    $scope.$apply();
  }
  
  controller.changeClass = function(block) {
    showPage('changeClass');
    $rootScope.requestedChange = controller.nextClasses[block];
  }
  
  controller.showChangeDialog = function(ev,block) {
    $rootScope.requestedChange = controller.nextClasses[block];
    $mdDialog.show({
      contentElement: '#change',
      parent: angular.element(document.body),
      targetEvent: ev,
      clickOutsideToClose: true,
      onRemoving: function() {
        $rootScope.changeClass.onClose();
      }
    });
  }
  
  controller.requestTeachers(false);
});