app.controller('mainViewController', function($scope, $rootScope, $timeout, $mdDialog) {
  var controller = this;
  controller.nextClasses = [{name:'Loading'}];
  controller.currentClasses = [{name:'Loading'}];
  controller.rootScope = $rootScope;
  $rootScope.mainView = this;
  $rootScope.mainViewScope = $scope;
  
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
      teachers[index].blocks = [teachers[index].Block1, teachers[index].Block2]
    });
    //teachers[1].Block1 = teachers[1].Block2;
    if (controller.currentWeek) {
      controller.currentClasses = teachers;
    } else {
      controller.nextClasses = teachers;
      controller.requestTeachers(true);
    }
    $scope.$apply();
  };

  controller.onCurrentChange = function(data) {
    if (data.Block == -1) {
      controller.nextClasses.forEach(function(cTeacher, index){
        data.Teachers.forEach(function(nTeacher) {
          if (cTeacher.Email == nTeacher.Email) {
            controller.nextClasses[index] = nTeacher
            controller.nextClasses[index].curBlock = index;
            controller.nextClasses[index].blocks = [nTeacher.Block1, nTeacher.Block2]
          }
        })
      })
    } else {
      if (data.Teacher) {
        controller.nextClasses[data.Block] = data.Teacher;
        controller.nextClasses[data.Block].blocks = [data.Teacher.Block1, data.Teacher.Block2]
      } else {
        controller.nextClasses[data.Block] = {Name:"Unassigned", ID: -1,blocks:[{BlockOpen:true}, {BlockOpen:true}]};
      }
      controller.nextClasses[data.Block].curBlock = data.Block;
    }
    
    $scope.$apply();
  }
  
  controller.openDialog = function(ev, block) {
  	controller.rootScope.ev = ev;
  	location.hash='#change'+block;
  };
  
  controller.requestTeachers(false);
});