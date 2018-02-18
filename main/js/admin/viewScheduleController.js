app.controller('viewScheduleController', function($scope, $rootScope, objService, $mdDialog) {
  var controller = this;
  controller.selectedTeacher;
  controller.nextTeacherList;
  controller.currentTeacherList;
  controller.studentList;
  controller.showBack = false;
  
  controller.rootScope = $rootScope;
  controller.rootScope.viewSchedule = this;
  
  controller.updateTeachers = function(currentWeek) {
    if (currentWeek) {
      controller.currentClasses = [{name:'Loading'}];
    } else {
      controller.nextClasses = [{name:'Loading'}];
    }
    controller.currentWeek = currentWeek;
    callMethod("getAllTeachers", currentWeek, controller.showTeachers);
  };
  
  controller.showTeachers = function(teachers) {
    $scope.$apply(function() {
      if (controller.currentWeek) {
      controller.currentTeacherList = teachers;
    } else {
      controller.nextTeacherList = teachers;
      controller.updateTeachers(true);
    }
    });
  };
  
  controller.updateStudents = function(teacher) {
    controller.studentList = [[{name:"Loading"}]];
    callMethod("getStudentsInClass", teacher.ID, controller.showStudents);
  };
  
  controller.showStudents = function(students) {
    $scope.$apply(function() {
      controller.studentList = students;
    });
  };
  
  controller.showTeacherStudentList = function(teacher, currentWeek) {
    controller.selectedTeacher = teacher;
    controller.selectedCurrentWeek = currentWeek;
    $("#viewTeachers").hide();
    $("#viewStudentList").show();
    controller.showBack = true;
    controller.updateStudents(teacher, currentWeek);
  };
  
  controller.hideTeacherStudentList = function() {
    $("#viewTeachers").show();
    $("#viewStudentList").hide();
    controller.showBack = false;
  };
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  };
  
  controller.updateTeachers(false);
});