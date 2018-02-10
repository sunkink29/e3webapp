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
  
  controller.updateStudents = function(teacher, currentWeek) {
    controller.studentList = [[{name:"Loading"}]];
    // google.script.run.withSuccessHandler(controller.showStudents).getStudents(false, undefined, currentWeek, teacher.id);
  };
  
  controller.showStudents = function(students) {
    $scope.$apply(function() {
      students.forEach(function callback(column, outerIndex) {
        column.forEach(function callback(student, innerIndex) {
          if (student != null) { students[outerIndex][innerIndex] = objService.getStudentObjFromList(student); }
        });
      });
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