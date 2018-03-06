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
    getMethod("/student/getall", {current: false}, controller.showStudents);
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
    postMethod("/teacher/addstudent", {Key: studentId, Block: block}, controller.showStudents);
    var studentTable = $rootScope.mainView.NextStudents;
    studentTable[block][studentTable[block].length] = selectedStudent;
    controller.closeDialog();
  };
  
  controller.changeGrade = function() {
    controller.updateStudents();
    controller.submitDisabled = true;
  };
  
  controller.GetPreviousOpen = function() {
  	getMethod("/student/open", {id: controller.selectedStudent.ID, Block: $rootScope.block}, 
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
	$scope.$apply();
  };
  
  controller.closeDialog = function() {
    $mdDialog.hide();
  };
  
  $(window).on('hashchange', function() {
	if (window.location.hash === "#assign0") {
		$rootScope.block = 0;
		controller.showDialog();
	} else if (window.location.hash === "#assign1") {
		$rootScope.block = 1;
		controller.showDialog();
	}
  });
  
  controller.showDialog = function() {
    $mdDialog.show({
      contentElement: '#assign',
      parent: angular.element(document.body),
      targetEvent: $rootScope.ev,
      clickOutsideToClose: true,
      onRemoving: function() {
        controller.onClose();
        window.location.hash = "teacher";
      },
    });
  };
  
  controller.onClose = function() {
    controller.grade = null;
    controller.selectDisabled = false;
    controller.selectedStudent = null;
    controller.submitDisabled = true;
  };
  controller.updateStudents()
});