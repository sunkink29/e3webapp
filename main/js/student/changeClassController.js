app.controller('changeClassController', function($scope, $rootScope, $mdDialog) {
  var controller = this;
  controller.rootScope = $rootScope;
  $rootScope.changeClass = this;
  controller.classes = [{name:'Loading'}];
  controller.blockChange = 0
  
  controller.requestTeachers = function () {
    controller.classes = [{name:'Loading'}];
    getMethod("/teacher/getall", {current:false}, controller.updateTeachers);
  };
  
  controller.updateTeachers = function (teachers) {
    teachers.forEach(function(teacher, index) {
      teachers[index].blocks = [teacher.Block1, teacher.Block2]
    })
    controller.classes = teachers;
    $scope.$apply();
  };

  controller.onClassEdit = function(teachers) {
    teachers.forEach(function(teacher, index) {
      controller.classes.forEach(function(element, index) {
        if (element.Email == teacher.Email) {
          controller.classes[index] = teacher
        }
      });
      teachers[index].blocks = [teacher.Block1, teacher.Block2]
    });
    $scope.$apply();
  }
  
  controller.selectClass = function (teacher) {
    teacher.curBlock = controller.blockChange;
    $rootScope.mainView.nextClasses[teacher.curBlock] = teacher;
    postMethod("/student/setteacher", {"ID": teacher.ID,"Block": teacher.curBlock}, controller.updateTeachers);
    controller.closeDialog();
  };
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  };
  
  $(window).on('hashchange', function() {
	if (window.location.hash === "#change0") {
		controller.blockChange = 0
		controller.showDialog();
	}else if (window.location.hash === "#change1") {
		controller.blockChange = 1
		controller.showDialog();
	}
    $scope.$apply();
  });
  
  controller.showDialog = function() {
    $mdDialog.show({
      contentElement: '#change',
      parent: angular.element(document.body),
      targetEvent: $rootScope.ev,
      clickOutsideToClose: true,
      onRemoving: function() {
        controller.onClose();
        window.location.hash = "student";
      }
    });
  };
  
  controller.onClose = function() {
    
  };
  
  controller.requestTeachers();
});