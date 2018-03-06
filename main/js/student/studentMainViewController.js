app.controller('mainViewController', function($scope, $rootScope, $timeout, $mdDialog) {
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
    getMethod("/student/getteachers", {current: currentWeek}, controller.updateTeachers);
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
  
  controller.openDialog = function(ev, block) {
  	controller.rootScope.ev = ev;
  	location.hash='#change'+block;
  };
  
  controller.requestTeachers(false);
});