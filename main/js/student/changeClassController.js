app.controller('changeClassController', function($scope, $rootScope, $mdDialog) {
  var controller = this;
  controller.rootScope = $rootScope;
  $rootScope.changeClass = $scope;
  controller.classes = [{name:'Loading'}];
  
  controller.requestTeachers = function () {
    controller.classes = [{name:'Loading'}];
    getMethod("/teacher/getall", {current:false}, controller.updateTeachers);
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
    postMethod("/student/setteacher", {"ID": teacher.ID,"Block": teacher.curBlock}, controller.updateTeachers);
    controller.closeDialog();
  };
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  };
  
  $(window).on('hashchange', function() {
	if (window.location.hash === "#change0") {
		$rootScope.requestedChange = $rootScope.mainView.nextClasses[0];
		$rootScope.requestedChange.curBlock = 0;
		controller.showDialog();
	}else if (window.location.hash === "#change1") {
		$rootScope.requestedChange = $rootScope.mainView.nextClasses[1];
		$rootScope.requestedChange.curBlock = 1;
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