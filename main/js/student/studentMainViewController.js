app.controller('mainViewController', function($scope, $rootScope,objService, $timeout, $mdDialog) {
  var controller = this;
  controller.nextClasses = [{name:'Loading'}];
  controller.currentClasses = [{name:'Loading'}];
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
    callMethod("getCurrentStudentBlocks", currentWeek, controller.updateTeachers);
  };
  
  controller.updateTeachers = function (teachers) {
    teachers.forEach(function(item, index) {
      if (item === null) {
        teachers[index] = {Name:"Unassigned", ID: -1,Block1:{BlockOpen:true},Block2:{BlockOpen:true}};
      }
      teachers[index].curBlock = index;
    });
    teachers[1].Block1 = teachers[1].Block2;
    if (controller.currentWeek) {
      controller.currentClasses = teachers;
    } else {
      controller.nextClasses = teachers;
      controller.requestTeachers(true);
    }
    $scope.$apply();
  };
  
  controller.changeClass = function(block) {
    showPage('changeClass');
    $rootScope.requestedChange = controller.nextClasses[block];
  };
  
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
  };
  
  controller.requestTeachers(false);
});