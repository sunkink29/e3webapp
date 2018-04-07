app.controller('teacherMainViewController', function($scope, $rootScope, $timeout, $mdDialog, $mdToast) {
  var controller = this;
  $rootScope.mainView = controller;
  controller.NextStudents = [[],[]];
  controller.currentStudents = [[],[]];
  controller.rootScope = $rootScope;
  
  controller.updateStudents = function (currentWeek) {
    if (currentWeek) {
      controller.currentStudents = [[],[]];
    } else {
      controller.NextStudents = [[],[]];
    }
    controller.currentWeek = currentWeek;
    getMethod("/teacher/getstudents", {current: currentWeek}, controller.showStudents);
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

  controller.onStudentUpdate = function(data) {
    if (data.Method == "add") {
      controller.NextStudents[data.Block].push(data.Student)
    } else if (data.Method == "remove") {
      for(var i = 0, len = controller.NextStudents[data.Block].length;i < len; i++){
        if (data.Student.Email === controller.NextStudents[data.Block][i].Email) {
          var list = controller.NextStudents;
          list[data.Block].splice(i, 1);
        }
      }
    }
    $scope.$apply();
  }
  
  controller.openDialog = function(ev, hash, block) {
  	controller.rootScope.ev = ev; 
  	window.location.hash = hash+block;
  };
  
  controller.removeStudent = function(index, block){
    var studentId = controller.NextStudents[block][index].ID;
    postMethod("/teacher/removestudent", {Key: studentId, Block: block}, controller.showStudents);
    var list = controller.NextStudents;
    list[block].splice(index, 1);
  };
  
  controller.updateStudents(false);
});