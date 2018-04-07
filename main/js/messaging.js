app.controller('messaging', function($mdDialog, $rootScope) {
    var ctrl = this;
    const messaging = firebase.messaging();
    messaging.usePublicVapidKey(firebaseKey);

    messaging.requestPermission().then(function() {
      console.log('Notification permission granted.');
      messaging.getToken().then(function(currentToken) {
        if (currentToken) {
          console.log(currentToken);
          postMethod("/registertoken", currentToken);
        } else {
          console.log('No Instance ID token available. Request permission to generate one.');
        }
      }).catch(function(err) {
        console.log('An error occurred while retrieving token. ', err);
      });
    }).catch(function(err) {
      console.log('Unable to get permission to notify.', err);
      ctrl.openWarning()
    });

    // Callback fired if Instance ID token is updated.
    messaging.onTokenRefresh(function() {
      messaging.getToken().then(function(refreshedToken) {
        console.log('Token refreshed.');
        postMethod("/registertoken", currentToken);
      }).catch(function(err) {
        console.log('Unable to retrieve refreshed token ', err);
      });
    });

    function onMessage (payload) {
      console.log('Message received. ', payload);
      var data = true
      if (payload.data.data) {
        data = JSON.parse(payload.data.data);
      }
      if (payload.data.event === "popup") {
        alert = $mdDialog.alert({
          title: data.title,
          textContent: data.message,
          ok: 'Close'
        });
        $mdDialog.show( alert );
      } else if (isTeacher) {
        if (payload.data.event === "studentUpdate") {
          $rootScope.mainView.onStudentUpdate(data)
        }
      } else {
        if (payload.data.event === "classEdit") {
          $rootScope.changeClass.onClassEdit(data);
          var data2 = {Block: -1, Teachers: data}
          $rootScope.mainView.onCurrentChange(data2);
        } else if (payload.data.event === "currentChange") {
          $rootScope.mainView.onCurrentChange(data);
        }
      }
      
    }
    
    messaging.onMessage(onMessage)
    // use page visibility api to update after focus
    
    ctrl.openWarning = function() {
      alert = $mdDialog.alert({
        title: 'Error',
        textContent: 'Permission to receive notifications is required to get automatic updates',
        ok: 'Close'
      });
      $mdDialog.show( alert );
    };
  });