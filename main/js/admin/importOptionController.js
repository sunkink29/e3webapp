function showImportOptions() {
    window.location.hash = "#importOptions"
}

app.controller('importOptionController', function($scope, $rootScope, $mdDialog) {
    var ctrl = this;
    ctrl.options = {tabName: "user data", headerNames: 
        {name: "Name", email: "Email", teacher: "Teacher", admin: "Admin", grade: "Grade"}}
    ctrl.completed = 0
    ctrl.total = 1
    ctrl.rootScope = $rootScope;
    ctrl.rootScope.importOptions = this;
    ctrl.scope = $scope;
    
    $(window).on('hashchange', function() {
      if (window.location.hash === "#importOptions") {
          ctrl.showDialog();
      }
    });
    
    ctrl.showDialog = function() {
      $mdDialog.show({
        contentElement: '#importOptions',
        parent: angular.element(document.body),
        targetEvent: $rootScope.ev,
        clickOutsideToClose: true,
        onRemoving: function() {
          window.location.hash = "admin";
        }
      });
    };
    
    ctrl.submitOptions = function() {
      options = ctrl.options;
      var output = {id: spreadSheetID, options: options};
      postMethod("/admin/importusers", output, redirect);
      ctrl.closeDialog();
      ctrl.openProgress();
    };
    
    ctrl.closeDialog = function() {
      $mdDialog.hide();
    };

    ctrl.openProgress = function() {
      $mdDialog.show({
        contentElement: '#importProgress',
        parent: angular.element(document.body),
        targetEvent: $rootScope.ev,
        clickOutsideToClose: true,
        onRemoving: function() {
          window.location.hash = "admin";
        }
      });
      setTimeout(function() { getMethod("/admin/getimportprogress", null, ctrl.updateProgress);}, 5000)
    };

    ctrl.updateProgress = function(data) {
      ctrl.completed = data.Completed;
      ctrl.total = data.Total;
      if (ctrl.completed < ctrl.total) {
        setTimeout(function() { getMethod("/admin/getimportprogress", null, ctrl.updateProgress);}, 5000)
      }
      // console.log("completed: "+ctrl.completed+"  total: "+ctrl.total)
      $scope.$apply();
    };
  });
  