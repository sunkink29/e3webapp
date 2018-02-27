app.controller('assignStudentController', function($scope, $rootScope, $timeout, $mdDialog) {
  var controller = this;
  $rootScope.assign = this;
  controller.studentList = [{id:-1,name:'Loading'}];
  controller.rootScope = $rootScope;
  controller.grade;
  controller.selectedStudent;
  controller.submitDisabled = true;
  controller.selectDisabled = false;
  
  controller.querySearch = function(query) {
    return query ? controller.studentList.filter( controller.createFilterFor(query) ) : controller.studentList;
  };
  
  controller.createFilterFor = function(query) {
    var lowercaseQuery = angular.lowercase(query);

    return function filterFn(student) {
      return angular.lowercase(student.Name).indexOf(lowercaseQuery) === 0;
    };
  };
  
  controller.updateStudents = function() {
    controller.selectDisabled = false;
    controller.selectedStudent = null;
    callMethod("getAllStudents", false, controller.showStudents);
  };
  
  controller.showStudents = function(students) {
    $scope.$apply(function() {
      students.forEach(function(item, index) {
        students[index].curBlock = index;
      });
      controller.studentList = students;
//      if (controller.grade != '') {
//        controller.selectDisabled = false;
//      }
    });
  };
  
  controller.addStudent = function() {
    var selectedStudent = controller.selectedStudent;
    var studentId = selectedStudent.ID;
    var block = $rootScope.block;
    callMethod("addStudentToClass", {key: studentId, Block: block}, controller.showStudents);
    var studentTable = $rootScope.mainView.NextStudents;
    studentTable[block][studentTable[block].length] = selectedStudent;
    controller.closeDialog();
  };
  
  controller.changeGrade = function() {
    controller.updateStudents();
    controller.submitDisabled = true;
  };
  
  controller.GetPreviousOpen = function() {
  	callMethod("studentClassOpen", {StdntID: controller.selectedStudent.ID, Block: $rootScope.block}, 
  			controller.changeStudent);
  }
  
  controller.changeStudent = function(previousOpen) {
    var selectedStudent = controller.selectedStudent;
    if (selectedStudent !== null) {
	    var block = $rootScope.block;
	    controller.submitDisabled = true;
	    var classFull = $rootScope.edit.currentBlockInfo[block].CurSize >= $rootScope.edit.currentBlockInfo[block].MaxSize;
	    if (!controller.selectDisabled && selectedStudent !== null) {
	      controller.submitDisabled = !previousOpen || classFull;
	    }
	}

  };
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  };
  
  controller.onClose = function() {
    controller.grade = null;
    controller.selectDisabled = false;
    controller.selectedStudent = null;
    controller.submitDisabled = true;
  };
  controller.updateStudents()
});